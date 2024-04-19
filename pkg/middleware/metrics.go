package middleware

import (
	"fmt"
	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type Metrics interface {
	RequestsTotal() *prometheus.CounterVec
	ResponseSize() *prometheus.SummaryVec
	ResponseDuration() *prometheus.SummaryVec
}

func MetricsMiddleware(metrics Metrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(next, w, r)

			go func() {
				select {
				case <-r.Context().Done():
					labels := []string{fmt.Sprintf("%d", m.Code), r.URL.Path}
					metrics.RequestsTotal().WithLabelValues(labels...).Inc()
					metrics.ResponseSize().WithLabelValues(labels...).Observe(float64(m.Written))
					metrics.ResponseDuration().WithLabelValues(labels...).Observe(m.Duration.Seconds())
				}
			}()
		}
		return http.HandlerFunc(fn)
	}
}
