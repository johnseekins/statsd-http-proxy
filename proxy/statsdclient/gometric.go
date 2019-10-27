package statsdclient

import GoMetricStatsdClient "github.com/GoMetric/go-statsd-client"

func NewGoMetricClient(
	statsdHost string,
	statsdPort int,
) StatsdClientInterface {
	return GoMetricStatsdClient.NewClient(statsdHost, statsdPort)
}
