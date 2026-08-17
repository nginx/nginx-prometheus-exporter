package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	nginxUp   = 1
	nginxDown = 0
)

func newGlobalMetric(namespace string, metricName string, docString string, constLabels map[string]string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, nil, constLabels)
}

func newUpMetric(namespace string, constLabels map[string]string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "up",
		Help:        "Status of the last metric scrape",
		ConstLabels: constLabels,
	})
}

func newScrapeSuccessMetric(namespace string, constLabels map[string]string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "scrape_success",
		Help:        "Whether the last scrape of NGINX metrics was successful",
		ConstLabels: constLabels,
	})
}

func newScrapeDurationMetric(namespace string, constLabels map[string]string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "scrape_duration_seconds",
		Help:        "Duration of the last scrape in seconds",
		ConstLabels: constLabels,
	})
}

func newScrapeErrorsTotalMetric(namespace string, constLabels map[string]string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Name:        "scrape_errors_total",
		Help:        "Total number of scrape errors by type",
		ConstLabels: constLabels,
	}, []string{"type"})
}

// MergeLabels merges two maps of labels.
func MergeLabels(a map[string]string, b map[string]string) map[string]string {
	c := make(map[string]string)

	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}

	return c
}

func isNetworkError(errorMsg string) bool {
	return strings.Contains(errorMsg, "failed to get") ||
		strings.Contains(errorMsg, "connection") ||
		strings.Contains(errorMsg, "timeout") ||
		strings.Contains(errorMsg, "refused")
}

func isHTTPError(errorMsg string) bool {
	return strings.Contains(errorMsg, "expected 200 response")
}
