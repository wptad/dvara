package main

import (
	"log"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

func NewDataDogStatsDClient(address string) DatadogStatsClient {
	c, err := statsd.New(address)
	if err != nil {
		log.Fatal(err)
	}
	return DatadogStatsClient{c}
}

type DatadogStatsClient struct {
	client *statsd.Client
}

func (c DatadogStatsClient) BumpAvg(key string, val float64) {
	// average can go up or down, so I gues gauge is best aproximate
	c.client.Gauge(sanitizeStatsKey(key), val, nil, 1)
}

func (c DatadogStatsClient) BumpHistogram(key string, val float64) {
	c.client.Histogram(sanitizeStatsKey(key), val, nil, 1)
}

func (c DatadogStatsClient) BumpSum(key string, val float64) {
	// Sum can go only up, so I gues Count is best aproximate, I'm not
	// sure how lossy is float to int conversion here
	// code grep indicates that method is usually called with value of 1
	c.client.Count(sanitizeStatsKey(key), int64(val), nil, 1)
}

func (c DatadogStatsClient) BumpTime(key string) interface {
	End()
} {
	return timeEnd{c, sanitizeStatsKey(key), time.Now()}
}

type timeEnd struct {
	dataDogStatsClient DatadogStatsClient
	key                string
	eventStartTime     time.Time
}

func (n timeEnd) End() {
	// Graphite default precision is millisecond I think, should switch later
	// to millisecond I guess
	n.dataDogStatsClient.client.Gauge(sanitizeStatsKey(n.key), float64(time.Since(n.eventStartTime).Nanoseconds()), nil, 1)
}

func sanitizeStatsKey(statKey string) string {
	// statsd format uses ":" as separator, so keys must not contain this param
	// no particular reason to use underscore, just since it is easy to read
	return strings.Replace(statKey, ":", "_", -1)
}
