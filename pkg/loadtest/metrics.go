package loadtest

import (
	"math"
	"sort"
	"time"
)

// Result represents a single request result
type Result struct {
	Latency    time.Duration
	Success    bool
	StatusCode int
}

// Metrics tracks and calculates performance metrics for the load test
type Metrics struct {
	latencies     []time.Duration
	failures      int
	totalReqs     int
	startTime     time.Time
	endTime       time.Time
	statusCodeMap map[int]int
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		latencies:     make([]time.Duration, 0, 1000),
		startTime:     time.Now(),
		statusCodeMap: make(map[int]int),
	}
}

// Start records the start time of the test
func (m *Metrics) Start() {
	m.startTime = time.Now()
}

// End records the end time of the test
func (m *Metrics) End() {
	m.endTime = time.Now()
}

// AddResult adds a result to the metrics
func (m *Metrics) AddResult(result Result) {
	m.totalReqs++
	m.statusCodeMap[result.StatusCode]++

	if result.Success {
		m.latencies = append(m.latencies, result.Latency)
	} else {
		m.failures++
	}
}

// Results calculates and returns the test results
func (m *Metrics) Results() TestResults {
	var avgLatency, p95Latency, p99Latency float64
	hadFailures := m.failures > 0
	totalReqs := m.totalReqs
	duration := m.endTime.Sub(m.startTime).Seconds()

	// Calculate latency statistics if we have successful requests
	if len(m.latencies) > 0 {
		// Calculate average latency
		var sum time.Duration
		for _, lat := range m.latencies {
			sum += lat
		}
		avgLatency = float64(sum) / float64(len(m.latencies)) / float64(time.Millisecond)

		// Sort latencies for percentile calculations
		sort.Slice(m.latencies, func(i, j int) bool {
			return m.latencies[i] < m.latencies[j]
		})

		// Calculate percentiles
		if len(m.latencies) > 1 {
			p95Index := int(math.Ceil(float64(len(m.latencies))*0.95)) - 1
			p99Index := int(math.Ceil(float64(len(m.latencies))*0.99)) - 1

			if p95Index >= 0 && p95Index < len(m.latencies) {
				p95Latency = float64(m.latencies[p95Index]) / float64(time.Millisecond)
			}

			if p99Index >= 0 && p99Index < len(m.latencies) {
				p99Latency = float64(m.latencies[p99Index]) / float64(time.Millisecond)
			}
		}
	}

	var successRate float64
	if totalReqs > 0 {
		successRate = float64(totalReqs-m.failures) / float64(totalReqs) * 100
	}

	return TestResults{
		TotalRequests:  totalReqs,
		FailedRequests: m.failures,
		AvgLatency:     avgLatency,
		P95Latency:     p95Latency,
		P99Latency:     p99Latency,
		HadFailures:    hadFailures,
		TestDuration:   duration,
		SuccessRate:    successRate,
		StatusCodes:    m.statusCodeMap,
	}
}

// TestResults contains the calculated metrics from a load test
type TestResults struct {
	TotalRequests  int
	FailedRequests int
	AvgLatency     float64
	P95Latency     float64
	P99Latency     float64
	HadFailures    bool
	TestDuration   float64
	SuccessRate    float64
	StatusCodes    map[int]int
}
