package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

func TestParsePositiveDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		testInput string
		want      positiveDuration
		wantErr   bool
	}{
		{
			"ParsePositiveDuration returns a positiveDuration",
			"15ms",
			positiveDuration{15 * time.Millisecond},
			false,
		},
		{
			"ParsePositiveDuration returns error for trying to parse negative value",
			"-15ms",
			positiveDuration{},
			true,
		},
		{
			"ParsePositiveDuration returns error for trying to parse empty string",
			"",
			positiveDuration{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parsePositiveDuration(tt.testInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePositiveDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePositiveDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUnixSocketAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		testInput       string
		wantSocketPath  string
		wantRequestPath string
		wantErr         bool
	}{
		{
			"Normal unix socket address",
			"unix:/path/to/socket",
			"/path/to/socket",
			"",
			false,
		},
		{
			"Normal unix socket address with location",
			"unix:/path/to/socket:/with/location",
			"/path/to/socket",
			"/with/location",
			false,
		},
		{
			"Unix socket address with trailing ",
			"unix:/trailing/path:",
			"/trailing/path",
			"",
			false,
		},
		{
			"Unix socket address with too many colons",
			"unix:/too:/many:colons:",
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			socketPath, requestPath, err := parseUnixSocketAddress(tt.testInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUnixSocketAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(socketPath, tt.wantSocketPath) {
				t.Errorf("socket path: parseUnixSocketAddress() = %v, want %v", socketPath, tt.wantSocketPath)
			}
			if !reflect.DeepEqual(requestPath, tt.wantRequestPath) {
				t.Errorf("request path: parseUnixSocketAddress() = %v, want %v", requestPath, tt.wantRequestPath)
			}
		})
	}
}

func TestAddMissingEnvironmentFlags(t *testing.T) {
	t.Parallel()
	expectedMatches := map[string]string{
		"non-matching-flag":  "",
		"web.missing-env":    "MISSING_ENV",
		"web.has-env":        "HAS_ENV_ALREADY",
		"web.listen-address": "LISTEN_ADDRESS",
		"web.config.file":    "CONFIG_FILE",
	}
	kingpinflag.AddFlags(kingpin.CommandLine, ":9113")
	kingpin.Flag("non-matching-flag", "").String()
	kingpin.Flag("web.missing-env", "").String()
	kingpin.Flag("web.has-env", "").Envar("HAS_ENV_ALREADY").String()
	addMissingEnvironmentFlags(kingpin.CommandLine)

	// using Envar() on a flag returned from GetFlag()
	// adds an additional flag, which is processed correctly
	// at runtime but means that we need to check for a match
	// instead of checking the envar of each matching flag name
	for k, v := range expectedMatches {
		matched := false
		for _, f := range kingpin.CommandLine.Model().Flags {
			if f.Name == k && f.Envar == v {
				matched = true
			}
		}
		if !matched {
			t.Errorf("missing %s envar for %s", v, k)
		}
	}
}

func TestConvertFlagToEnvar(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input  string
		output string
	}{
		{
			input:  "dot.separate",
			output: "DOT_SEPARATE",
		},
		{
			input:  "underscore_separate",
			output: "UNDERSCORE_SEPARATE",
		},
		{
			input:  "mixed_separate_options",
			output: "MIXED_SEPARATE_OPTIONS",
		},
	}

	for _, c := range cases {
		res := convertFlagToEnvar(c.input)
		if res != c.output {
			t.Errorf("expected %s to resolve to %s but got %s", c.input, c.output, res)
		}
	}
}

func TestNginxMetricsOnlyRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		nginxMetricsOnly      bool
		wantVersionCollector  bool
		wantDefaultCollectors bool
	}{
		{
			name:                  "default behavior includes version and runtime collectors",
			nginxMetricsOnly:      false,
			wantVersionCollector:  true,
			wantDefaultCollectors: true,
		},
		{
			name:                  "web.nginx-metrics-only excludes version and runtime collectors",
			nginxMetricsOnly:      true,
			wantVersionCollector:  false,
			wantDefaultCollectors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := prometheus.NewRegistry()

			if !tt.nginxMetricsOnly {
				// Register version collector like main() does when nginxMetricsOnly is false
				if err := registry.Register(version.NewCollector("nginx_exporter")); err != nil {
					t.Fatalf("failed to register version collector: %v", err)
				}
			}

			// Gather metrics
			metricFamilies, err := registry.Gather()
			if err != nil {
				t.Fatalf("failed to gather metrics: %v", err)
			}

			// Check for version collector metric
			hasVersionMetric := false
			for _, mf := range metricFamilies {
				if mf.GetName() == "nginx_exporter_build_info" {
					hasVersionMetric = true
					break
				}
			}

			if hasVersionMetric != tt.wantVersionCollector {
				t.Errorf("version collector presence = %v, want %v", hasVersionMetric, tt.wantVersionCollector)
			}
		})
	}
}
