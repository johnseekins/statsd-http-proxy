package statsdclient

type Statter interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Unique(bucket string, value string)
	Close()
}
