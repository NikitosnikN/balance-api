package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
	"time"
)

var once sync.Once
var metrics *Metrics
var maxAge = time.Minute
var objectives = map[float64]float64{
	0.5:  0.02,
	0.9:  0.01,
	0.95: 0.01,
	0.99: 0.001,
}

type Metrics struct {
	requestsTotal       *prometheus.CounterVec
	responseSize        *prometheus.SummaryVec
	responseDuration    *prometheus.SummaryVec
	nodePoolBlockHeight *prometheus.GaugeVec
	nodePoolLiveness    *prometheus.GaugeVec
}

func (m *Metrics) RequestsTotal() *prometheus.CounterVec {
	return m.requestsTotal
}

func (m *Metrics) ResponseSize() *prometheus.SummaryVec {
	return m.responseSize
}

func (m *Metrics) ResponseDuration() *prometheus.SummaryVec {
	return m.responseDuration
}

func (m *Metrics) NodePoolBlockHeight() *prometheus.GaugeVec {
	return m.nodePoolBlockHeight
}

func (m *Metrics) NodePoolLiveness() *prometheus.GaugeVec {
	return m.nodePoolLiveness
}

func NewMetrics() (http.HandlerFunc, *Metrics) {
	handler := promhttp.Handler().ServeHTTP

	once.Do(func() {
		httpLabels := []string{"status", "path"}
		poolLabels := []string{"name"}

		m := &Metrics{
			requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "requests_total",
				Help: "Number of requests",
			}, httpLabels),
			responseSize: prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Name:       "response_size",
				Help:       "response size in bytes",
				Objectives: objectives,
				MaxAge:     maxAge,
			}, httpLabels),
			responseDuration: prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Name:       "response_duration_seconds",
				Help:       "response size in seconds",
				Objectives: objectives,
				MaxAge:     maxAge,
			}, httpLabels),
			nodePoolBlockHeight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: "node_pool_block_height",
				Help: "Highest block of the node",
			}, poolLabels),
			nodePoolLiveness: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: "node_pool_liveness",
				Help: "Node liveness",
			}, poolLabels),
		}

		prometheus.MustRegister(m.requestsTotal)
		prometheus.MustRegister(m.responseSize)
		prometheus.MustRegister(m.responseDuration)
		prometheus.MustRegister(m.nodePoolBlockHeight)
		prometheus.MustRegister(m.nodePoolLiveness)

		metrics = m
	})

	return handler, metrics
}
