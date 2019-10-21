package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

// StatsD connection params
const defaultStatsDHost = "127.0.0.1"
const defaultStatsDPort = 8125

// declare command line options
var httpHost = flag.String("http-host", defaultHTTPHost, "HTTP Host")
var httpPort = flag.Int("http-port", defaultHTTPPort, "HTTP Port")
var tlsCert = flag.String("tls-cert", "", "TLS certificate to enable HTTPS")
var tlsKey = flag.String("tls-key", "", "TLS private key  to enable HTTPS")
var statsdHost = flag.String("statsd-host", defaultStatsDHost, "StatsD Host")
var statsdPort = flag.Int("statsd-port", defaultStatsDPort, "StatsD Port")
var metricPrefix = flag.String("metric-prefix", "", "Prefix of metric name")
var tokenSecret = flag.String("jwt-secret", "", "Secret to encrypt JWT")
var verbose = flag.Bool("verbose", false, "Verbose")
var version = flag.Bool("version", false, "Show version")

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

	// configure verbosity of logging
	if *verbose == true {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	// prepare metric prefix
	if *metricPrefix != "" && (*metricPrefix)[len(*metricPrefix)-1:] != "." {
		*metricPrefix = *metricPrefix + "."
	}

	proxyServer := proxy.NewServer(
		*httpHost,
		*httpPort,
		*statsdHost,
		*statsdPort,
		*tlsCert,
		*tlsKey,
		*tokenSecret,
		*metricPrefix,
	)

	proxyServer.Listen()
}
