package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		errorMsg string
		want     bool
	}{
		{
			name:     "network connection error",
			errorMsg: "failed to get http://localhost:8080/stub_status: connection refused",
			want:     true,
		},
		{
			name:     "network timeout error",
			errorMsg: "failed to get http://localhost:8080/stub_status: timeout",
			want:     true,
		},
		{
			name:     "HTTP error",
			errorMsg: "expected 200 response, got 404",
			want:     false,
		},
		{
			name:     "parse error",
			errorMsg: "failed to parse response body",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNetworkError(tt.errorMsg); got != tt.want {
				t.Errorf("isNetworkError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsHTTPError(t *testing.T) {
	tests := []struct {
		name     string
		errorMsg string
		want     bool
	}{
		{
			name:     "HTTP 404 error",
			errorMsg: "expected 200 response, got 404",
			want:     true,
		},
		{
			name:     "HTTP 500 error",
			errorMsg: "expected 200 response, got 500",
			want:     true,
		},
		{
			name:     "network error",
			errorMsg: "failed to get http://localhost:8080/stub_status: connection refused",
			want:     false,
		},
		{
			name:     "parse error",
			errorMsg: "failed to parse response body",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHTTPError(tt.errorMsg); got != tt.want {
				t.Errorf("isHTTPError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewScrapeSuccessMetric(t *testing.T) {
	metric := newScrapeSuccessMetric("nginx", map[string]string{"job": "nginx"})

	if metric == nil {
		t.Error("newScrapeSuccessMetric() returned nil")
	}

	desc := metric.Desc().String()
	if !strings.Contains(desc, "nginx_scrape_success") {
		t.Errorf("metric description should contain 'nginx_scrape_success', got: %s", desc)
	}
}

func TestNewScrapeDurationMetric(t *testing.T) {
	metric := newScrapeDurationMetric("nginx", map[string]string{"job": "nginx"})

	if metric == nil {
		t.Error("newScrapeDurationMetric() returned nil")
	}

	desc := metric.Desc().String()
	if !strings.Contains(desc, "nginx_scrape_duration_seconds") {
		t.Errorf("metric description should contain 'nginx_scrape_duration_seconds', got: %s", desc)
	}
}

func TestNewScrapeErrorsTotalMetric(t *testing.T) {
	metric := newScrapeErrorsTotalMetric("nginx", map[string]string{"job": "nginx"})

	if metric == nil {
		t.Error("newScrapeErrorsTotalMetric() returned nil")
	}

	ch := make(chan *prometheus.Desc, 10)
	metric.Describe(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}

	if count == 0 {
		t.Error("metric should have descriptions")
	}
}
