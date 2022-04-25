package statsd

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const metricTypeCount = "c"
const metricTypeGauge = "g"
const metricTypeTiming = "ms"
const metricTypeSet = "s"

// The Client type
type Client struct {
	host                string
	port                int
	conn                net.Conn   // UDP connection to StatsD server
	rand                *rand.Rand // rand generator to skip messages by sample rate
	buffer              []string   // array of messages to send
	bufferLength        int
	bufferLock          sync.RWMutex // mutex to lock buffer of keys
	bufferFlushTicker   *time.Ticker
	bufferFlushInterval *time.Duration
	metricPrefix        string // used to prefix all metrics before send
}

// NewClient creates new StatsD client with disabled buffer.
func NewClient(host string, port int) *Client {
	client := Client{
		host:         host,
		port:         port,
		rand:         rand.New(rand.NewSource(time.Now().Unix())),
		metricPrefix: "",
	}

	return &client
}

// NewBufferedClient creates new StatsD client with enabled buffer.
// If flush interval not defined, manual call of Flush() required to send metrics to StatsD server.
func NewBufferedClient(host string, port int) *Client {
	client := Client{
		host:         host,
		port:         port,
		rand:         rand.New(rand.NewSource(time.Now().Unix())),
		buffer:       make([]string, 0),
		bufferLength: 0,
		metricPrefix: "",
	}

	return &client
}

func (client *Client) SetFlushInterval(flushInterval time.Duration) {
	client.bufferFlushInterval = &flushInterval
}

func (client *Client) flushOnTick() {
	for ; true; <-client.bufferFlushTicker.C {
		client.Flush()
	}
}

// SetPrefix adds prefix to all keys
func (client *Client) SetPrefix(metricPrefix string) {
	if metricPrefix != "" && (metricPrefix)[len(metricPrefix)-1:] != "." {
		metricPrefix = metricPrefix + "."
	}

	client.metricPrefix = metricPrefix
}

// Open UDP connection to statsd server
func (client *Client) Open() {
	// start UDP connection to StatsD Server
	connectionString := fmt.Sprintf("%s:%d", client.host, client.port)

	conn, err := net.Dial("udp", connectionString)

	if err != nil {
		log.Println(err)
	}

	client.conn = conn

	// start flush ticker if buffered client
	if client.buffer != nil && client.bufferFlushInterval != nil {
		client.bufferFlushTicker = time.NewTicker(*client.bufferFlushInterval)

		go client.flushOnTick()
	}
}

// Close UDP connection to statsd server
func (client *Client) Close() {
	// send buffer and stop flushing
	if client.buffer != nil {
		// flush metrics in buffer
		client.Flush()

		// stop flush ticker
		if client.bufferFlushTicker != nil {
			client.bufferFlushTicker.Stop()
		}
	}

	// close UDP connection
	client.conn.Close()
	client.conn = nil
}

// Timing track in milliseconds with sampling
func (client *Client) Timing(key string, time int64, sampleRate float32) {
	metricValue := fmt.Sprintf("%d|%s", time, metricTypeTiming)
	if sampleRate < 1 {
		if client.isSendAcceptedBySampleRate(sampleRate) {
			metricValue = fmt.Sprintf("%s|@%g", metricValue, sampleRate)
		} else {
			return
		}
	}

	client.addToBuffer(key, metricValue)
}

// Count tack
func (client *Client) Count(key string, value int, sampleRate float32) {
	metricValue := fmt.Sprintf("%d|%s", value, metricTypeCount)
	if sampleRate < 1 {
		if client.isSendAcceptedBySampleRate(sampleRate) {
			metricValue = fmt.Sprintf("%s|@%g", metricValue, sampleRate)
		} else {
			return
		}
	}

	client.addToBuffer(key, metricValue)
}

// Gauge track
// To set a gauge to a negative number you need first set it to 0, because negative value interprets as negative shift.
func (client *Client) Gauge(key string, value int) {
	metricValue := fmt.Sprintf("%d|%s", value, metricTypeGauge)
	client.addToBuffer(key, metricValue)
}

// GaugeShift decrease previously set value if negative value passed, and increase if positive.
func (client *Client) GaugeShift(key string, value int) {
	metricValue := fmt.Sprintf("%+d|%s", value, metricTypeGauge)
	client.addToBuffer(key, metricValue)
}

// Set tracking
func (client *Client) Set(key string, value int) {
	metricValue := fmt.Sprintf("%d|%s", value, metricTypeSet)
	client.addToBuffer(key, metricValue)
}

// add to buffer and flush if auto flush enabled
func (client *Client) addToBuffer(key string, metricValue string) {
	// build metric
	metric := fmt.Sprintf("%s:%s", client.metricPrefix+key, metricValue)

	// flush
	if client.buffer == nil {
		// send metric now
		client.send(metric)
	} else {
		client.bufferLock.Lock()

		newBufferLength := client.bufferLength + len(metric)
		isBufferOverflow := newBufferLength > 1400
		if !isBufferOverflow {
			client.bufferLength = client.bufferLength + newBufferLength
			client.buffer = append(client.buffer, metric)
		}

		client.bufferLock.Unlock()

		if isBufferOverflow {
			client.send(metric)
			client.Flush()
		}
	}
}

// Check if acceptable by sample rate
func (client *Client) isSendAcceptedBySampleRate(sampleRate float32) bool {
	if sampleRate >= 1 {
		return true
	}
	randomNumber := client.rand.Float32()
	return randomNumber <= sampleRate
}

// Flush buffer to statsd daemon by UDP when buffer disabled
func (client *Client) Flush() error {
	// check if buffer enabled
	if client.buffer == nil {
		return errors.New("Invalid call of flush in unbuffered mode")
	}

	// check if buffer has metrics
	if len(client.buffer) == 0 {
		return nil
	}

	// lock
	client.bufferLock.Lock()

	// build packet
	metricPacket := strings.Join(client.buffer, "\n")

	// clear key buffer
	client.buffer = make([]string, 0)
	client.bufferLength = 0

	// lock
	client.bufferLock.Unlock()

	// send packet
	client.send(metricPacket)

	return nil
}

// Send StatsD packet
func (client *Client) send(metricPacket string) {
	// send metric packet
	_, err := fmt.Fprintf(client.conn, metricPacket)
	if err != nil {
		log.Println(err)
	}
}
