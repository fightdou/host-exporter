package pkg

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type CpuCollector struct {
	hostCPUTempStatus *prometheus.Desc
	logger            log.Logger
}

func NewCpuCollector(promLog log.Logger) *CpuCollector {
	return &CpuCollector{
		hostCPUTempStatus: prometheus.NewDesc(
			"host_cpu_temp_status",
			"The host cpu temp health status check(0=abnormal, 1=normal)",
			[]string{"name"},
			nil,
		),
		logger: promLog,
	}
}

type CPUData struct {
	ID    int64
	Name  string
	State string
}

func (c *CpuCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hostCPUTempStatus
}

func (c *CpuCollector) Collect(ch chan<- prometheus.Metric) {
	args := []string{
		"-Q",
		"--ignore-unrecognized-events",
		"--comma-separated-output",
		"--no-header-output",
		"--sdr-cache-recreate",
		"--output-event-bitmask",
	}

	results, err := Execute("ipmimonitoring", args...)
	if err != nil {
		level.Error(c.logger).Log("msg", "Exec command ipmimonitoring failed", "error", err)
		return
	}

	cpuResult, err := getCPUSensorData(results)
	if err != nil {
		level.Error(c.logger).Log("msg", "Get cpu sensor data failed", "error", err)
	}
	for _, data := range cpuResult {
		var state float64
		switch data.State {
		case "Nominal":
			state = 1
		case "Warning":
			state = 0
		case "Critical":
			state = 0
		case "N/A":
			state = math.NaN()
		default:
			level.Error(c.logger).Log("msg", "Unknown sensor state", "state", data.State, "error", err)
			state = math.NaN()
		}
		level.Debug(c.logger).Log("msg", "Got values", "data", fmt.Sprintf("%+v", data))
		ch <- prometheus.MustNewConstMetric(
			c.hostCPUTempStatus,
			prometheus.GaugeValue,
			state,
			data.Name,
		)
	}
	level.Info(c.logger).Log("msg", "collectd cpu temp data success")
}

func getCPUSensorData(results []byte) ([]CPUData, error) {
	var cpuRes []CPUData
	r := csv.NewReader(bytes.NewReader(results))
	fields, err := r.ReadAll()
	if err != nil {
		return cpuRes, err
	}
	for _, line := range fields {
		var data CPUData
		data.Name = line[1]
		if strings.Contains(data.Name, "CPU") {
			data.ID, err = strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				return cpuRes, err
			}
			data.State = line[3]
			cpuRes = append(cpuRes, data)
		} else {
			continue
		}
	}
	return cpuRes, nil
}
