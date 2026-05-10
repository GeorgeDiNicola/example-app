package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestMetricsRecordsLatencyHistogram(t *testing.T) {
	const route = "/test"

	histogramBefore := histogramValue(t, httpRequestDurationSeconds, route)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, route, nil)
	recorder := httptest.NewRecorder()

	// Wrap the handler with the metrics middleware and run it
	Metrics(route, handler).ServeHTTP(recorder, req)

	histogramAfter := histogramValue(t, httpRequestDurationSeconds, route)

	if histogramAfter.GetSampleCount() != histogramBefore.GetSampleCount()+1 {
		t.Fatalf("expected sample count to increase by 1")
	}
}

func TestMetricsRecordsRequestCounterIncrease(t *testing.T) {
	const route = "/error"

	statusStr := strconv.Itoa(http.StatusInternalServerError)
	totalBefore := counterValue(t, httpRequestsTotal, route, statusStr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, route, nil)
	recorder := httptest.NewRecorder()

	// Wrap the handler with the metrics middleware and run it
	Metrics(route, handler).ServeHTTP(recorder, req)

	totalAfter := counterValue(t, httpRequestsTotal, route, statusStr)
	if totalAfter != totalBefore+1 {
		t.Fatalf("unexpected total request count: got %v want %v", totalAfter, totalBefore+1)
	}
}

func counterValue(t *testing.T, collector *prometheus.CounterVec, labels ...string) float64 {
	t.Helper()

	metric, err := collector.GetMetricWithLabelValues(labels...)
	if err != nil {
		t.Fatalf("failed to read counter metric: %v", err)
	}

	dtoMetric := &dto.Metric{}
	if err := metric.Write(dtoMetric); err != nil {
		t.Fatalf("failed to collect counter metric: %v", err)
	}

	return dtoMetric.GetCounter().GetValue()
}

func histogramValue(t *testing.T, collector *prometheus.HistogramVec, labels ...string) *dto.Histogram {
	t.Helper()

	observer, err := collector.GetMetricWithLabelValues(labels...)
	if err != nil {
		t.Fatalf("failed to read histogram metric: %v", err)
	}

	exportable, ok := observer.(prometheus.Metric)
	if !ok {
		t.Fatal("metric does not support data export")
	}

	metricDTO := &dto.Metric{}
	if err := exportable.Write(metricDTO); err != nil {
		t.Fatalf("failed to collect metric state: %v", err)
	}

	histogram := metricDTO.GetHistogram()
	if histogram == nil {
		t.Fatal("collected metric is not a histogram")
	}

	return histogram
}
