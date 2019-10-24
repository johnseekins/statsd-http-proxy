package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoMetric/statsd-http-proxy/proxy"
)

// Version is a current git commit hash and tag
// Injected by compilation flag
var Version = "Unknown"

// BuildNumber is a current commit hash
// Injected by compilation flag
var BuildNumber = "Unknown"

// BuildDate is a date of build
// Injected by compilation flag
var BuildDate = "Unknown"

// HTTP connection params
const defaultHTTPHost = "127.0.0.1"
const defaultHTTPPort = 8825
const defaultHTTPReadTimeout = 1
const defaultHTTPWriteTimeout = 1
const defaultHTTPIdleTimeout = 1

// StatsD connection params
const defaultStatsDHost = "127.0.0.1"
const defaultStatsDPort = 8125

// declare command line options
var httpHost = flag.String("http-host", defaultHTTPHost, "HTTP Host")
var httpPort = flag.Int("http-port", defaultHTTPPort, "HTTP Port")
var httpReadTimeout = flag.Int("http-timeout-read", defaultHTTPReadTimeout, "The maximum duration in seconds for reading the entire request, including the body")
var httpWriteTimeout = flag.Int("http-timeout-write", defaultHTTPWriteTimeout, "The maximum duration in seconds before timing out writes of the respons")
var httpIdleTimeout = flag.Int("http-timeout-idle", defaultHTTPIdleTimeout, "The maximum amount of time in seconds to wait for the next request when keep-alives are enabled")
var tlsCert = flag.String("tls-cert", "", "TLS certificate to enable HTTPS")
var tlsKey = flag.String("tls-key", "", "TLS private key  to enable HTTPS")
var statsdHost = flag.String("statsd-host", defaultStatsDHost, "StatsD Host")
var statsdPort = flag.Int("statsd-port", defaultStatsDPort, "StatsD Port")
var metricPrefix = flag.String("metric-prefix", "", "Prefix of metric name")
var tokenSecret = flag.String("jwt-secret", "", "Secret to encrypt JWT")
var verbose = flag.Bool("verbose", false, "Verbose")
var version = flag.Bool("version", false, "Show version")
var httpRouterName = flag.String("http-router-name", "GorillaMux", "Type of HTTP router. Allowed values are GorillaMux and HttpRouter. Do not use in production.")

func main() {
	// get flags
	flag.Parse()

	// show version and exit
	if *version == true {
		fmt.Printf(
			"StatsD HTTP Proxy v.%s, build %s from %s\n",
			Version,
			BuildNumber,
			BuildDate,
		)
		os.Exit(0)
	}

	// prepare metric prefix
	if *metricPrefix != "" && (*metricPrefix)[len(*metricPrefix)-1:] != "." {
		*metricPrefix = *metricPrefix + "."
	}

	proxyServer := proxy.NewServer(
		*httpHost,
		*httpPort,
		*httpReadTimeout,
		*httpWriteTimeout,
		*httpIdleTimeout,
		*statsdHost,
		*statsdPort,
		*tlsCert,
		*tlsKey,
		*metricPrefix,
		*tokenSecret,
		*verbose,
		*httpRouterName,
	)

	proxyServer.Listen()
}
