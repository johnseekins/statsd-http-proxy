package proxy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InjectiveLabs/statsd-http-proxy/proxy/routehandler"
	"github.com/InjectiveLabs/statsd-http-proxy/proxy/router"
	"github.com/InjectiveLabs/statsd-http-proxy/proxy/statsdclient"
)

// Server is a proxy server between HTTP REST API and UDP Connection to StatsD
type Server struct {
	httpAddress string
	httpServer  *http.Server
	tlsCert     string
	tlsKey      string
}

// NewServer creates new instance of StatsD HTTP Proxy
func NewServer(
	httpHost string,
	httpPort int,
	httpReadTimeout int,
	httpWriteTimeout int,
	httpIdleTimeout int,
	statsdHost string,
	statsdPort int,
	tlsCert string,
	tlsKey string,
	metricPrefix string,
	tokenSecret string,
	verbose bool,
) *Server {
	// configure logging
	var logOutput io.Writer
	if verbose == true {
		logOutput = os.Stderr
	} else {
		logOutput = ioutil.Discard
	}

	log.SetOutput(logOutput)

	logger := log.New(logOutput, "", log.LstdFlags)

	// prepare metric prefix
	if metricPrefix != "" && (metricPrefix)[len(metricPrefix)-1:] != "." {
		metricPrefix = metricPrefix + "."
	}

	// create StatsD Client
	statsdclient.Init(fmt.Sprintf("%s:%d", statsdHost, statsdPort), metricPrefix)

	// build route handler
	routeHandler := routehandler.NewRouteHandler(
		statsdclient.Client(),
	)

	// build router
	httpServerHandler := router.NewHTTPRouter(routeHandler, tokenSecret)

	// get HTTP server address to bind
	httpAddress := fmt.Sprintf("%s:%d", httpHost, httpPort)

	// create http server
	httpServer := &http.Server{
		Addr:           httpAddress,
		Handler:        httpServerHandler,
		ErrorLog:       logger,
		ReadTimeout:    time.Duration(httpReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(httpWriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(httpIdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 11,
	}

	statsdHTTPProxyServer := Server{
		httpAddress,
		httpServer,
		tlsCert,
		tlsKey,
	}

	return &statsdHTTPProxyServer
}

// Listen starts listening HTTP connections
func (proxyServer *Server) Listen() {
	// prepare for gracefull shutdown
	gracefullStopSignalHandler := make(chan os.Signal, 1)
	signal.Notify(gracefullStopSignalHandler, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start HTTP/HTTPS proxy to StatsD
	go func() {
		log.Printf("Starting HTTP server at %s", proxyServer.httpAddress)

		// open StatsD connection
		defer statsdclient.Close()

		// open HTTP connection
		var err error
		if len(proxyServer.tlsCert) > 0 && len(proxyServer.tlsKey) > 0 {
			err = proxyServer.httpServer.ListenAndServeTLS(proxyServer.tlsCert, proxyServer.tlsKey)
		} else {
			err = proxyServer.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			log.Fatal("Can not start HTTP server")
		}
	}()

	<-gracefullStopSignalHandler

	// Gracefull shutdown
	log.Printf("Stopping HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := proxyServer.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP Server Shutdown Failed:%+v", err)
	}

	log.Printf("HTTP server stopped successfully")
}
