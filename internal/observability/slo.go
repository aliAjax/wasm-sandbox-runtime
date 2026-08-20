package observability

import "time"

type SLO struct {
	Availability float64
	LatencyP99   time.Duration
	ErrorBudget  float64
}

func DefaultSLO() SLO {
	return SLO{Availability: 0.999, LatencyP99: 500 * time.Millisecond, ErrorBudget: 0.001}
}
func (s SLO) Valid() bool {
	return s.Availability > 0 && s.Availability <= 1 && s.LatencyP99 > 0 && s.ErrorBudget >= 0 && s.ErrorBudget < 1
}
func (s SLO) Within(availability float64, latency time.Duration) bool {
	return availability >= s.Availability && latency <= s.LatencyP99
}
