package collector

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	plusclient "github.com/nginx/nginx-plus-go-client/v3/client"
	"github.com/prometheus/client_golang/prometheus"
)

// LabelUpdater updates the labels of upstream server and server zone metrics.
type LabelUpdater interface {
	UpdateUpstreamServerPeerLabels(upstreamServerPeerLabels map[string][]string)
	DeleteUpstreamServerPeerLabels(peers []string)
	UpdateUpstreamServerLabels(upstreamServerLabelValues map[string][]string)
	DeleteUpstreamServerLabels(upstreamNames []string)
	UpdateStreamUpstreamServerPeerLabels(streamUpstreamServerPeerLabels map[string][]string)
	DeleteStreamUpstreamServerPeerLabels(peers []string)
	UpdateStreamUpstreamServerLabels(streamUpstreamServerPeerLabels map[string][]string)
	DeleteStreamUpstreamServerLabels(peers []string)
	UpdateServerZoneLabels(serverZoneLabelValues map[string][]string)
	DeleteServerZoneLabels(zoneNames []string)
	UpdateStreamServerZoneLabels(streamServerZoneLabelValues map[string][]string)
	DeleteStreamServerZoneLabels(zoneNames []string)
	UpdateCacheZoneLabels(cacheLabelValues map[string][]string)
	DeleteCacheZoneLabels(cacheNames []string)
}

// NginxPlusCollector collects NGINX Plus metrics. It implements prometheus.Collector interface.
type NginxPlusCollector struct {
	upMetric                       prometheus.Gauge
	logger                         *slog.Logger
	cacheZoneMetrics               map[string]*prometheus.Desc
	workerMetrics                  map[string]*prometheus.Desc
	nginxClient                    *plusclient.NginxClient
	streamServerZoneMetrics        map[string]*prometheus.Desc
	streamZoneSyncMetrics          map[string]*prometheus.Desc
	streamUpstreamMetrics          map[string]*prometheus.Desc
	streamUpstreamServerMetrics    map[string]*prometheus.Desc
	locationZoneMetrics            map[string]*prometheus.Desc
	resolverMetrics                map[string]*prometheus.Desc
	limitRequestMetrics            map[string]*prometheus.Desc
	limitConnectionMetrics         map[string]*prometheus.Desc
	streamLimitConnectionMetrics   map[string]*prometheus.Desc
	upstreamServerMetrics          map[string]*prometheus.Desc
	upstreamMetrics                map[string]*prometheus.Desc
	streamUpstreamServerPeerLabels map[string][]string
	serverZoneMetrics              map[string]*prometheus.Desc
	upstreamServerLabels           map[string][]string
	streamUpstreamServerLabels     map[string][]string
	serverZoneLabels               map[string][]string
	streamServerZoneLabels         map[string][]string
	upstreamServerPeerLabels       map[string][]string
	cacheZoneLabels                map[string][]string
	totalMetrics                   map[string]*prometheus.Desc
	variableLabelNames             VariableLabelNames
	variableLabelsMutex            sync.RWMutex
	mutex                          sync.Mutex
	createdAtFunc                  func(*plusclient.Stats) time.Time
}

// UpdateUpstreamServerPeerLabels updates the Upstream Server Peer Labels.
func (c *NginxPlusCollector) UpdateUpstreamServerPeerLabels(upstreamServerPeerLabels map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range upstreamServerPeerLabels {
		c.upstreamServerPeerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteUpstreamServerPeerLabels deletes the Upstream Server Peer Labels.
func (c *NginxPlusCollector) DeleteUpstreamServerPeerLabels(peers []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range peers {
		delete(c.upstreamServerPeerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateStreamUpstreamServerPeerLabels updates the Upstream Server Peer Labels.
func (c *NginxPlusCollector) UpdateStreamUpstreamServerPeerLabels(streamUpstreamServerPeerLabels map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range streamUpstreamServerPeerLabels {
		c.streamUpstreamServerPeerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteStreamUpstreamServerPeerLabels deletes the Upstream Server Peer Labels.
func (c *NginxPlusCollector) DeleteStreamUpstreamServerPeerLabels(peers []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range peers {
		delete(c.streamUpstreamServerPeerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateUpstreamServerLabels updates the Upstream Server Labels.
func (c *NginxPlusCollector) UpdateUpstreamServerLabels(upstreamServerLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range upstreamServerLabelValues {
		c.upstreamServerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteUpstreamServerLabels deletes the Upstream Server Labels.
func (c *NginxPlusCollector) DeleteUpstreamServerLabels(upstreamNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range upstreamNames {
		delete(c.upstreamServerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateStreamUpstreamServerLabels updates the Upstream Server Labels.
func (c *NginxPlusCollector) UpdateStreamUpstreamServerLabels(streamUpstreamServerLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range streamUpstreamServerLabelValues {
		c.streamUpstreamServerLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteStreamUpstreamServerLabels deletes the Upstream Server Labels.
func (c *NginxPlusCollector) DeleteStreamUpstreamServerLabels(streamUpstreamNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range streamUpstreamNames {
		delete(c.streamUpstreamServerLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateServerZoneLabels updates the Server Zone Labels.
func (c *NginxPlusCollector) UpdateServerZoneLabels(serverZoneLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range serverZoneLabelValues {
		c.serverZoneLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteServerZoneLabels deletes the Server Zone Labels.
func (c *NginxPlusCollector) DeleteServerZoneLabels(zoneNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range zoneNames {
		delete(c.serverZoneLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateStreamServerZoneLabels updates the Stream Server Zone Labels.
func (c *NginxPlusCollector) UpdateStreamServerZoneLabels(streamServerZoneLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range streamServerZoneLabelValues {
		c.streamServerZoneLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteStreamServerZoneLabels deletes the Stream Server Zone Labels.
func (c *NginxPlusCollector) DeleteStreamServerZoneLabels(zoneNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range zoneNames {
		delete(c.streamServerZoneLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

// UpdateCacheZoneLabels updates the Upstream Cache Zone labels.
func (c *NginxPlusCollector) UpdateCacheZoneLabels(cacheZoneLabelValues map[string][]string) {
	c.variableLabelsMutex.Lock()
	for k, v := range cacheZoneLabelValues {
		c.cacheZoneLabels[k] = v
	}
	c.variableLabelsMutex.Unlock()
}

// DeleteCacheZoneLabels deletes the Cache Zone Labels.
func (c *NginxPlusCollector) DeleteCacheZoneLabels(cacheZoneNames []string) {
	c.variableLabelsMutex.Lock()
	for _, k := range cacheZoneNames {
		delete(c.cacheZoneLabels, k)
	}
	c.variableLabelsMutex.Unlock()
}

func (c *NginxPlusCollector) getUpstreamServerLabelValues(upstreamName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.upstreamServerLabels[upstreamName]
}

func (c *NginxPlusCollector) getStreamUpstreamServerLabelValues(upstreamName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.streamUpstreamServerLabels[upstreamName]
}

func (c *NginxPlusCollector) getServerZoneLabelValues(zoneName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.serverZoneLabels[zoneName]
}

func (c *NginxPlusCollector) getStreamServerZoneLabelValues(zoneName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.streamServerZoneLabels[zoneName]
}

func (c *NginxPlusCollector) getUpstreamServerPeerLabelValues(peer string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.upstreamServerPeerLabels[peer]
}

func (c *NginxPlusCollector) getStreamUpstreamServerPeerLabelValues(peer string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.streamUpstreamServerPeerLabels[peer]
}

func (c *NginxPlusCollector) getCacheZoneLabelValues(cacheName string) []string {
	c.variableLabelsMutex.RLock()
	defer c.variableLabelsMutex.RUnlock()
	return c.cacheZoneLabels[cacheName]
}

// VariableLabelNames holds all the variable label names for the different metrics.
type VariableLabelNames struct {
	UpstreamServerVariableLabelNames           []string
	ServerZoneVariableLabelNames               []string
	UpstreamServerPeerVariableLabelNames       []string
	StreamUpstreamServerPeerVariableLabelNames []string
	StreamServerZoneVariableLabelNames         []string
	StreamUpstreamServerVariableLabelNames     []string
	CacheZoneVariableLabelNames                []string
}

// NewVariableLabelNames NewVariableLabels creates a new struct for VariableNames for the collector.
func NewVariableLabelNames(upstreamServerVariableLabelNames []string, serverZoneVariableLabelNames []string, upstreamServerPeerVariableLabelNames []string,
	streamUpstreamServerVariableLabelNames []string, streamServerZoneLabels []string, streamUpstreamServerPeerVariableLabelNames []string, cacheZoneVariableLabelNames []string,
) VariableLabelNames {
	return VariableLabelNames{
		UpstreamServerVariableLabelNames:           upstreamServerVariableLabelNames,
		ServerZoneVariableLabelNames:               serverZoneVariableLabelNames,
		UpstreamServerPeerVariableLabelNames:       upstreamServerPeerVariableLabelNames,
		StreamUpstreamServerVariableLabelNames:     streamUpstreamServerVariableLabelNames,
		StreamServerZoneVariableLabelNames:         streamServerZoneLabels,
		StreamUpstreamServerPeerVariableLabelNames: streamUpstreamServerPeerVariableLabelNames,
		CacheZoneVariableLabelNames:                cacheZoneVariableLabelNames,
	}
}

// NewNginxPlusCollector creates an NginxPlusCollector.
func NewNginxPlusCollector(nginxClient *plusclient.NginxClient, namespace string, variableLabelNames VariableLabelNames, constLabels map[string]string, logger *slog.Logger, ctSource CTSource) *NginxPlusCollector {
	upstreamServerVariableLabelNames := variableLabelNames.UpstreamServerVariableLabelNames
	streamUpstreamServerVariableLabelNames := variableLabelNames.StreamUpstreamServerVariableLabelNames

	upstreamServerVariableLabelNames = append(upstreamServerVariableLabelNames, variableLabelNames.UpstreamServerPeerVariableLabelNames...)
	streamUpstreamServerVariableLabelNames = append(streamUpstreamServerVariableLabelNames, variableLabelNames.StreamUpstreamServerPeerVariableLabelNames...)
	c := &NginxPlusCollector{
		variableLabelNames:             variableLabelNames,
		upstreamServerLabels:           make(map[string][]string),
		serverZoneLabels:               make(map[string][]string),
		streamServerZoneLabels:         make(map[string][]string),
		upstreamServerPeerLabels:       make(map[string][]string),
		streamUpstreamServerPeerLabels: make(map[string][]string),
		streamUpstreamServerLabels:     make(map[string][]string),
		cacheZoneLabels:                make(map[string][]string),
		nginxClient:                    nginxClient,
		logger:                         logger,
		totalMetrics: map[string]*prometheus.Desc{
			"connections_accepted":           newGlobalMetric(namespace, "connections_accepted", "Accepted client connections", constLabels),
			"connections_dropped":            newGlobalMetric(namespace, "connections_dropped", "Dropped client connections", constLabels),
			"connections_active":             newGlobalMetric(namespace, "connections_active", "Active client connections", constLabels),
			"connections_idle":               newGlobalMetric(namespace, "connections_idle", "Idle client connections", constLabels),
			"http_requests_total":            newGlobalMetric(namespace, "http_requests_total", "Total http requests", constLabels),
			"http_requests_current":          newGlobalMetric(namespace, "http_requests_current", "Current http requests", constLabels),
			"ssl_handshakes":                 newGlobalMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", constLabels),
			"ssl_handshakes_failed":          newGlobalMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", constLabels),
			"ssl_session_reuses":             newGlobalMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", constLabels),
			"license_active_till":            newGlobalMetric(namespace, "license_expiration_timestamp_seconds", "License expiration date (expressed as Unix Epoch Time)", constLabels),
			"license_reporting_healthy":      newGlobalMetric(namespace, "license_reporting_healthy", "Indicates whether the reporting state is still considered healthy despite recent failed attempts", constLabels),
			"license_reporting_fails":        newGlobalMetric(namespace, "license_reporting_fails_count", "Number of failed reporting attempts, reset each time the usage report is successfully sent", constLabels),
			"license_reporting_grace_period": newGlobalMetric(namespace, "license_reporting_grace_period_seconds", "Number of seconds before traffic processing is stopped after unsuccessful report attempt", constLabels),
		},
		serverZoneMetrics: map[string]*prometheus.Desc{
			"processing":            newServerZoneMetric(namespace, "processing", "Client requests that are currently being processed", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"requests":              newServerZoneMetric(namespace, "requests", "Total client requests", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"responses_1xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "1xx"})),
			"responses_2xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"responses_3xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "3xx"})),
			"responses_4xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"responses_5xx":         newServerZoneMetric(namespace, "responses", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"discarded":             newServerZoneMetric(namespace, "discarded", "Requests completed without sending a response", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"received":              newServerZoneMetric(namespace, "received", "Bytes received from clients", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"sent":                  newServerZoneMetric(namespace, "sent", "Bytes sent to clients", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"codes_100":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "100"})),
			"codes_101":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "101"})),
			"codes_102":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "102"})),
			"codes_200":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "200"})),
			"codes_201":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "201"})),
			"codes_202":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "202"})),
			"codes_204":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "204"})),
			"codes_206":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "206"})),
			"codes_300":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "300"})),
			"codes_301":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "301"})),
			"codes_302":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "302"})),
			"codes_303":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "303"})),
			"codes_304":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "304"})),
			"codes_307":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "307"})),
			"codes_400":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "400"})),
			"codes_401":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "401"})),
			"codes_403":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "403"})),
			"codes_404":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "404"})),
			"codes_405":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "405"})),
			"codes_408":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "408"})),
			"codes_409":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "409"})),
			"codes_411":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "411"})),
			"codes_412":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "412"})),
			"codes_413":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "413"})),
			"codes_414":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "414"})),
			"codes_415":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "415"})),
			"codes_416":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "416"})),
			"codes_429":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "429"})),
			"codes_444":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "444"})),
			"codes_494":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "494"})),
			"codes_495":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "495"})),
			"codes_496":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "496"})),
			"codes_497":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "497"})),
			"codes_499":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "499"})),
			"codes_500":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "500"})),
			"codes_501":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "501"})),
			"codes_502":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "502"})),
			"codes_503":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "503"})),
			"codes_504":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "504"})),
			"codes_507":             newServerZoneMetric(namespace, "responses_codes", "Total responses sent to clients", variableLabelNames.ServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "507"})),
			"ssl_handshakes":        newServerZoneMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"ssl_handshakes_failed": newServerZoneMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
			"ssl_session_reuses":    newServerZoneMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", variableLabelNames.ServerZoneVariableLabelNames, constLabels),
		},
		streamServerZoneMetrics: map[string]*prometheus.Desc{
			"processing":            newStreamServerZoneMetric(namespace, "processing", "Client connections that are currently being processed", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"connections":           newStreamServerZoneMetric(namespace, "connections", "Total connections", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"sessions_2xx":          newStreamServerZoneMetric(namespace, "sessions", "Total sessions completed", variableLabelNames.StreamServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"sessions_4xx":          newStreamServerZoneMetric(namespace, "sessions", "Total sessions completed", variableLabelNames.StreamServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"sessions_5xx":          newStreamServerZoneMetric(namespace, "sessions", "Total sessions completed", variableLabelNames.StreamServerZoneVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"discarded":             newStreamServerZoneMetric(namespace, "discarded", "Connections completed without creating a session", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"received":              newStreamServerZoneMetric(namespace, "received", "Bytes received from clients", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"sent":                  newStreamServerZoneMetric(namespace, "sent", "Bytes sent to clients", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"ssl_handshakes":        newStreamServerZoneMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"ssl_handshakes_failed": newStreamServerZoneMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
			"ssl_session_reuses":    newStreamServerZoneMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", variableLabelNames.StreamServerZoneVariableLabelNames, constLabels),
		},
		upstreamMetrics: map[string]*prometheus.Desc{
			"keepalive": newUpstreamMetric(namespace, "keepalive", "Idle keepalive connections", constLabels),
			"zombies":   newUpstreamMetric(namespace, "zombies", "Servers removed from the group but still processing active client requests", constLabels),
		},
		streamUpstreamMetrics: map[string]*prometheus.Desc{
			"zombies": newStreamUpstreamMetric(namespace, "zombies", "Servers removed from the group but still processing active client connections", constLabels),
		},
		upstreamServerMetrics: map[string]*prometheus.Desc{
			"state":                   newUpstreamServerMetric(namespace, "state", "Current state", upstreamServerVariableLabelNames, constLabels),
			"active":                  newUpstreamServerMetric(namespace, "active", "Active connections", upstreamServerVariableLabelNames, constLabels),
			"limit":                   newUpstreamServerMetric(namespace, "limit", "Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit", upstreamServerVariableLabelNames, constLabels),
			"requests":                newUpstreamServerMetric(namespace, "requests", "Total client requests", upstreamServerVariableLabelNames, constLabels),
			"responses_1xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "1xx"})),
			"responses_2xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"responses_3xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "3xx"})),
			"responses_4xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"responses_5xx":           newUpstreamServerMetric(namespace, "responses", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"sent":                    newUpstreamServerMetric(namespace, "sent", "Bytes sent to this server", upstreamServerVariableLabelNames, constLabels),
			"received":                newUpstreamServerMetric(namespace, "received", "Bytes received to this server", upstreamServerVariableLabelNames, constLabels),
			"fails":                   newUpstreamServerMetric(namespace, "fails", "Number of unsuccessful attempts to communicate with the server", upstreamServerVariableLabelNames, constLabels),
			"unavail":                 newUpstreamServerMetric(namespace, "unavail", "How many times the server became unavailable for client requests (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold", upstreamServerVariableLabelNames, constLabels),
			"header_time":             newUpstreamServerMetric(namespace, "header_time", "Average time to get the response header from the server", upstreamServerVariableLabelNames, constLabels),
			"response_time":           newUpstreamServerMetric(namespace, "response_time", "Average time to get the full response from the server", upstreamServerVariableLabelNames, constLabels),
			"health_checks_checks":    newUpstreamServerMetric(namespace, "health_checks_checks", "Total health check requests", upstreamServerVariableLabelNames, constLabels),
			"health_checks_fails":     newUpstreamServerMetric(namespace, "health_checks_fails", "Failed health checks", upstreamServerVariableLabelNames, constLabels),
			"health_checks_unhealthy": newUpstreamServerMetric(namespace, "health_checks_unhealthy", "How many times the server became unhealthy (state 'unhealthy')", upstreamServerVariableLabelNames, constLabels),
			"codes_100":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "100"})),
			"codes_101":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "101"})),
			"codes_102":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "102"})),
			"codes_200":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "200"})),
			"codes_201":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "201"})),
			"codes_202":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "202"})),
			"codes_204":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "204"})),
			"codes_206":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "206"})),
			"codes_300":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "300"})),
			"codes_301":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "301"})),
			"codes_302":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "302"})),
			"codes_303":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "303"})),
			"codes_304":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "304"})),
			"codes_307":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "307"})),
			"codes_400":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "400"})),
			"codes_401":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "401"})),
			"codes_403":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "403"})),
			"codes_404":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "404"})),
			"codes_405":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "405"})),
			"codes_408":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "408"})),
			"codes_409":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "409"})),
			"codes_411":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "411"})),
			"codes_412":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "412"})),
			"codes_413":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "413"})),
			"codes_414":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "414"})),
			"codes_415":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "415"})),
			"codes_416":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "416"})),
			"codes_429":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "429"})),
			"codes_444":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "444"})),
			"codes_494":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "494"})),
			"codes_495":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "495"})),
			"codes_496":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "496"})),
			"codes_497":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "497"})),
			"codes_499":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "499"})),
			"codes_500":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "500"})),
			"codes_501":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "501"})),
			"codes_502":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "502"})),
			"codes_503":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "503"})),
			"codes_504":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "504"})),
			"codes_507":               newUpstreamServerMetric(namespace, "responses_codes", "Total responses sent to clients", upstreamServerVariableLabelNames, MergeLabels(constLabels, prometheus.Labels{"code": "507"})),
			"ssl_handshakes":          newUpstreamServerMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", upstreamServerVariableLabelNames, constLabels),
			"ssl_handshakes_failed":   newUpstreamServerMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", upstreamServerVariableLabelNames, constLabels),
			"ssl_session_reuses":      newUpstreamServerMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", upstreamServerVariableLabelNames, constLabels),
		},
		streamUpstreamServerMetrics: map[string]*prometheus.Desc{
			"state":                   newStreamUpstreamServerMetric(namespace, "state", "Current state", streamUpstreamServerVariableLabelNames, constLabels),
			"active":                  newStreamUpstreamServerMetric(namespace, "active", "Active connections", streamUpstreamServerVariableLabelNames, constLabels),
			"limit":                   newStreamUpstreamServerMetric(namespace, "limit", "Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit", streamUpstreamServerVariableLabelNames, constLabels),
			"sent":                    newStreamUpstreamServerMetric(namespace, "sent", "Bytes sent to this server", streamUpstreamServerVariableLabelNames, constLabels),
			"received":                newStreamUpstreamServerMetric(namespace, "received", "Bytes received from this server", streamUpstreamServerVariableLabelNames, constLabels),
			"fails":                   newStreamUpstreamServerMetric(namespace, "fails", "Number of unsuccessful attempts to communicate with the server", streamUpstreamServerVariableLabelNames, constLabels),
			"unavail":                 newStreamUpstreamServerMetric(namespace, "unavail", "How many times the server became unavailable for client connections (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold", streamUpstreamServerVariableLabelNames, constLabels),
			"connections":             newStreamUpstreamServerMetric(namespace, "connections", "Total number of client connections forwarded to this server", streamUpstreamServerVariableLabelNames, constLabels),
			"connect_time":            newStreamUpstreamServerMetric(namespace, "connect_time", "Average time to connect to the upstream server", streamUpstreamServerVariableLabelNames, constLabels),
			"first_byte_time":         newStreamUpstreamServerMetric(namespace, "first_byte_time", "Average time to receive the first byte of data", streamUpstreamServerVariableLabelNames, constLabels),
			"response_time":           newStreamUpstreamServerMetric(namespace, "response_time", "Average time to receive the last byte of data", streamUpstreamServerVariableLabelNames, constLabels),
			"health_checks_checks":    newStreamUpstreamServerMetric(namespace, "health_checks_checks", "Total health check requests", streamUpstreamServerVariableLabelNames, constLabels),
			"health_checks_fails":     newStreamUpstreamServerMetric(namespace, "health_checks_fails", "Failed health checks", streamUpstreamServerVariableLabelNames, constLabels),
			"health_checks_unhealthy": newStreamUpstreamServerMetric(namespace, "health_checks_unhealthy", "How many times the server became unhealthy (state 'unhealthy')", streamUpstreamServerVariableLabelNames, constLabels),
			"ssl_handshakes":          newStreamUpstreamServerMetric(namespace, "ssl_handshakes", "Successful SSL handshakes", streamUpstreamServerVariableLabelNames, constLabels),
			"ssl_handshakes_failed":   newStreamUpstreamServerMetric(namespace, "ssl_handshakes_failed", "Failed SSL handshakes", streamUpstreamServerVariableLabelNames, constLabels),
			"ssl_session_reuses":      newStreamUpstreamServerMetric(namespace, "ssl_session_reuses", "Session reuses during SSL handshake", streamUpstreamServerVariableLabelNames, constLabels),
		},
		streamZoneSyncMetrics: map[string]*prometheus.Desc{
			"bytes_in":        newStreamZoneSyncMetric(namespace, "bytes_in", "Bytes received by this node", constLabels),
			"bytes_out":       newStreamZoneSyncMetric(namespace, "bytes_out", "Bytes sent by this node", constLabels),
			"msgs_in":         newStreamZoneSyncMetric(namespace, "msgs_in", "Total messages received by this node", constLabels),
			"msgs_out":        newStreamZoneSyncMetric(namespace, "msgs_out", "Total messages sent by this node", constLabels),
			"nodes_online":    newStreamZoneSyncMetric(namespace, "nodes_online", "Number of peers this node is connected to", constLabels),
			"records_pending": newStreamZoneSyncZoneMetric(namespace, "records_pending", "The number of records that need to be sent to the cluster", constLabels),
			"records_total":   newStreamZoneSyncZoneMetric(namespace, "records_total", "The total number of records stored in the shared memory zone", constLabels),
		},
		locationZoneMetrics: map[string]*prometheus.Desc{
			"requests":      newLocationZoneMetric(namespace, "requests", "Total client requests", constLabels),
			"responses_1xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "1xx"})),
			"responses_2xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "2xx"})),
			"responses_3xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "3xx"})),
			"responses_4xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "4xx"})),
			"responses_5xx": newLocationZoneMetric(namespace, "responses", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "5xx"})),
			"discarded":     newLocationZoneMetric(namespace, "discarded", "Requests completed without sending a response", constLabels),
			"received":      newLocationZoneMetric(namespace, "received", "Bytes received from clients", constLabels),
			"sent":          newLocationZoneMetric(namespace, "sent", "Bytes sent to clients", constLabels),
			"codes_100":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "100"})),
			"codes_101":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "101"})),
			"codes_102":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "102"})),
			"codes_200":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "200"})),
			"codes_201":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "201"})),
			"codes_202":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "202"})),
			"codes_204":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "204"})),
			"codes_206":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "206"})),
			"codes_300":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "300"})),
			"codes_301":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "301"})),
			"codes_302":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "302"})),
			"codes_303":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "303"})),
			"codes_304":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "304"})),
			"codes_307":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "307"})),
			"codes_400":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "400"})),
			"codes_401":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "401"})),
			"codes_403":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "403"})),
			"codes_404":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "404"})),
			"codes_405":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "405"})),
			"codes_408":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "408"})),
			"codes_409":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "409"})),
			"codes_411":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "411"})),
			"codes_412":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "412"})),
			"codes_413":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "413"})),
			"codes_414":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "414"})),
			"codes_415":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "415"})),
			"codes_416":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "416"})),
			"codes_429":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "429"})),
			"codes_444":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "444"})),
			"codes_494":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "494"})),
			"codes_495":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "495"})),
			"codes_496":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "496"})),
			"codes_497":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "497"})),
			"codes_499":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "499"})),
			"codes_500":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "500"})),
			"codes_501":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "501"})),
			"codes_502":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "502"})),
			"codes_503":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "503"})),
			"codes_504":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "504"})),
			"codes_507":     newLocationZoneMetric(namespace, "responses_codes", "Total responses sent to clients", MergeLabels(constLabels, prometheus.Labels{"code": "507"})),
		},
		resolverMetrics: map[string]*prometheus.Desc{
			"name":     newResolverMetric(namespace, "name", "Total requests to resolve names to addresses", constLabels),
			"srv":      newResolverMetric(namespace, "srv", "Total requests to resolve SRV records", constLabels),
			"addr":     newResolverMetric(namespace, "addr", "Total requests to resolve addresses to names", constLabels),
			"noerror":  newResolverMetric(namespace, "noerror", "Total number of successful responses", constLabels),
			"formerr":  newResolverMetric(namespace, "formerr", "Total number of FORMERR responses", constLabels),
			"servfail": newResolverMetric(namespace, "servfail", "Total number of SERVFAIL responses", constLabels),
			"nxdomain": newResolverMetric(namespace, "nxdomain", "Total number of NXDOMAIN responses", constLabels),
			"notimp":   newResolverMetric(namespace, "notimp", "Total number of NOTIMP responses", constLabels),
			"refused":  newResolverMetric(namespace, "refused", "Total number of REFUSED responses", constLabels),
			"timedout": newResolverMetric(namespace, "timedout", "Total number of timed out requests", constLabels),
			"unknown":  newResolverMetric(namespace, "unknown", "Total requests completed with an unknown error", constLabels),
		},
		limitRequestMetrics: map[string]*prometheus.Desc{
			"passed":           newLimitRequestMetric(namespace, "passed", "Total number of requests that were neither limited nor accounted as limited", constLabels),
			"delayed":          newLimitRequestMetric(namespace, "delayed", "Total number of requests that were delayed", constLabels),
			"rejected":         newLimitRequestMetric(namespace, "rejected", "Total number of requests that were rejected", constLabels),
			"delayed_dry_run":  newLimitRequestMetric(namespace, "delayed_dry_run", "Total number of requests accounted as delayed in the dry run mode", constLabels),
			"rejected_dry_run": newLimitRequestMetric(namespace, "rejected_dry_run", "Total number of requests accounted as rejected in the dry run mode", constLabels),
		},
		limitConnectionMetrics: map[string]*prometheus.Desc{
			"passed":           newLimitConnectionMetric(namespace, "passed", "Total number of connections that were neither limited nor accounted as limited", constLabels),
			"rejected":         newLimitConnectionMetric(namespace, "rejected", "Total number of connections that were rejected", constLabels),
			"rejected_dry_run": newLimitConnectionMetric(namespace, "rejected_dry_run", "Total number of connections accounted as rejected in the dry run mode", constLabels),
		},
		streamLimitConnectionMetrics: map[string]*prometheus.Desc{
			"passed":           newStreamLimitConnectionMetric(namespace, "passed", "Total number of connections that were neither limited nor accounted as limited", constLabels),
			"rejected":         newStreamLimitConnectionMetric(namespace, "rejected", "Total number of connections that were rejected", constLabels),
			"rejected_dry_run": newStreamLimitConnectionMetric(namespace, "rejected_dry_run", "Total number of connections accounted as rejected in the dry run mode", constLabels),
		},
		upMetric: newUpMetric(namespace, constLabels),
		cacheZoneMetrics: map[string]*prometheus.Desc{
			"size":                      newCacheZoneMetric(namespace, "size", "Total size of the cache", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"max_size":                  newCacheZoneMetric(namespace, "max_size", "Maximum size of the cache", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"cold":                      newCacheZoneMetric(namespace, "cold", "Is the cache considered cold", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"hit_responses":             newCacheZoneMetric(namespace, "hit_responses", "Total number of cache hits", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"hit_bytes":                 newCacheZoneMetric(namespace, "hit_bytes", "Total number of bytes returned from cache", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"stale_responses":           newCacheZoneMetric(namespace, "stale_responses", "Total number of stale cache hits", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"stale_bytes":               newCacheZoneMetric(namespace, "stale_bytes", "Total number of bytes returned from stale cache", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"updating_responses":        newCacheZoneMetric(namespace, "updating_responses", "Total number of cache hits while cache is updating", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"updating_bytes":            newCacheZoneMetric(namespace, "updating_bytes", "Total number of bytes returned from cache while cache is updating", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"revalidated_responses":     newCacheZoneMetric(namespace, "revalidated_responses", "Total number of cache revalidations", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"revalidated_bytes":         newCacheZoneMetric(namespace, "revalidated_bytes", "Total number of bytes returned from cache revalidations", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"miss_responses":            newCacheZoneMetric(namespace, "miss_responses", "Total number of cache misses", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"miss_bytes":                newCacheZoneMetric(namespace, "miss_bytes", "Total number of bytes returned from cache misses", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"expired_responses":         newCacheZoneMetric(namespace, "expired_responses", "Total number of cache hits with expired TTL", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"expired_bytes":             newCacheZoneMetric(namespace, "expired_bytes", "Total number of bytes returned from cache hits with expired TTL", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"expired_responses_written": newCacheZoneMetric(namespace, "expired_responses_written", "Total number of cache hits with expired TTL written to cache", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"expired_bytes_written":     newCacheZoneMetric(namespace, "expired_bytes_written", "Total number of bytes written to cache from cache hits with expired TTL", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"bypass_responses":          newCacheZoneMetric(namespace, "bypass_responses", "Total number of cache bypasses", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"bypass_bytes":              newCacheZoneMetric(namespace, "bypass_bytes", "Total number of bytes returned from cache bypasses", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"bypass_responses_written":  newCacheZoneMetric(namespace, "bypass_responses_written", "Total number of cache bypasses written to cache", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
			"bypass_bytes_written":      newCacheZoneMetric(namespace, "bypass_bytes_written", "Total number of bytes written to cache from cache bypasses", variableLabelNames.CacheZoneVariableLabelNames, constLabels),
		},
		workerMetrics: map[string]*prometheus.Desc{
			"connection_accepted":   newWorkerMetric(namespace, "connection_accepted", "The total number of accepted client connections", constLabels),
			"connection_dropped":    newWorkerMetric(namespace, "connection_dropped", "The total number of dropped client connections", constLabels),
			"connection_active":     newWorkerMetric(namespace, "connection_active", "The current number of active client connections", constLabels),
			"connection_idle":       newWorkerMetric(namespace, "connection_idle", "The current number of idle client connections", constLabels),
			"http_requests_total":   newWorkerMetric(namespace, "http_requests_total", "The total number of client requests received by the worker process", constLabels),
			"http_requests_current": newWorkerMetric(namespace, "http_requests_current", "The current number of client requests that are currently being processed by the worker process", constLabels),
		},
	}
	if ctSource == CTSourceProcess {
		now := nowFunc()
		c.createdAtFunc = func(_ *plusclient.Stats) time.Time { return now }
	} else if ctSource == CTSourceStats {
		c.createdAtFunc = func(stats *plusclient.Stats) time.Time {
			parsed, err := time.Parse(time.RFC3339Nano, stats.NginxInfo.LoadTimestamp)
			if err != nil {
				c.logger.Warn("error parsing load_timestamp for created timestamp", "load_timestamp", stats.NginxInfo.LoadTimestamp, "error", err.Error())
				return time.Time{}
			}
			return parsed
		}
	}
	return c
}

// Describe sends the super-set of all possible descriptors of NGINX Plus metrics
// to the provided channel.
func (c *NginxPlusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upMetric.Desc()

	for _, m := range c.totalMetrics {
		ch <- m
	}
	for _, m := range c.serverZoneMetrics {
		ch <- m
	}
	for _, m := range c.upstreamMetrics {
		ch <- m
	}
	for _, m := range c.upstreamServerMetrics {
		ch <- m
	}
	for _, m := range c.streamServerZoneMetrics {
		ch <- m
	}
	for _, m := range c.streamUpstreamMetrics {
		ch <- m
	}
	for _, m := range c.streamUpstreamServerMetrics {
		ch <- m
	}
	for _, m := range c.streamZoneSyncMetrics {
		ch <- m
	}
	for _, m := range c.locationZoneMetrics {
		ch <- m
	}
	for _, m := range c.resolverMetrics {
		ch <- m
	}
	for _, m := range c.limitRequestMetrics {
		ch <- m
	}
	for _, m := range c.limitConnectionMetrics {
		ch <- m
	}
	for _, m := range c.streamLimitConnectionMetrics {
		ch <- m
	}
	for _, m := range c.cacheZoneMetrics {
		ch <- m
	}
	for _, m := range c.workerMetrics {
		ch <- m
	}
}

func (c *NginxPlusCollector) newCounterMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, createdAt time.Time, value float64, labelValues ...string) {
	if createdAt.IsZero() {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, value, labelValues...)
	} else {
		ch <- prometheus.MustNewConstMetricWithCreatedTimestamp(desc, prometheus.CounterValue, value, createdAt, labelValues...)
	}
}

func (c *NginxPlusCollector) newGaugeMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, labelValues...)
}

// Collect fetches metrics from NGINX Plus and sends them to the provided channel.
func (c *NginxPlusCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	// FIXME: https://github.com/nginx/nginx-prometheus-exporter/issues/858
	stats, err := c.nginxClient.GetStats(context.TODO())
	if err != nil {
		c.upMetric.Set(nginxDown)
		ch <- c.upMetric
		c.logger.Warn("error getting stats", "error", err.Error())
		return
	}

	c.upMetric.Set(nginxUp)
	ch <- c.upMetric

	createdAt := c.createdAtFunc(stats)

	c.newCounterMetric(ch, c.totalMetrics["connections_accepted"], createdAt,
		float64(stats.Connections.Accepted))
	c.newCounterMetric(ch, c.totalMetrics["connections_dropped"], createdAt,
		float64(stats.Connections.Dropped))
	c.newGaugeMetric(ch, c.totalMetrics["connections_active"],
		float64(stats.Connections.Active))
	c.newGaugeMetric(ch, c.totalMetrics["connections_idle"],
		float64(stats.Connections.Idle))
	c.newCounterMetric(ch, c.totalMetrics["http_requests_total"], createdAt,
		float64(stats.HTTPRequests.Total))
	c.newGaugeMetric(ch, c.totalMetrics["http_requests_current"],
		float64(stats.HTTPRequests.Current))
	c.newCounterMetric(ch, c.totalMetrics["ssl_handshakes"], createdAt,
		float64(stats.SSL.Handshakes))
	c.newCounterMetric(ch, c.totalMetrics["ssl_handshakes_failed"], createdAt,
		float64(stats.SSL.HandshakesFailed))
	c.newCounterMetric(ch, c.totalMetrics["ssl_session_reuses"], createdAt,
		float64(stats.SSL.SessionReuses))

	license, err := c.nginxClient.GetNginxLicense(context.TODO())
	if err != nil {
		c.logger.Warn("error getting license information", "error", err.Error())
	} else {
		c.newGaugeMetric(ch, c.totalMetrics["license_active_till"],
			float64(license.ActiveTill))

		if license.Reporting != nil {
			if license.Reporting.Healthy {
				c.newGaugeMetric(ch, c.totalMetrics["license_reporting_healthy"],
					float64(1))
			} else {
				c.newGaugeMetric(ch, c.totalMetrics["license_reporting_healthy"],
					float64(0))
			}
			c.newGaugeMetric(ch, c.totalMetrics["license_reporting_fails"],
				float64(license.Reporting.Fails))

			c.newGaugeMetric(ch, c.totalMetrics["license_reporting_grace_period"],
				float64(license.Reporting.Grace))
		}
	}

	for name, zone := range stats.ServerZones {
		labelValues := []string{name}
		varLabelValues := c.getServerZoneLabelValues(name)

		if c.variableLabelNames.ServerZoneVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.ServerZoneVariableLabelNames) {
			c.logger.Warn("wrong number of labels for http zone, empty labels will be used instead", "zone", name, "expected", len(c.variableLabelNames.ServerZoneVariableLabelNames), "got", len(varLabelValues))
			for range c.variableLabelNames.ServerZoneVariableLabelNames {
				labelValues = append(labelValues, "")
			}
		} else {
			labelValues = append(labelValues, varLabelValues...)
		}

		c.newGaugeMetric(ch, c.serverZoneMetrics["processing"],
			float64(zone.Processing), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["requests"], createdAt,
			float64(zone.Requests), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["responses_1xx"], createdAt,
			float64(zone.Responses.Responses1xx), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["responses_2xx"], createdAt,
			float64(zone.Responses.Responses2xx), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["responses_3xx"], createdAt,
			float64(zone.Responses.Responses3xx), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["responses_4xx"], createdAt,
			float64(zone.Responses.Responses4xx), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["responses_5xx"], createdAt,
			float64(zone.Responses.Responses5xx), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["discarded"], createdAt,
			float64(zone.Discarded), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["received"], createdAt,
			float64(zone.Received), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["sent"], createdAt,
			float64(zone.Sent), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_100"], createdAt,
			float64(zone.Responses.Codes.HTTPContinue), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_101"], createdAt,
			float64(zone.Responses.Codes.HTTPSwitchingProtocols), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_102"], createdAt,
			float64(zone.Responses.Codes.HTTPProcessing), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_200"], createdAt,
			float64(zone.Responses.Codes.HTTPOk), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_201"], createdAt,
			float64(zone.Responses.Codes.HTTPCreated), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_202"], createdAt,
			float64(zone.Responses.Codes.HTTPAccepted), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_204"], createdAt,
			float64(zone.Responses.Codes.HTTPNoContent), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_206"], createdAt,
			float64(zone.Responses.Codes.HTTPPartialContent), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_300"], createdAt,
			float64(zone.Responses.Codes.HTTPSpecialResponse), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_301"], createdAt,
			float64(zone.Responses.Codes.HTTPMovedPermanently), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_302"], createdAt,
			float64(zone.Responses.Codes.HTTPMovedTemporarily), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_303"], createdAt,
			float64(zone.Responses.Codes.HTTPSeeOther), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_304"], createdAt,
			float64(zone.Responses.Codes.HTTPNotModified), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_307"], createdAt,
			float64(zone.Responses.Codes.HTTPTemporaryRedirect), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_400"], createdAt,
			float64(zone.Responses.Codes.HTTPBadRequest), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_401"], createdAt,
			float64(zone.Responses.Codes.HTTPUnauthorized), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_403"], createdAt,
			float64(zone.Responses.Codes.HTTPForbidden), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_404"], createdAt,
			float64(zone.Responses.Codes.HTTPNotFound), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_405"], createdAt,
			float64(zone.Responses.Codes.HTTPNotAllowed), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_408"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestTimeOut), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_409"], createdAt,
			float64(zone.Responses.Codes.HTTPConflict), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_411"], createdAt,
			float64(zone.Responses.Codes.HTTPLengthRequired), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_412"], createdAt,
			float64(zone.Responses.Codes.HTTPPreconditionFailed), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_413"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestEntityTooLarge), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_414"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestURITooLarge), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_415"], createdAt,
			float64(zone.Responses.Codes.HTTPUnsupportedMediaType), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_416"], createdAt,
			float64(zone.Responses.Codes.HTTPRangeNotSatisfiable), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_429"], createdAt,
			float64(zone.Responses.Codes.HTTPTooManyRequests), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_444"], createdAt,
			float64(zone.Responses.Codes.HTTPClose), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_494"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestHeaderTooLarge), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_495"], createdAt,
			float64(zone.Responses.Codes.HTTPSCertError), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_496"], createdAt,
			float64(zone.Responses.Codes.HTTPSNoCert), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_497"], createdAt,
			float64(zone.Responses.Codes.HTTPToHTTPS), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_499"], createdAt,
			float64(zone.Responses.Codes.HTTPClientClosedRequest), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_500"], createdAt,
			float64(zone.Responses.Codes.HTTPInternalServerError), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_501"], createdAt,
			float64(zone.Responses.Codes.HTTPNotImplemented), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_502"], createdAt,
			float64(zone.Responses.Codes.HTTPBadGateway), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_503"], createdAt,
			float64(zone.Responses.Codes.HTTPServiceUnavailable), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_504"], createdAt,
			float64(zone.Responses.Codes.HTTPGatewayTimeOut), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["codes_507"], createdAt,
			float64(zone.Responses.Codes.HTTPInsufficientStorage), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["ssl_handshakes"], createdAt,
			float64(zone.SSL.Handshakes), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["ssl_handshakes_failed"], createdAt,
			float64(zone.SSL.HandshakesFailed), labelValues...)
		c.newCounterMetric(ch, c.serverZoneMetrics["ssl_session_reuses"], createdAt,
			float64(zone.SSL.SessionReuses), labelValues...)
	}

	for name, zone := range stats.StreamServerZones {
		labelValues := []string{name}
		varLabelValues := c.getStreamServerZoneLabelValues(name)

		if c.variableLabelNames.StreamServerZoneVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.StreamServerZoneVariableLabelNames) {
			c.logger.Warn("wrong number of labels for stream server zone, empty labels will be used instead", "zone", name, "expected", len(c.variableLabelNames.StreamServerZoneVariableLabelNames), "got", len(varLabelValues))
			for range c.variableLabelNames.StreamServerZoneVariableLabelNames {
				labelValues = append(labelValues, "")
			}
		} else {
			labelValues = append(labelValues, varLabelValues...)
		}
		c.newGaugeMetric(ch, c.streamServerZoneMetrics["processing"],
			float64(zone.Processing), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["connections"], createdAt,
			float64(zone.Connections), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["sessions_2xx"], createdAt,
			float64(zone.Sessions.Sessions2xx), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["sessions_4xx"], createdAt,
			float64(zone.Sessions.Sessions4xx), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["sessions_5xx"], createdAt,
			float64(zone.Sessions.Sessions5xx), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["discarded"], createdAt,
			float64(zone.Discarded), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["received"], createdAt,
			float64(zone.Received), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["sent"], createdAt,
			float64(zone.Sent), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["ssl_handshakes"], createdAt,
			float64(zone.SSL.Handshakes), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["ssl_handshakes_failed"], createdAt,
			float64(zone.SSL.HandshakesFailed), labelValues...)
		c.newCounterMetric(ch, c.streamServerZoneMetrics["ssl_session_reuses"], createdAt,
			float64(zone.SSL.SessionReuses), labelValues...)
	}

	for name, upstream := range stats.Upstreams {
		for _, peer := range upstream.Peers {
			labelValues := []string{name, peer.Server}
			varLabelValues := c.getUpstreamServerLabelValues(name)

			if c.variableLabelNames.UpstreamServerVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.UpstreamServerVariableLabelNames) {
				c.logger.Warn("wrong number of labels for upstream, empty labels will be used instead", "upstream", name, "expected", len(c.variableLabelNames.UpstreamServerVariableLabelNames), "got", len(varLabelValues))
				for range c.variableLabelNames.UpstreamServerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varLabelValues...)
			}

			upstreamServer := fmt.Sprintf("%v/%v", name, peer.Server)
			varPeerLabelValues := c.getUpstreamServerPeerLabelValues(upstreamServer)
			if c.variableLabelNames.UpstreamServerPeerVariableLabelNames != nil && len(varPeerLabelValues) != len(c.variableLabelNames.UpstreamServerPeerVariableLabelNames) {
				c.logger.Warn("wrong number of labels for upstream peer, empty labels will be used instead", "upstream", name, "peer", peer.Server, "expected", len(c.variableLabelNames.UpstreamServerPeerVariableLabelNames), "got", len(varPeerLabelValues))
				for range c.variableLabelNames.UpstreamServerPeerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varPeerLabelValues...)
			}

			c.newGaugeMetric(ch, c.upstreamServerMetrics["state"],
				upstreamServerStates[peer.State], labelValues...)
			c.newGaugeMetric(ch, c.upstreamServerMetrics["active"],
				float64(peer.Active), labelValues...)
			c.newGaugeMetric(ch, c.upstreamServerMetrics["limit"],
				float64(peer.MaxConns), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["requests"], createdAt,
				float64(peer.Requests), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["responses_1xx"], createdAt,
				float64(peer.Responses.Responses1xx), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["responses_2xx"], createdAt,
				float64(peer.Responses.Responses2xx), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["responses_3xx"], createdAt,
				float64(peer.Responses.Responses3xx), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["responses_4xx"], createdAt,
				float64(peer.Responses.Responses4xx), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["responses_5xx"], createdAt,
				float64(peer.Responses.Responses5xx), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["sent"], createdAt,
				float64(peer.Sent), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["received"], createdAt,
				float64(peer.Received), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["fails"], createdAt,
				float64(peer.Fails), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["unavail"], createdAt,
				float64(peer.Unavail), labelValues...)
			c.newGaugeMetric(ch, c.upstreamServerMetrics["header_time"],
				float64(peer.HeaderTime), labelValues...)
			c.newGaugeMetric(ch, c.upstreamServerMetrics["response_time"],
				float64(peer.ResponseTime), labelValues...)

			if peer.HealthChecks != (plusclient.HealthChecks{}) {
				c.newCounterMetric(ch, c.upstreamServerMetrics["health_checks_checks"], createdAt,
					float64(peer.HealthChecks.Checks), labelValues...)
				c.newCounterMetric(ch, c.upstreamServerMetrics["health_checks_fails"], createdAt,
					float64(peer.HealthChecks.Fails), labelValues...)
				c.newCounterMetric(ch, c.upstreamServerMetrics["health_checks_unhealthy"], createdAt,
					float64(peer.HealthChecks.Unhealthy), labelValues...)
			}
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_100"], createdAt,
				float64(peer.Responses.Codes.HTTPContinue), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_101"], createdAt,
				float64(peer.Responses.Codes.HTTPSwitchingProtocols), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_102"], createdAt,
				float64(peer.Responses.Codes.HTTPProcessing), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_200"], createdAt,
				float64(peer.Responses.Codes.HTTPOk), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_201"], createdAt,
				float64(peer.Responses.Codes.HTTPCreated), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_202"], createdAt,
				float64(peer.Responses.Codes.HTTPAccepted), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_204"], createdAt,
				float64(peer.Responses.Codes.HTTPNoContent), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_206"], createdAt,
				float64(peer.Responses.Codes.HTTPPartialContent), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_300"], createdAt,
				float64(peer.Responses.Codes.HTTPSpecialResponse), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_301"], createdAt,
				float64(peer.Responses.Codes.HTTPMovedPermanently), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_302"], createdAt,
				float64(peer.Responses.Codes.HTTPMovedTemporarily), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_303"], createdAt,
				float64(peer.Responses.Codes.HTTPSeeOther), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_304"], createdAt,
				float64(peer.Responses.Codes.HTTPNotModified), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_307"], createdAt,
				float64(peer.Responses.Codes.HTTPTemporaryRedirect), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_400"], createdAt,
				float64(peer.Responses.Codes.HTTPBadRequest), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_401"], createdAt,
				float64(peer.Responses.Codes.HTTPUnauthorized), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_403"], createdAt,
				float64(peer.Responses.Codes.HTTPForbidden), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_404"], createdAt,
				float64(peer.Responses.Codes.HTTPNotFound), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_405"], createdAt,
				float64(peer.Responses.Codes.HTTPNotAllowed), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_408"], createdAt,
				float64(peer.Responses.Codes.HTTPRequestTimeOut), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_409"], createdAt,
				float64(peer.Responses.Codes.HTTPConflict), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_411"], createdAt,
				float64(peer.Responses.Codes.HTTPLengthRequired), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_412"], createdAt,
				float64(peer.Responses.Codes.HTTPPreconditionFailed), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_413"], createdAt,
				float64(peer.Responses.Codes.HTTPRequestEntityTooLarge), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_414"], createdAt,
				float64(peer.Responses.Codes.HTTPRequestURITooLarge), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_415"], createdAt,
				float64(peer.Responses.Codes.HTTPUnsupportedMediaType), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_416"], createdAt,
				float64(peer.Responses.Codes.HTTPRangeNotSatisfiable), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_429"], createdAt,
				float64(peer.Responses.Codes.HTTPTooManyRequests), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_444"], createdAt,
				float64(peer.Responses.Codes.HTTPClose), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_494"], createdAt,
				float64(peer.Responses.Codes.HTTPRequestHeaderTooLarge), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_495"], createdAt,
				float64(peer.Responses.Codes.HTTPSCertError), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_496"], createdAt,
				float64(peer.Responses.Codes.HTTPSNoCert), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_497"], createdAt,
				float64(peer.Responses.Codes.HTTPToHTTPS), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_499"], createdAt,
				float64(peer.Responses.Codes.HTTPClientClosedRequest), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_500"], createdAt,
				float64(peer.Responses.Codes.HTTPInternalServerError), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_501"], createdAt,
				float64(peer.Responses.Codes.HTTPNotImplemented), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_502"], createdAt,
				float64(peer.Responses.Codes.HTTPBadGateway), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_503"], createdAt,
				float64(peer.Responses.Codes.HTTPServiceUnavailable), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_504"], createdAt,
				float64(peer.Responses.Codes.HTTPGatewayTimeOut), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["codes_507"], createdAt,
				float64(peer.Responses.Codes.HTTPInsufficientStorage), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["ssl_handshakes"], createdAt,
				float64(peer.SSL.Handshakes), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["ssl_handshakes_failed"], createdAt,
				float64(peer.SSL.HandshakesFailed), labelValues...)
			c.newCounterMetric(ch, c.upstreamServerMetrics["ssl_session_reuses"], createdAt,
				float64(peer.SSL.SessionReuses), labelValues...)
		}
		c.newGaugeMetric(ch, c.upstreamMetrics["keepalive"],
			float64(upstream.Keepalive), name)
		c.newGaugeMetric(ch, c.upstreamMetrics["zombies"],
			float64(upstream.Zombies), name)
	}

	for name, upstream := range stats.StreamUpstreams {
		for _, peer := range upstream.Peers {
			labelValues := []string{name, peer.Server}
			varLabelValues := c.getStreamUpstreamServerLabelValues(name)

			if c.variableLabelNames.StreamUpstreamServerVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.StreamUpstreamServerVariableLabelNames) {
				c.logger.Warn("wrong number of labels for stream server, empty labels will be used instead", "server", name, "labels", c.variableLabelNames.StreamUpstreamServerVariableLabelNames, "values", varLabelValues)
				for range c.variableLabelNames.StreamUpstreamServerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varLabelValues...)
			}

			upstreamServer := fmt.Sprintf("%v/%v", name, peer.Server)
			varPeerLabelValues := c.getStreamUpstreamServerPeerLabelValues(upstreamServer)
			if c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames != nil && len(varPeerLabelValues) != len(c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames) {
				c.logger.Warn("wrong number of labels for stream upstream peer, empty labels will be used instead", "server", upstreamServer, "labels", c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames, "values", varPeerLabelValues)
				for range c.variableLabelNames.StreamUpstreamServerPeerVariableLabelNames {
					labelValues = append(labelValues, "")
				}
			} else {
				labelValues = append(labelValues, varPeerLabelValues...)
			}

			c.newGaugeMetric(ch, c.streamUpstreamServerMetrics["state"],
				upstreamServerStates[peer.State], labelValues...)
			c.newGaugeMetric(ch, c.streamUpstreamServerMetrics["active"],
				float64(peer.Active), labelValues...)
			c.newGaugeMetric(ch, c.streamUpstreamServerMetrics["limit"],
				float64(peer.MaxConns), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["connections"], createdAt,
				float64(peer.Connections), labelValues...)
			c.newGaugeMetric(ch, c.streamUpstreamServerMetrics["connect_time"],
				float64(peer.ConnectTime), labelValues...)
			c.newGaugeMetric(ch, c.streamUpstreamServerMetrics["first_byte_time"],
				float64(peer.FirstByteTime), labelValues...)
			c.newGaugeMetric(ch, c.streamUpstreamServerMetrics["response_time"],
				float64(peer.ResponseTime), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["sent"], createdAt,
				float64(peer.Sent), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["received"], createdAt,
				float64(peer.Received), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["fails"], createdAt,
				float64(peer.Fails), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["unavail"], createdAt,
				float64(peer.Unavail), labelValues...)
			if peer.HealthChecks != (plusclient.HealthChecks{}) {
				c.newCounterMetric(ch, c.streamUpstreamServerMetrics["health_checks_checks"], createdAt,
					float64(peer.HealthChecks.Checks), labelValues...)
				c.newCounterMetric(ch, c.streamUpstreamServerMetrics["health_checks_fails"], createdAt,
					float64(peer.HealthChecks.Fails), labelValues...)
				c.newCounterMetric(ch, c.streamUpstreamServerMetrics["health_checks_unhealthy"], createdAt,
					float64(peer.HealthChecks.Unhealthy), labelValues...)
			}
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["ssl_handshakes"], createdAt,
				float64(peer.SSL.Handshakes), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["ssl_handshakes_failed"], createdAt,
				float64(peer.SSL.HandshakesFailed), labelValues...)
			c.newCounterMetric(ch, c.streamUpstreamServerMetrics["ssl_session_reuses"], createdAt,
				float64(peer.SSL.SessionReuses), labelValues...)
		}
		c.newGaugeMetric(ch, c.streamUpstreamMetrics["zombies"],
			float64(upstream.Zombies), name)
	}

	if stats.StreamZoneSync != nil {
		for name, zone := range stats.StreamZoneSync.Zones {
			c.newGaugeMetric(ch, c.streamZoneSyncMetrics["records_pending"],
				float64(zone.RecordsPending), name)
			c.newGaugeMetric(ch, c.streamZoneSyncMetrics["records_total"],
				float64(zone.RecordsTotal), name)
		}

		c.newCounterMetric(ch, c.streamZoneSyncMetrics["bytes_in"], createdAt,
			float64(stats.StreamZoneSync.Status.BytesIn))
		c.newCounterMetric(ch, c.streamZoneSyncMetrics["bytes_out"], createdAt,
			float64(stats.StreamZoneSync.Status.BytesOut))
		c.newCounterMetric(ch, c.streamZoneSyncMetrics["msgs_in"], createdAt,
			float64(stats.StreamZoneSync.Status.MsgsIn))
		c.newCounterMetric(ch, c.streamZoneSyncMetrics["msgs_out"], createdAt,
			float64(stats.StreamZoneSync.Status.MsgsOut))
		c.newGaugeMetric(ch, c.streamZoneSyncMetrics["nodes_online"],
			float64(stats.StreamZoneSync.Status.NodesOnline))
	}

	for name, zone := range stats.LocationZones {
		c.newCounterMetric(ch, c.locationZoneMetrics["requests"], createdAt,
			float64(zone.Requests), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["responses_1xx"], createdAt,
			float64(zone.Responses.Responses1xx), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["responses_2xx"], createdAt,
			float64(zone.Responses.Responses2xx), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["responses_3xx"], createdAt,
			float64(zone.Responses.Responses3xx), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["responses_4xx"], createdAt,
			float64(zone.Responses.Responses4xx), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["responses_5xx"], createdAt,
			float64(zone.Responses.Responses5xx), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["discarded"], createdAt,
			float64(zone.Discarded), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["received"], createdAt,
			float64(zone.Received), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["sent"], createdAt,
			float64(zone.Sent), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_100"], createdAt,
			float64(zone.Responses.Codes.HTTPContinue), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_101"], createdAt,
			float64(zone.Responses.Codes.HTTPSwitchingProtocols), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_102"], createdAt,
			float64(zone.Responses.Codes.HTTPProcessing), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_200"], createdAt,
			float64(zone.Responses.Codes.HTTPOk), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_201"], createdAt,
			float64(zone.Responses.Codes.HTTPCreated), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_202"], createdAt,
			float64(zone.Responses.Codes.HTTPAccepted), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_204"], createdAt,
			float64(zone.Responses.Codes.HTTPNoContent), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_206"], createdAt,
			float64(zone.Responses.Codes.HTTPPartialContent), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_300"], createdAt,
			float64(zone.Responses.Codes.HTTPSpecialResponse), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_301"], createdAt,
			float64(zone.Responses.Codes.HTTPMovedPermanently), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_302"], createdAt,
			float64(zone.Responses.Codes.HTTPMovedTemporarily), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_303"], createdAt,
			float64(zone.Responses.Codes.HTTPSeeOther), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_304"], createdAt,
			float64(zone.Responses.Codes.HTTPNotModified), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_307"], createdAt,
			float64(zone.Responses.Codes.HTTPTemporaryRedirect), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_400"], createdAt,
			float64(zone.Responses.Codes.HTTPBadRequest), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_401"], createdAt,
			float64(zone.Responses.Codes.HTTPUnauthorized), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_403"], createdAt,
			float64(zone.Responses.Codes.HTTPForbidden), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_404"], createdAt,
			float64(zone.Responses.Codes.HTTPNotFound), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_405"], createdAt,
			float64(zone.Responses.Codes.HTTPNotAllowed), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_408"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestTimeOut), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_409"], createdAt,
			float64(zone.Responses.Codes.HTTPConflict), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_411"], createdAt,
			float64(zone.Responses.Codes.HTTPLengthRequired), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_412"], createdAt,
			float64(zone.Responses.Codes.HTTPPreconditionFailed), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_413"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestEntityTooLarge), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_414"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestURITooLarge), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_415"], createdAt,
			float64(zone.Responses.Codes.HTTPUnsupportedMediaType), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_416"], createdAt,
			float64(zone.Responses.Codes.HTTPRangeNotSatisfiable), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_429"], createdAt,
			float64(zone.Responses.Codes.HTTPTooManyRequests), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_444"], createdAt,
			float64(zone.Responses.Codes.HTTPClose), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_494"], createdAt,
			float64(zone.Responses.Codes.HTTPRequestHeaderTooLarge), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_495"], createdAt,
			float64(zone.Responses.Codes.HTTPSCertError), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_496"], createdAt,
			float64(zone.Responses.Codes.HTTPSNoCert), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_497"], createdAt,
			float64(zone.Responses.Codes.HTTPToHTTPS), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_499"], createdAt,
			float64(zone.Responses.Codes.HTTPClientClosedRequest), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_500"], createdAt,
			float64(zone.Responses.Codes.HTTPInternalServerError), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_501"], createdAt,
			float64(zone.Responses.Codes.HTTPNotImplemented), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_502"], createdAt,
			float64(zone.Responses.Codes.HTTPBadGateway), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_503"], createdAt,
			float64(zone.Responses.Codes.HTTPServiceUnavailable), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_504"], createdAt,
			float64(zone.Responses.Codes.HTTPGatewayTimeOut), name)
		c.newCounterMetric(ch, c.locationZoneMetrics["codes_507"], createdAt,
			float64(zone.Responses.Codes.HTTPInsufficientStorage), name)
	}

	for name, zone := range stats.Resolvers {
		c.newCounterMetric(ch, c.resolverMetrics["name"], createdAt,
			float64(zone.Requests.Name), name)
		c.newCounterMetric(ch, c.resolverMetrics["srv"], createdAt,
			float64(zone.Requests.Srv), name)
		c.newCounterMetric(ch, c.resolverMetrics["addr"], createdAt,
			float64(zone.Requests.Addr), name)
		c.newCounterMetric(ch, c.resolverMetrics["noerror"], createdAt,
			float64(zone.Responses.Noerror), name)
		c.newCounterMetric(ch, c.resolverMetrics["formerr"], createdAt,
			float64(zone.Responses.Formerr), name)
		c.newCounterMetric(ch, c.resolverMetrics["servfail"], createdAt,
			float64(zone.Responses.Servfail), name)
		c.newCounterMetric(ch, c.resolverMetrics["nxdomain"], createdAt,
			float64(zone.Responses.Nxdomain), name)
		c.newCounterMetric(ch, c.resolverMetrics["notimp"], createdAt,
			float64(zone.Responses.Notimp), name)
		c.newCounterMetric(ch, c.resolverMetrics["refused"], createdAt,
			float64(zone.Responses.Refused), name)
		c.newCounterMetric(ch, c.resolverMetrics["timedout"], createdAt,
			float64(zone.Responses.Timedout), name)
		c.newCounterMetric(ch, c.resolverMetrics["unknown"], createdAt,
			float64(zone.Responses.Unknown), name)
	}

	for name, zone := range stats.HTTPLimitRequests {
		c.newCounterMetric(ch, c.limitRequestMetrics["passed"], createdAt, float64(zone.Passed), name)
		c.newCounterMetric(ch, c.limitRequestMetrics["rejected"], createdAt, float64(zone.Rejected), name)
		c.newCounterMetric(ch, c.limitRequestMetrics["delayed"], createdAt, float64(zone.Delayed), name)
		c.newCounterMetric(ch, c.limitRequestMetrics["rejected_dry_run"], createdAt, float64(zone.RejectedDryRun), name)
		c.newCounterMetric(ch, c.limitRequestMetrics["delayed_dry_run"], createdAt, float64(zone.DelayedDryRun), name)
	}

	for name, zone := range stats.HTTPLimitConnections {
		c.newCounterMetric(ch, c.limitConnectionMetrics["passed"], createdAt, float64(zone.Passed), name)
		c.newCounterMetric(ch, c.limitConnectionMetrics["rejected"], createdAt, float64(zone.Rejected), name)
		c.newCounterMetric(ch, c.limitConnectionMetrics["rejected_dry_run"], createdAt, float64(zone.RejectedDryRun), name)
	}

	for name, zone := range stats.StreamLimitConnections {
		c.newCounterMetric(ch, c.streamLimitConnectionMetrics["passed"], createdAt, float64(zone.Passed), name)
		c.newCounterMetric(ch, c.streamLimitConnectionMetrics["rejected"], createdAt, float64(zone.Rejected), name)
		c.newCounterMetric(ch, c.streamLimitConnectionMetrics["rejected_dry_run"], createdAt, float64(zone.RejectedDryRun), name)
	}

	for name, zone := range stats.Caches {
		labelValues := []string{name}
		varLabelValues := c.getCacheZoneLabelValues(name)

		if c.variableLabelNames.CacheZoneVariableLabelNames != nil && len(varLabelValues) != len(c.variableLabelNames.CacheZoneVariableLabelNames) {
			c.logger.Warn("wrong number of labels for cache zone, empty labels will be used instead", "zone", name, "labels", c.variableLabelNames.CacheZoneVariableLabelNames, "values", varLabelValues)
			for range c.variableLabelNames.CacheZoneVariableLabelNames {
				labelValues = append(labelValues, "")
			}
		} else {
			labelValues = append(labelValues, varLabelValues...)
		}

		c.newGaugeMetric(ch, c.cacheZoneMetrics["size"], float64(zone.Size), labelValues...)
		c.newGaugeMetric(ch, c.cacheZoneMetrics["max_size"], float64(zone.MaxSize), labelValues...)
		c.newGaugeMetric(ch, c.cacheZoneMetrics["cold"], booleanToFloat64[zone.Cold], labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["hit_responses"], createdAt, float64(zone.Hit.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["hit_bytes"], createdAt, float64(zone.Hit.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["stale_responses"], createdAt, float64(zone.Stale.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["stale_bytes"], createdAt, float64(zone.Stale.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["updating_responses"], createdAt, float64(zone.Updating.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["updating_bytes"], createdAt, float64(zone.Updating.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["revalidated_responses"], createdAt, float64(zone.Revalidated.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["revalidated_bytes"], createdAt, float64(zone.Revalidated.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["miss_responses"], createdAt, float64(zone.Miss.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["miss_bytes"], createdAt, float64(zone.Miss.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["expired_responses"], createdAt, float64(zone.Expired.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["expired_bytes"], createdAt, float64(zone.Expired.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["expired_responses_written"], createdAt, float64(zone.Expired.ResponsesWritten), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["expired_bytes_written"], createdAt, float64(zone.Expired.BytesWritten), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["bypass_responses"], createdAt, float64(zone.Bypass.Responses), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["bypass_bytes"], createdAt, float64(zone.Bypass.Bytes), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["bypass_responses_written"], createdAt, float64(zone.Bypass.ResponsesWritten), labelValues...)
		c.newCounterMetric(ch, c.cacheZoneMetrics["bypass_bytes_written"], createdAt, float64(zone.Bypass.BytesWritten), labelValues...)
	}

	for id, worker := range stats.Workers {
		workerID := strconv.FormatInt(int64(id), 10)
		workerPID := strconv.FormatUint(worker.ProcessID, 10)
		c.newCounterMetric(ch, c.workerMetrics["connection_accepted"], createdAt, float64(worker.Connections.Accepted), workerID, workerPID)
		c.newCounterMetric(ch, c.workerMetrics["connection_dropped"], createdAt, float64(worker.Connections.Dropped), workerID, workerPID)
		c.newGaugeMetric(ch, c.workerMetrics["connection_active"], float64(worker.Connections.Active), workerID, workerPID)
		c.newGaugeMetric(ch, c.workerMetrics["connection_idle"], float64(worker.Connections.Idle), workerID, workerPID)
		c.newCounterMetric(ch, c.workerMetrics["http_requests_total"], createdAt, float64(worker.HTTP.HTTPRequests.Total), workerID, workerPID)
		c.newGaugeMetric(ch, c.workerMetrics["http_requests_current"], float64(worker.HTTP.HTTPRequests.Current), workerID, workerPID)
	}
}

var upstreamServerStates = map[string]float64{
	"up":        1.0,
	"draining":  2.0,
	"down":      3.0,
	"unavail":   4.0,
	"checking":  5.0,
	"unhealthy": 6.0,
}

var booleanToFloat64 = map[bool]float64{
	true:  1.0,
	false: 0.0,
}

func newServerZoneMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"server_zone"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "server_zone", metricName), docString, labels, constLabels)
}

func newStreamServerZoneMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"server_zone"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_server_zone", metricName), docString, labels, constLabels)
}

func newUpstreamMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "upstream", metricName), docString, []string{"upstream"}, constLabels)
}

func newStreamUpstreamMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_upstream", metricName), docString, []string{"upstream"}, constLabels)
}

func newUpstreamServerMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"upstream", "server"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "upstream_server", metricName), docString, labels, constLabels)
}

func newStreamUpstreamServerMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"upstream", "server"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_upstream_server", metricName), docString, labels, constLabels)
}

func newStreamZoneSyncMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_zone_sync_status", metricName), docString, nil, constLabels)
}

func newStreamZoneSyncZoneMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_zone_sync_zone", metricName), docString, []string{"zone"}, constLabels)
}

func newLocationZoneMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "location_zone", metricName), docString, []string{"location_zone"}, constLabels)
}

func newResolverMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "resolver", metricName), docString, []string{"resolver"}, constLabels)
}

func newLimitRequestMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "limit_request", metricName), docString, []string{"zone"}, constLabels)
}

func newLimitConnectionMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "limit_connection", metricName), docString, []string{"zone"}, constLabels)
}

func newStreamLimitConnectionMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "stream_limit_connection", metricName), docString, []string{"zone"}, constLabels)
}

func newCacheZoneMetric(namespace string, metricName string, docString string, variableLabelNames []string, constLabels prometheus.Labels) *prometheus.Desc {
	labels := []string{"zone"}
	labels = append(labels, variableLabelNames...)
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "cache", metricName), docString, labels, constLabels)
}

func newWorkerMetric(namespace string, metricName string, docString string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "worker", metricName), docString, []string{"id", "pid"}, constLabels)
}
