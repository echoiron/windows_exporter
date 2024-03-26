//go:build windows

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus-community/windows_exporter/custom/router"
	"github.com/prometheus-community/windows_exporter/custom/util"
	winlog "github.com/prometheus-community/windows_exporter/pkg/log"
	"github.com/prometheus-community/windows_exporter/pkg/types"
	"github.com/prometheus-community/windows_exporter/pkg/utils"
	"github.com/prometheus-community/windows_exporter/pkg/wmi"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/prometheus-community/windows_exporter/pkg/collector"
	"github.com/prometheus-community/windows_exporter/pkg/config"
	"github.com/prometheus-community/windows_exporter/pkg/log/flag"
	"github.com/prometheus/common/version"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

// Same struct prometheus uses for their /version endpoint.
// Separate copy to avoid pulling all of prometheus as a dependency
type prometheusVersion struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

func main() {
	app := kingpin.New("windows_exporter", "A metrics collector for Windows.")
	var (
		configFile = app.Flag(
			"config.file",
			"YAML configuration file to use. Values set in this file will be overridden by CLI flags.",
		).String()
		insecure_skip_verify = app.Flag(
			"config.file.insecure-skip-verify",
			"Skip TLS verification in loading YAML configuration.",
		).Default("false").Bool()
		webConfig   = webflag.AddFlags(app, ":9182")
		metricsPath = app.Flag(
			"telemetry.path",
			"URL path for surfacing collected metrics.",
		).Default("/metrics").String()
		disableExporterMetrics = app.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Bool()
		maxRequests = app.Flag(
			"telemetry.max-requests",
			"Maximum number of concurrent requests. 0 to disable.",
		).Default("5").Int()
		enabledCollectors = app.Flag(
			"collectors.enabled",
			"Comma-separated list of collectors to use. Use '[defaults]' as a placeholder for all the collectors enabled by default.").
			Default(types.DefaultCollectors).String()
		printCollectors = app.Flag(
			"collectors.print",
			"If true, print available collectors and exit.",
		).Bool()
		timeoutMargin = app.Flag(
			"scrape.timeout-margin",
			"Seconds to subtract from the timeout allowed by the client. Tune to allow for overhead or high loads.",
		).Default("0.5").Float64()
	)

	winlogConfig := &winlog.Config{}
	flag.AddFlags(app, winlogConfig)

	app.Version(version.Print("windows_exporter"))
	app.HelpFlag.Short('h')

	// Initialize collectors before loading and parsing CLI arguments
	collectors := collector.NewWithFlags(app)

	// Load values from configuration file(s). Executable flags must first be parsed, in order
	// to load the specified file(s).
	kingpin.MustParse(app.Parse(os.Args[1:]))
	logger, err := winlog.New(winlogConfig)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	_ = level.Debug(logger).Log("msg", "Logging has Started")

	if *configFile != "" {
		resolver, err := config.NewResolver(*configFile, logger, *insecure_skip_verify)
		if err != nil {
			_ = level.Error(logger).Log("msg", "could not load config file", "err", err)
			os.Exit(1)
		}
		err = resolver.Bind(app, os.Args[1:])
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			os.Exit(1)
		}

		// NOTE: This is temporary fix for issue #1092, calling kingpin.Parse
		// twice makes slices flags duplicate its value, this clean up
		// the first parse before the second call.
		*webConfig.WebListenAddresses = (*webConfig.WebListenAddresses)[1:]

		// Parse flags once more to include those discovered in configuration file(s).
		kingpin.MustParse(app.Parse(os.Args[1:]))

		logger, err = winlog.New(winlogConfig)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			os.Exit(1)
		}
	}

	if *printCollectors {
		collectorNames := collector.Available()
		sort.Strings(collectorNames)

		fmt.Printf("Available collectors:\n")
		for _, n := range collectorNames {
			fmt.Printf(" - %s\n", n)
		}

		return
	}

	if err = wmi.InitWbem(logger); err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	enabledCollectorList := utils.ExpandEnabledCollectors(*enabledCollectors)
	collectors.Enable(enabledCollectorList)
	collectors.SetLogger(logger)

	// Initialize collectors before loading
	err = collectors.Build()
	if err != nil {
		_ = level.Error(logger).Log("msg", "Couldn't load collectors", "err", err)
		os.Exit(1)
	}
	err = collectors.SetPerfCounterQuery()
	if err != nil {
		_ = level.Error(logger).Log("msg", "Couldn't set performance counter query", "err", err)
		os.Exit(1)
	}

	if u, err := user.Current(); err != nil {
		_ = level.Warn(logger).Log("msg", "Unable to determine which user is running this exporter. More info: https://github.com/golang/go/issues/37348")
	} else {
		_ = level.Info(logger).Log("msg", fmt.Sprintf("Running as %v", u.Username))

		if strings.Contains(u.Username, "ContainerAdministrator") || strings.Contains(u.Username, "ContainerUser") {
			_ = level.Warn(logger).Log("msg", "Running as a preconfigured Windows Container user. This may mean you do not have Windows HostProcess containers configured correctly and some functionality will not work as expected.")
		}
	}

	_ = level.Info(logger).Log("msg", fmt.Sprintf("Enabled collectors: %v", strings.Join(enabledCollectorList, ", ")))

	r := gin.Default()
	// init validator
	util.RegValid()
	r.GET(*metricsPath, func(c *gin.Context) {
		withConcurrencyLimit(*maxRequests, collectors.BuildServeHTTP(*disableExporterMetrics, *timeoutMargin)).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":   version.Version,
			"revision":  version.Revision,
			"branch":    version.Branch,
			"buildUser": version.BuildUser,
			"buildDate": version.BuildDate,
			"goVersion": version.GoVersion,
		})
	})

	router.RegisterRoutes(r)
	runMode := SetEnvStr("GIN_MODE", "debug")
	_ = os.Setenv("GIN_MODE", runMode)

	if err = r.Run(*webConfig.WebListenAddresses...); err != nil {
		_ = level.Error(logger).Log("msg", "cannot start windows_exporter", "err", err)
		os.Exit(1)
	}
}

func SetEnvStr(env string, value string) string {
	if e := os.Getenv(env); e != "" {
		return e
	}
	return value
}

func withConcurrencyLimit(n int, next http.HandlerFunc) http.HandlerFunc {
	if n <= 0 {
		return next
	}

	sem := make(chan struct{}, n)
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Too many concurrent requests"))
			return
		}
		next(w, r)
	}
}
