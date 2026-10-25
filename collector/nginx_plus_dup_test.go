package collector

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	plusclient "github.com/nginx/nginx-plus-go-client/v3/client"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNginxPlusCollectorDuplicateUpstreamServers(t *testing.T) {
	t.Parallel()

	peers := `[
		{"id":0,"server":"10.0.0.1:80","name":"10.0.0.1:80","state":"up"},
		{"id":1,"server":"10.0.0.1:80","name":"10.0.0.1:80","state":"up"}
	]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		p := strings.TrimPrefix(r.URL.Path, "/api/9")
		switch {
		case r.URL.Path == "/api/":
			_, _ = io.WriteString(w, `[1,2,3,4,5,6,7,8,9]`)
		case p == "/" || p == "":
			_, _ = io.WriteString(w, `["nginx","processes","connections","slabs","http","stream","resolvers","ssl","workers"]`)
		case p == "/http":
			_, _ = io.WriteString(w, `["requests","server_zones","location_zones","caches","limit_reqs","limit_conns","upstreams"]`)
		case p == "/stream":
			_, _ = io.WriteString(w, `["server_zones","limit_conns","upstreams","zone_sync"]`)
		case p == "/http/upstreams":
			_, _ = io.WriteString(w, `{"backend":{"peers":`+peers+`,"zone":"backend"}}`)
		case p == "/stream/upstreams":
			_, _ = io.WriteString(w, `{"backend":{"peers":`+peers+`,"zone":"backend"}}`)
		case p == "/workers":
			_, _ = io.WriteString(w, `[]`)
		default:
			_, _ = io.WriteString(w, `{}`)
		}
	}))
	defer srv.Close()

	client, err := plusclient.NewNginxClient(srv.URL+"/api", plusclient.WithAPIVersion(9))
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(NewNginxPlusCollector(client, "nginxplus", VariableLabelNames{}, nil, slog.New(slog.DiscardHandler)))
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}
	found := 0
	for _, mf := range mfs {
		if mf.GetName() == "nginxplus_upstream_server_state" || mf.GetName() == "nginxplus_stream_upstream_server_state" {
			if n := len(mf.GetMetric()); n != 2 {
				t.Errorf("%s: got %d series, want 2", mf.GetName(), n)
			}
			found++
		}
	}
	if found != 2 {
		t.Fatalf("upstream server metrics missing, got %d of 2 families", found)
	}
}
