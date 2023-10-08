package pkg

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type DiskStatusCollector struct {
	hostPhysicalDiskStatus *prometheus.Desc
	hostRaidStatus         *prometheus.Desc
	logger                 log.Logger
}

type Response struct {
	Controllers []struct {
		CommandStatus struct {
			Controller int    `json:"Controller"`
			Status     string `json:"Status"`
		} `json:"Command Status"`
		ResponseData struct {
			PhysicalDrives int `json:"Physical Drives"`
			PDList         []struct {
				Device   int    `json:"DID"`
				Position string `json:"EID:Slt"`
				State    string `json:"State"`
				Media    string `json:"Med"`
				Model    string `json:"Model"`
				Size     string `json:"Size"`
			} `json:"PD LIST"`
		} `json:"Response Data"`
	} `json:"Controllers"`
}

func NewDiskStatusCollector(promLog log.Logger) *DiskStatusCollector {
	return &DiskStatusCollector{
		hostPhysicalDiskStatus: prometheus.NewDesc(
			"host_physical_drives_status",
			"The host physical drives status check (0=abnormal, 1=normal)",
			[]string{"controller", "slot", "device", "model", "state", "media", "size"},
			nil,
		),
		hostRaidStatus: prometheus.NewDesc(
			"host_raid_card_status",
			"The host raid status check(0=abnormal, 1=normal)",
			[]string{"controller"},
			nil,
		),
		logger: promLog,
	}
}

func (d *DiskStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- d.hostPhysicalDiskStatus
	ch <- d.hostRaidStatus
}

func (d *DiskStatusCollector) Collect(ch chan<- prometheus.Metric) {
	response, err := getDiskStatusInfo()
	if err != nil {
		level.Error(d.logger).Log("msg", "Failed to fetch StorCLI output")
		return
	}

	var value float64
	value = 1
	for _, controller := range response.Controllers {
		for _, physicalDrive := range controller.ResponseData.PDList {
			if physicalDrive.State == "UBUnsp" {
				value = 0
			}
			ch <- prometheus.MustNewConstMetric(
				d.hostPhysicalDiskStatus, prometheus.GaugeValue, value,
				strconv.Itoa(controller.CommandStatus.Controller), physicalDrive.Position, strconv.Itoa(physicalDrive.Device), strings.TrimSpace(physicalDrive.Model),
				physicalDrive.State, physicalDrive.Media, physicalDrive.Size,
			)
		}
		if controller.CommandStatus.Status != "Success" {
			value = 0
		}
		ch <- prometheus.MustNewConstMetric(
			d.hostRaidStatus,
			prometheus.GaugeValue,
			value,
			strconv.Itoa(controller.CommandStatus.Controller),
		)
	}
	level.Info(d.logger).Log("msg", "collectd disk status success")
}

func getDiskStatusInfo() (resp Response, err error) {
	args := []string{
		"/call",
		"show",
		"J",
		"nolog",
	}
	results, err := Execute("storcli64", args...)
	if err != nil {
		return Response{}, fmt.Errorf("Execute storcli command failed %s ", err)
	}
	var response Response
	err = json.Unmarshal(results, &response)
	if err != nil {
		return Response{}, fmt.Errorf("Failed to unmarshal json %s ", err)
	}
	return response, nil
}
