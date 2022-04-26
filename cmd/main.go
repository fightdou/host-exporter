package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/digineo/go-ping"
	mon "github.com/digineo/go-ping/monitor"
	"github.com/fightdou/host-exporter/config"
	"github.com/fightdou/host-exporter/pkg"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"

	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	configFile    = kingpin.Flag("config.path", "Path to config file").Default("/opt/config.yml").String()
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics.").Default(":9490").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metric").Default("/metrics").String()
	pingInterval  = kingpin.Flag("ping.interval", "Interval for ICMP echo requests").Default("5s").Duration()
	pingTimeout   = kingpin.Flag("ping.timeout", "Timeout for ICMP echo request").Default("4s").Duration()
	pingSize      = kingpin.Flag("ping.size", "Payload size for ICMP echo requests").Default("56").Uint16()
	historySize   = kingpin.Flag("ping.history-size", "Number of results to remember per target").Default("10").Int()
	targets       = kingpin.Arg("targets", "A list of targets to ping").Strings()
	ipmiTimeout   = kingpin.Flag("ipmi.timeout", "Timeout for ICMP echo request").Default("3s").Duration()
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print("ipmi_exporter"))
	kingpin.Parse()
	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting ipmi_exporter", "version", version.Info())

	cfg, err := loadConfig()
	if err != nil {
		kingpin.FatalUsage("could not load config.path: %v", err)
	}

	if cfg.Ping.History < 1 {
		kingpin.FatalUsage("ping.history-size must be greater than 0")
	}

	if cfg.Ping.Size > 65500 {
		kingpin.FatalUsage("ping.size must be between 0 and 65500")
	}

	if len(cfg.Targets) == 0 {
		kingpin.FatalUsage("No targets specified")
	}

	m := startMonitor(cfg, logger)

	cpu := pkg.NewCpuCollector(logger, *ipmiTimeout)
	disk := pkg.NewDiskStatusCollector(logger)
	nic := pkg.NewNicOnline(logger)
	net := pkg.NewNetPing(logger, m)
	diskIo := pkg.NewDiskIOUtil(logger)
	prometheus.MustRegister(cpu)
	prometheus.MustRegister(disk)
	prometheus.MustRegister(nic)
	prometheus.MustRegister(net)
	prometheus.MustRegister(diskIo)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>Device Exporter</title></head>
             <body>
             <h1>Device Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func startMonitor(cfg *config.Config, promLog log.Logger) *mon.Monitor {
	var bind4 string
	if ln, err := net.Listen("tcp4", "127.0.0.1:0"); err == nil {
		// ipv4 enabled
		ln.Close()
		bind4 = "0.0.0.0"
	}
	pinger, _ := ping.New(bind4, "")
	if pinger.PayloadSize() != cfg.Ping.Size {
		pinger.SetPayloadSize(cfg.Ping.Size)
	}

	monitor := mon.New(pinger,
		cfg.Ping.Interval.Duration(),
		cfg.Ping.Timeout.Duration())
	monitor.HistorySize = cfg.Ping.History

	targets := make([]*pkg.Target, len(cfg.Targets))

	for i, host := range cfg.Targets {
		t := &pkg.Target{
			Host:  host,
			Delay: time.Duration(10*i) * time.Millisecond,
		}
		targets[i] = t
		err := t.AddOrUpdateMonitor(monitor)
		if err != nil {
			level.Error(promLog).Log("msg", err)
		}
	}
	return monitor
}

func loadConfig() (*config.Config, error) {
	if *configFile == "" {
		cfg := config.Config{}
		addFlagToConfig(&cfg)
		return &cfg, nil
	}
	f, err := os.Open(*configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load config file: %w", err)
	}
	defer f.Close()

	cfg, err := config.FromYAML(f)
	if err == nil {
		addFlagToConfig(cfg)
	}
	return cfg, err
}

func addFlagToConfig(cfg *config.Config) {
	if len(cfg.Targets) == 0 {
		cfg.Targets = *targets
	}
	if cfg.Ping.History == 0 {
		cfg.Ping.History = *historySize
	}
	if cfg.Ping.Interval == 0 {
		cfg.Ping.Interval.Set(*pingInterval)
	}
	if cfg.Ping.Timeout == 0 {
		cfg.Ping.Timeout.Set(*pingTimeout)
	}
	if cfg.Ping.Size == 0 {
		cfg.Ping.Size = *pingSize
	}
}
