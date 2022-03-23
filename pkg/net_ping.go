package pkg

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	mon "github.com/digineo/go-ping/monitor"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type ipVersion uint8

type Target struct {
	Host      string
	Addresses []net.IPAddr
	Delay     time.Duration
	resolver  *net.Resolver
	mutex     sync.Mutex
}

const (
	ipv4 ipVersion = 4
	ipv6 ipVersion = 6
)

type NetPing struct {
	hostNetPingLoss *prometheus.Desc
	hostNetConn     *prometheus.Desc
	logger          log.Logger
	monitor         *mon.Monitor
	metrics         map[string]*mon.Metrics
}

func NewNetPing(promLog log.Logger, mon *mon.Monitor) *NetPing {
	return &NetPing{
		hostNetPingLoss: prometheus.NewDesc(
			"host_net_ping_loss_percent",
			"The host network packet loss rate",
			[]string{"target", "ip", "ip_version"},
			nil,
		),
		hostNetConn: prometheus.NewDesc(
			"host_net_conn_alive",
			"The host network Whether can reach",
			[]string{"target", "ip", "ip_version"},
			nil,
		),
		logger:  promLog,
		monitor: mon,
	}
}

func (n *NetPing) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.hostNetConn
	ch <- n.hostNetPingLoss
}

func (n *NetPing) Collect(ch chan<- prometheus.Metric) {
	if m := n.monitor.Export(); len(m) > 0 {
		n.metrics = m
	}

	for target, metrics := range n.metrics {
		l := strings.SplitN(target, " ", 3)
		loss := float64(metrics.PacketsLost) / float64(metrics.PacketsSent)
		res := 0
		isActive := (metrics.PacketsSent - metrics.PacketsLost) > 0
		if isActive {
			res = 1
		}
		ch <- prometheus.MustNewConstMetric(n.hostNetPingLoss, prometheus.GaugeValue, loss, l...)
		ch <- prometheus.MustNewConstMetric(n.hostNetConn, prometheus.GaugeValue, float64(res), l...)
	}
}

func (t *Target) AddOrUpdateMonitor(monitor *mon.Monitor) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	addrs, err := t.resolver.LookupIPAddr(context.Background(), t.Host)
	if err != nil {
		return fmt.Errorf("error resolving target: %w", err)
	}

	for _, addr := range addrs {
		err := t.addIfNew(addr, monitor)
		if err != nil {
			return err
		}
	}

	t.cleanUp(addrs, monitor)
	t.Addresses = addrs
	return nil
}

func (t *Target) addIfNew(addr net.IPAddr, monitor *mon.Monitor) error {
	if isIPAddrInSlice(addr, t.Addresses) {
		return nil
	}

	return t.add(addr, monitor)
}

func (t *Target) cleanUp(addr []net.IPAddr, monitor *mon.Monitor) {
	for _, o := range t.Addresses {
		if !isIPAddrInSlice(o, addr) {
			name := t.nameForIP(o)
			monitor.RemoveTarget(name)
		}
	}
}

func (t *Target) add(addr net.IPAddr, monitor *mon.Monitor) error {
	name := t.nameForIP(addr)
	return monitor.AddTargetDelayed(name, addr, t.Delay)
}

func (t *Target) nameForIP(addr net.IPAddr) string {
	return fmt.Sprintf("%s %s %s", t.Host, addr.IP, getIPVersion(addr))
}

func isIPAddrInSlice(ipa net.IPAddr, slice []net.IPAddr) bool {
	for _, x := range slice {
		if x.IP.Equal(ipa.IP) {
			return true
		}
	}

	return false
}

// getIPVersion returns the version of IP protocol used for a given address
func getIPVersion(addr net.IPAddr) ipVersion {
	if addr.IP.To4() == nil {
		return ipv6
	}

	return ipv4
}

// String converts ipVersion to a string represention of the IP version used (i.e. "4" or "6")
func (ipv ipVersion) String() string {
	return strconv.Itoa(int(ipv))
}
