package pkg

import (
	"bufio"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type DiskStatusCollector struct {
	hostDiskStatus *prometheus.Desc
	hostRaidStatus *prometheus.Desc
	logger         log.Logger
}

type deviceInfo struct {
	slotNumber string
	state      string
}

func NewDiskStatusCollector(promLog log.Logger) *DiskStatusCollector {
	return &DiskStatusCollector{
		hostDiskStatus: prometheus.NewDesc(
			"host_disk_status",
			"The host disk status check (state=UBUnsp) instructions disk abnormal",
			[]string{"slotNumber", "state"},
			nil,
		),
		hostRaidStatus: prometheus.NewDesc(
			"host_raid_status",
			"The host raid status check(0=abnormal, 1=normal)",
			[]string{"raidCardName"},
			nil,
		),
		logger: promLog,
	}
}

func (d *DiskStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- d.hostDiskStatus
	ch <- d.hostRaidStatus
}

func (d *DiskStatusCollector) Collect(ch chan<- prometheus.Metric) {
	raidName, raidStatus, diskInfo := d.getDiskStatusInfo()
	level.Debug(d.logger).Log("msg", "Get raid card name", raidName, "Get raid current status", raidStatus)
	for _, va := range diskInfo {
		ch <- prometheus.MustNewConstMetric(
			d.hostDiskStatus,
			prometheus.GaugeValue,
			1,
			va.slotNumber,
			va.state,
		)
	}
	ch <- prometheus.MustNewConstMetric(
		d.hostRaidStatus,
		prometheus.GaugeValue,
		raidStatus,
		raidName,
	)
	level.Info(d.logger).Log("msg", "collectd disk status success")
}

func (d *DiskStatusCollector) getDiskStatusInfo() (string, float64, map[string]deviceInfo) {
	args := []string{
		"/c0",
		"show",
		"nolog",
	}
	results, err := Execute("storcli64", args...)
	level.Debug(d.logger).Log("msg", "Exec command storcli64", "command args", args)
	if err != nil {
		level.Error(d.logger).Log("msg", "Exec command storcli64 failed")
	}
	raidName, raidStatus, diskInfo := parseDiskInfo(results)
	return raidName, raidStatus, diskInfo
}

func parseDiskInfo(results []byte) (string, float64, map[string]deviceInfo) {
	var diskInfo map[string]deviceInfo
	var devices []string
	record := 1
	parse := false

	var raidStatus float64
	var raidName string

	file, err := ioutil.TempFile("/tmp/", "*")
	if err != nil {
		return "", 0, nil
	}
	defer os.Remove(file.Name())
	if _, err = file.Write(results); err != nil {
		return "", 0, nil
	}
	f, err := os.Open(file.Name())
	if err != nil {
		return "", 0, nil
	}
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.ContainsAny(line, "=") {
			element := strings.Split(line, "=")
			if strings.HasPrefix(element[0], "Status") {
				status := strings.TrimSpace(element[1])
				if status == "Success" {
					raidStatus = 1
				} else {
					raidStatus = 0
				}
			}
			if strings.HasPrefix(element[0], "Product Name") {
				element := strings.Split(line, "=")
				name := strings.TrimSpace(element[1])
				raidName = strings.TrimSpace(name)
			}
		}
		if strings.HasPrefix(line, "PD LIST") {
			parse = true
			record = 1
			devices = []string{}
		}
		if strings.HasPrefix(line, "------------") {
			record += 1
		} else if record == 3 {
			devices = append(devices, line)
		} else if record == 4 {
			if parse {
				diskInfo = parseDevice(devices)
			}
			parse = false
		}
	}
	return raidName, raidStatus, diskInfo
}

func parseDevice(devices []string) map[string]deviceInfo {
	deviceResult := make(map[string]deviceInfo)
	deInfo := deviceInfo{}
	for _, device := range devices {
		device = deleteExtraSpace(device)
		element := strings.Split(device, " ")
		if len(element) < 13 {
			continue
		}
		slotNumber := strings.TrimSpace(strings.Split(element[0], ":")[1])
		deInfo.slotNumber = slotNumber
		deInfo.state = strings.TrimSpace(element[2])
		deviceResult[slotNumber] = deInfo
	}
	return deviceResult
}

func deleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "  ", " ", -1)
	regStr := "\\s{2,}"
	reg, _ := regexp.Compile(regStr)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spcIndex := reg.FindStringIndex(string(s2))
	for len(spcIndex) > 0 {
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...)
		spcIndex = reg.FindStringIndex(string(s2))
	}
	return string(s2)
}
