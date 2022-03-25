package pkg

import (
	"encoding/json"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type diskIOUtil struct {
	hostDiskIOUtil *prometheus.Desc
	hostDiskIOWait *prometheus.Desc
	logger         log.Logger
}

type DiskIOInfo struct {
	Sysstat struct {
		Hosts []struct {
			Statistics []struct {
				Disk []struct {
					DiskDevice string  `json:"disk_device"`
					Util       float64 `json:"util"`
					aWait      float64 `json:"await"`
				} `json:"disk"`
			} `json:"statistics"`
		} `json:"hosts"`
	} `json:"sysstat"`
}

func NewDiskIOUtil(promLog log.Logger) *diskIOUtil {
	return &diskIOUtil{
		hostDiskIOUtil: prometheus.NewDesc(
			"host_disk_io_util_percent",
			"The host disk io util percent",
			[]string{"disk"},
			nil,
		),
		hostDiskIOWait: prometheus.NewDesc(
			"host_disk_io_await",
			"The host disk io await time(ms)",
			[]string{"disk"},
			nil,
		),
		logger: promLog,
	}
}

func (d *diskIOUtil) Describe(ch chan<- *prometheus.Desc) {
	ch <- d.hostDiskIOUtil
	ch <- d.hostDiskIOWait
}

func (d *diskIOUtil) Collect(ch chan<- prometheus.Metric) {
	diskData := d.getDiskIOUtil()
	for _, hosts := range diskData.Sysstat.Hosts {
		for _, s := range hosts.Statistics {
			for _, disk := range s.Disk {
				ch <- prometheus.MustNewConstMetric(
					d.hostDiskIOUtil,
					prometheus.GaugeValue,
					disk.Util,
					disk.DiskDevice,
				)
				ch <- prometheus.MustNewConstMetric(
					d.hostDiskIOWait,
					prometheus.GaugeValue,
					disk.aWait,
					disk.DiskDevice,
				)
			}
		}
	}
}

func (d *diskIOUtil) getDiskIOUtil() DiskIOInfo {
	diskInfo := DiskIOInfo{}
	result, err := Execute("iostat", "-dxs", "-o", "JSON")
	if err != nil {
		level.Error(d.logger).Log("msg", "Exec iostat command failed", "err", err)
		return diskInfo
	}
	jErr := json.Unmarshal(result, &diskInfo)
	if jErr != nil {
		level.Error(d.logger).Log("msg", "Json Unmarshal failed", "err", jErr)
		return diskInfo
	}
	level.Info(d.logger).Log("msg", "Command exec success")
	return diskInfo
}
