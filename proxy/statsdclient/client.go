package statsdclient

import (
	"log"
	"sync"

	"github.com/alexcesaro/statsd"
)

var client Statter
var clientMux = new(sync.RWMutex)

func Close() {
	clientMux.RLock()
	defer clientMux.RUnlock()
	if client == nil {
		return
	}
	client.Close()
}

func Client() Statter {
	clientMux.RLock()
	defer clientMux.RUnlock()

	return client
}

func Init(addr string, prefix string) {
	statter, err := statsd.New(
		statsd.Address(addr),
		statsd.Prefix(prefix),
		statsd.TagsFormat(statsd.InfluxDB),
		statsd.ErrorHandler(errHandler),
	)
	if err != nil {
		log.Fatalln("[WARN] statsd client init error", err)
	}

	clientMux.Lock()
	client = statter
	clientMux.Unlock()
}

func errHandler(err error) {
	log.Println("[WARN] statsd client error", err)
}
