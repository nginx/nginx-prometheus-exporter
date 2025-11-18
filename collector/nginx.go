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
	upMetric             prometheus.Gauge
	scrapeSuccessMetric  prometheus.Gauge
	scrapeDurationMetric prometheus.Gauge
	scrapeErrorsTotal    *prometheus.CounterVec
	logger               *slog.Logger
	nginxClient          *client.NginxClient
	metrics              map[string]*prometheus.Desc
	mutex                sync.Mutex
}

// NewNginxCollector creates an NginxCollector.
func NewNginxCollector(nginxClient *client.NginxClient, namespace string, constLabels map[string]string, logger *slog.Logger) *NginxCollector {
	return &NginxCollector{
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
		upMetric:             newUpMetric(namespace, constLabels),
		scrapeSuccessMetric:  newScrapeSuccessMetric(namespace, constLabels),
		scrapeDurationMetric: newScrapeDurationMetric(namespace, constLabels),
		scrapeErrorsTotal:    newScrapeErrorsTotalMetric(namespace, constLabels),
	}
}

// Describe sends the super-set of all possible descriptors of NGINX metrics
// to the provided channel.
func (c *NginxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upMetric.Desc()
	ch <- c.scrapeSuccessMetric.Desc()
	ch <- c.scrapeDurationMetric.Desc()
	c.scrapeErrorsTotal.Describe(ch)

	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect fetches metrics from NGINX and sends them to the provided channel.
func (c *NginxCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	start := time.Now()
	stats, err := c.nginxClient.GetStubStats()
	duration := time.Since(start).Seconds()
	c.scrapeDurationMetric.Set(duration)
	ch <- c.scrapeDurationMetric

	if err != nil {
		c.handleScrapeError(ch, err)
		return
	}

	c.handleScrapeSuccess(ch, stats)
}

func (c *NginxCollector) handleScrapeError(ch chan<- prometheus.Metric, err error) {
	errorMsg := err.Error()
	var errorType string

	if isNetworkError(errorMsg) {
		c.upMetric.Set(nginxDown)
		errorType = "network"
	} else if isHTTPError(errorMsg) {
		c.upMetric.Set(nginxUp)
		errorType = "http"
	} else {
		c.upMetric.Set(nginxUp)
		errorType = "parse"
	}

	c.scrapeErrorsTotal.WithLabelValues(errorType).Inc()
	c.scrapeSuccessMetric.Set(0)

	ch <- c.upMetric
	ch <- c.scrapeSuccessMetric
	c.scrapeErrorsTotal.Collect(ch)

	c.logger.Error("error getting stats", "error", err.Error(), "type", errorType)
}

func (c *NginxCollector) handleScrapeSuccess(ch chan<- prometheus.Metric, stats *client.StubStats) {
	c.upMetric.Set(nginxUp)
	c.scrapeSuccessMetric.Set(1)

	ch <- c.upMetric
	ch <- c.scrapeSuccessMetric
	c.scrapeErrorsTotal.Collect(ch)

	ch <- prometheus.MustNewConstMetric(c.metrics["connections_active"],
		prometheus.GaugeValue, float64(stats.Connections.Active))
	ch <- prometheus.MustNewConstMetric(c.metrics["connections_accepted"],
		prometheus.CounterValue, float64(stats.Connections.Accepted))
	ch <- prometheus.MustNewConstMetric(c.metrics["connections_handled"],
		prometheus.CounterValue, float64(stats.Connections.Handled))
	ch <- prometheus.MustNewConstMetric(c.metrics["connections_reading"],
		prometheus.GaugeValue, float64(stats.Connections.Reading))
	ch <- prometheus.MustNewConstMetric(c.metrics["connections_writing"],
		prometheus.GaugeValue, float64(stats.Connections.Writing))
	ch <- prometheus.MustNewConstMetric(c.metrics["connections_waiting"],
		prometheus.GaugeValue, float64(stats.Connections.Waiting))
	ch <- prometheus.MustNewConstMetric(c.metrics["http_requests_total"],
		prometheus.CounterValue, float64(stats.Requests))
}
