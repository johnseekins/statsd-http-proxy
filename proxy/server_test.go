package proxy

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	proxyServer := NewServer(
		"127.0.0.1",
		8080,
		1,
		1,
		1,
		"127.0.0.1",
		8125,
		"tls.crt",
		"tls.key",
		"someMetricPrefix",
		"someTokenSecret",
		true,
		"GorillaMux",
	)

	require.Equal(t, "*proxy.Server", reflect.TypeOf(proxyServer).String())

	require.Equal(t, "127.0.0.1:8080", proxyServer.httpAddress)

	require.Equal(t, "*statsd.Client", reflect.TypeOf(proxyServer.statsdClient).String())

	require.Equal(t, "tls.crt", proxyServer.tlsCert)

	require.Equal(t, "tls.key", proxyServer.tlsKey)
}
