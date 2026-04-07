package collector

import (
	"log/slog"
	"sync"
	"time"

	"github.com/nginx/nginx-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

// NginxCollector collects NGINX metrics. It implements prometheus.Collector interface.
type NginxCollector struct {
	upMetric    prometheus.Gauge
	logger      *slog.Logger
	nginxClient *client.NginxClient
	metrics     map[string]*prometheus.Desc
	mutex       sync.Mutex
	createdAt   time.Time
}

// NewNginxCollector creates an NginxCollector.
func NewNginxCollector(nginxClient *client.NginxClient, namespace string, constLabels map[string]string, logger *slog.Logger, ctSource CTSource) *NginxCollector {
	c := &NginxCollector{
		nginxClient: nginxClient,
		logger:      logger,
		metrics: map[string]*prometheus.Desc{
			"connections_active":   newGlobalMetric(namespace, "connections_active", "Active client connections", constLabels),
			"connections_accepted": newGlobalMetric(namespace, "connections_accepted", "Accepted client connections", constLabels),
			"connections_handled":  newGlobalMetric(namespace, "connections_handled", "Handled client connections", constLabels),
			"connections_reading":  newGlobalMetric(namespace, "connections_reading", "Connections where NGINX is reading the request header", constLabels),
			"connections_writing":  newGlobalMetric(namespace, "connections_writing", "Connections where NGINX is writing the response back to the client", constLabels),
			"connections_waiting":  newGlobalMetric(namespace, "connections_waiting", "Idle client connections", constLabels),
			"http_requests_total":  newGlobalMetric(namespace, "http_requests_total", "Total http requests", constLabels),
		},
		upMetric: newUpMetric(namespace, constLabels),
	}
	if ctSource == CTSourceProcess {
		c.createdAt = nowFunc()
	}
	return c
}

// Describe sends the super-set of all possible descriptors of NGINX metrics
// to the provided channel.
func (c *NginxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upMetric.Desc()

	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *NginxCollector) newGaugeMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, labelValues...)
}

func (c *NginxCollector) newCounterMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, labelValues ...string) {
	if c.createdAt.IsZero() {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, value, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetricWithCreatedTimestamp(desc, prometheus.CounterValue, value, c.createdAt, labelValues...)
	}
}

// Collect fetches metrics from NGINX and sends them to the provided channel.
func (c *NginxCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	stats, err := c.nginxClient.GetStubStats()
	if err != nil {
		c.upMetric.Set(nginxDown)
		ch <- c.upMetric
		c.logger.Error("error getting stats", "uri", c.nginxClient.GetAPIEndpoint(), "error", err)
		return
	}

	c.upMetric.Set(nginxUp)
	ch <- c.upMetric

	c.newGaugeMetric(ch, c.metrics["connections_active"],
		float64(stats.Connections.Active))
	c.newCounterMetric(ch, c.metrics["connections_accepted"],
		float64(stats.Connections.Accepted))
	c.newCounterMetric(ch, c.metrics["connections_handled"],
		float64(stats.Connections.Handled))
	c.newGaugeMetric(ch, c.metrics["connections_reading"],
		float64(stats.Connections.Reading))
	c.newGaugeMetric(ch, c.metrics["connections_writing"],
		float64(stats.Connections.Writing))
	c.newGaugeMetric(ch, c.metrics["connections_waiting"],
		float64(stats.Connections.Waiting))
	c.newCounterMetric(ch, c.metrics["http_requests_total"],
		float64(stats.Requests))
}
