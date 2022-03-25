package pkg

import (
	"context"
	"fmt"
	"github.com/fightdou/host-exporter/config"
	"github.com/go-kit/log/level"
	"net"
	"strings"
	"time"

	mon "github.com/digineo/go-ping/monitor"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Target struct {
	Host     string
	Delay    time.Duration
	resolver *net.Resolver
	cfg      *config.Config
}

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
			"The host network packet loss percent",
			[]string{"target"},
			nil,
		),
		hostNetConn: prometheus.NewDesc(
			"host_net_target_conn_status",
			"The host network target Whether can reach(0=abnormal, 1=normal)",
			[]string{"target"},
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
	level.Info(n.logger).Log("msg", "collectd net conn status success")
}

func (t *Target) AddOrUpdateMonitor(monitor *mon.Monitor) error {
	addrs, err := t.resolver.LookupIPAddr(context.Background(), t.Host)
	if err != nil {
		return fmt.Errorf("error resolving target: %w", err)
	}

	for _, addr := range addrs {
		monitor.AddTargetDelayed(t.Host, addr, t.Delay)
		if err != nil {
			return err
		}
	}
	return nil
}
