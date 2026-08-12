package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewMetricLabelPrealloc(t *testing.T) {
	t.Parallel()

	type ctor func(namespace, metricName, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc

	tests := []struct {
		fn             ctor
		name           string
		wantFixedLabel []string
	}{
		{name: "newServerZoneMetric", fn: newServerZoneMetric, wantFixedLabel: []string{"server_zone"}},
		{name: "newStreamServerZoneMetric", fn: newStreamServerZoneMetric, wantFixedLabel: []string{"server_zone"}},
		{name: "newUpstreamServerMetric", fn: newUpstreamServerMetric, wantFixedLabel: []string{"upstream", "server"}},
		{name: "newStreamUpstreamServerMetric", fn: newStreamUpstreamServerMetric, wantFixedLabel: []string{"upstream", "server"}},
		{name: "newCacheZoneMetric", fn: newCacheZoneMetric, wantFixedLabel: []string{"zone"}},
	}

	variable := []string{"host", "pod"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			desc := tt.fn("nginxplus", "metric", "help", variable, prometheus.Labels{"const": "v"})
			if desc == nil {
				t.Fatal("expected non-nil desc")
			}
			s := desc.String()
			for _, l := range tt.wantFixedLabel {
				if !strings.Contains(s, l) {
					t.Errorf("desc %q missing fixed label %q", s, l)
				}
			}
			for _, l := range variable {
				if !strings.Contains(s, l) {
					t.Errorf("desc %q missing variable label %q", s, l)
				}
			}
		})
	}
}
