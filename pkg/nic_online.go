package pkg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/safchain/ethtool"
)

type NicOnline struct {
	hostNicOnline *prometheus.Desc
	logger        log.Logger
}

func NewNicOnline(promLog log.Logger) *NicOnline {
	return &NicOnline{
		hostNicOnline: prometheus.NewDesc(
			"host_nic_on_line",
			"The host nic online status",
			[]string{"interface"},
			nil,
		),
		logger: promLog,
	}
}

func (n *NicOnline) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.hostNicOnline
}

func (n *NicOnline) Collect(ch chan<- prometheus.Metric) {
	nicState, err := getNicStatus()
	if err != nil {
		level.Error(n.logger).Log("msg", "Get Nic status failed")
	}
	for k, v := range nicState {
		ch <- prometheus.MustNewConstMetric(
			n.hostNicOnline,
			prometheus.GaugeValue,
			v,
			k,
		)
	}
	level.Info(n.logger).Log("msg", "collectd nic online status success")
}

func getNicStatus() (map[string]float64, error) {
	nicStatus := map[string]float64{}
	nicName, err := getNicName()
	if err != nil {
		return nil, err
	}
	ethHandle, err := ethtool.NewEthtool()
	if err != nil {
		return nil, err
	}
	defer ethHandle.Close()
	for _, nic := range nicName {
		stats, err := ethHandle.LinkState(nic)
		if err != nil {
			return nil, err
		}
		nicStatus[nic] = float64(stats)
	}
	return nicStatus, nil
}

func getNicName() ([]string, error) {
	var pciNumber []string
	var nicNames []string
	res, err := Execute("/bin/sh", "-c", `lspci | grep -i Ethernet`)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(res), "\n") {
		pciInfo := strings.Split(line, " ")
		if pciInfo[0] == "" {
			continue
		}
		pciNumber = append(pciNumber, pciInfo[0])
	}

	for _, pci := range pciNumber {
		path := fmt.Sprintf("/sys/bus/pci/devices/0000:%s/net/*/mtu", pci)
		paths, fErr := filepath.Glob(path)
		if fErr != nil {
			return nil, err
		}
		dirName, _ := filepath.Split(paths[0])
		nic := strings.Split(dirName, "net")[1]
		nicName := strings.Trim(nic, "/")
		nicNames = append(nicNames, nicName)
	}
	return nicNames, nil
}
