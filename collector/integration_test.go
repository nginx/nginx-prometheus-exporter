package collector

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	plusclient "github.com/nginx/nginx-plus-go-client/v3/client"
	"github.com/nginx/nginx-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const stubStatusResponse = `Active connections: 1457
server accepts handled requests
 6717066 6717066 65844359
Reading: 1 Writing: 8 Waiting: 1448
`

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newOSSMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, stubStatusResponse)
	}))
}

func newPlusMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.RequestURI {
		case "/":
			fmt.Fprint(w, `[9]`)
		case "/9/":
			fmt.Fprint(w, `["nginx","processes","connections","slabs","http","resolvers","ssl"]`)
		case "/9/nginx":
			fmt.Fprint(w, `{"version":"1.25.0","build":"nginx-plus-r30","pid":123,"load_timestamp":"2025-06-15T12:00:00.000Z"}`)
		case "/9/connections":
			fmt.Fprint(w, `{"accepted":100,"dropped":5,"active":10,"idle":20}`)
		case "/9/http/requests":
			fmt.Fprint(w, `{"total":1000,"current":5}`)
		case "/9/ssl":
			fmt.Fprint(w, `{"handshakes":50,"handshakes_failed":2,"session_reuses":30}`)
		case "/9/processes":
			fmt.Fprint(w, `{}`)
		case "/9/slabs":
			fmt.Fprint(w, `{}`)
		case "/9/http/caches":
			fmt.Fprint(w, `{}`)
		case "/9/http/server_zones":
			fmt.Fprint(w, `{}`)
		case "/9/http/upstreams":
			fmt.Fprint(w, `{}`)
		case "/9/http/location_zones":
			fmt.Fprint(w, `{}`)
		case "/9/resolvers":
			fmt.Fprint(w, `{}`)
		case "/9/http/limit_reqs":
			fmt.Fprint(w, `{}`)
		case "/9/http/limit_conns":
			fmt.Fprint(w, `{}`)
		case "/9/workers":
			fmt.Fprint(w, `[]`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// gatherCreatedTimestamps calls Gather() on the registry and returns a map of
// metric family name → created timestamp for all counter metrics that have a
// non-nil CreatedTimestamp.
func gatherCreatedTimestamps(t *testing.T, registry *prometheus.Registry) map[string]time.Time {
	t.Helper()
	families, err := registry.Gather()
	if err != nil {
		t.Fatal(err)
	}

	result := make(map[string]time.Time)
	for _, mf := range families {
		if mf.GetType() != dto.MetricType_COUNTER {
			continue
		}
		for _, m := range mf.GetMetric() {
			ct := m.GetCounter().GetCreatedTimestamp()
			if ct != nil {
				result[mf.GetName()] = ct.AsTime()
				break // one per family is enough
			}
		}
	}
	return result
}

func TestCreatedTimestamp_NginxOSS_Process(t *testing.T) {
	backend := newOSSMockServer(t)
	defer backend.Close()

	now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return now }
	defer func() { nowFunc = time.Now }()
	ossClient := client.NewNginxClient(backend.Client(), backend.URL+"/stub_status")
	col := NewNginxCollector(ossClient, "nginx", nil, newTestLogger(),
		CTSourceProcess)

	registry := prometheus.NewRegistry()
	registry.MustRegister(col)

	cts := gatherCreatedTimestamps(t, registry)

	wantCounters := []string{
		"nginx_connections_accepted",
		"nginx_connections_handled",
		"nginx_http_requests_total",
	}
	for _, name := range wantCounters {
		got, ok := cts[name]
		if !ok {
			t.Errorf("expected created timestamp for %s but found none", name)
			continue
		}
		if !got.Equal(now) {
			t.Errorf("metric %s: got created time %v, want %v", name, got, now)
		}
	}

	// Gauge metrics should NOT have created timestamps.
	gauges := []string{
		"nginx_connections_active",
		"nginx_connections_reading",
		"nginx_connections_writing",
		"nginx_connections_waiting",
	}
	for _, name := range gauges {
		if _, ok := cts[name]; ok {
			t.Errorf("gauge metric %s should not have a created timestamp", name)
		}
	}
}

func TestCreatedTimestamp_NginxPlus_Stats(t *testing.T) {
	backend := newPlusMockServer(t)
	defer backend.Close()

	plusClient, err := plusclient.NewNginxClient(backend.URL, plusclient.WithAPIVersion(9))
	if err != nil {
		t.Fatal(err)
	}

	variableLabels := NewVariableLabelNames(nil, nil, nil, nil, nil, nil, nil)
	col := NewNginxPlusCollector(plusClient, "nginxplus", variableLabels, nil, newTestLogger(), CTSourceStats)

	registry := prometheus.NewRegistry()
	registry.MustRegister(col)

	cts := gatherCreatedTimestamps(t, registry)

	wantTime := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	wantCounters := []string{
		"nginxplus_connections_accepted",
		"nginxplus_connections_dropped",
		"nginxplus_http_requests_total",
	}
	for _, name := range wantCounters {
		got, ok := cts[name]
		if !ok {
			t.Errorf("expected created timestamp for %s but found none", name)
			continue
		}
		if !got.Equal(wantTime) {
			t.Errorf("metric %s: got created time %v, want %v", name, got, wantTime)
		}
	}

	// Gauge metrics should NOT have created timestamps.
	gauges := []string{
		"nginxplus_connections_active",
		"nginxplus_connections_idle",
	}
	for _, name := range gauges {
		if _, ok := cts[name]; ok {
			t.Errorf("gauge metric %s should not have a created timestamp", name)
		}
	}
}

func TestCreatedTimestamp_NginxPlus_Process(t *testing.T) {
	backend := newPlusMockServer(t)
	defer backend.Close()

	now := time.Date(2025, 7, 1, 10, 30, 0, 0, time.UTC)
	nowFunc = func() time.Time { return now }
	defer func() { nowFunc = time.Now }()
	plusClient, err := plusclient.NewNginxClient(backend.URL, plusclient.WithAPIVersion(9))
	if err != nil {
		t.Fatal(err)
	}

	variableLabels := NewVariableLabelNames(nil, nil, nil, nil, nil, nil, nil)
	col := NewNginxPlusCollector(plusClient, "nginxplus", variableLabels, nil, newTestLogger(), CTSourceProcess)

	registry := prometheus.NewRegistry()
	registry.MustRegister(col)

	cts := gatherCreatedTimestamps(t, registry)

	wantCounters := []string{
		"nginxplus_connections_accepted",
		"nginxplus_connections_dropped",
		"nginxplus_http_requests_total",
	}
	for _, name := range wantCounters {
		got, ok := cts[name]
		if !ok {
			t.Errorf("expected created timestamp for %s but found none", name)
			continue
		}
		if !got.Equal(now) {
			t.Errorf("metric %s: got created time %v, want %v", name, got, now)
		}
	}
}

func TestCreatedTimestamp_None(t *testing.T) {
	backend := newOSSMockServer(t)
	defer backend.Close()

	ossClient := client.NewNginxClient(backend.Client(), backend.URL+"/stub_status")
	col := NewNginxCollector(ossClient, "nginx", nil, newTestLogger(), CTSourceNone)

	registry := prometheus.NewRegistry()
	registry.MustRegister(col)

	cts := gatherCreatedTimestamps(t, registry)

	if len(cts) != 0 {
		t.Errorf("expected no created timestamps with source=none, got %v", cts)
	}
}
