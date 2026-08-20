package quota

import "time"

type Window struct {
	Start      time.Time
	Duration   time.Duration
	CPU        float64
	Executions int
}

func (w Window) Contains(t time.Time) bool {
	return !t.Before(w.Start) && t.Before(w.Start.Add(w.Duration))
}
func (w Window) CPUAllowed(v float64) bool { return w.CPU <= 0 || v <= w.CPU }
func (w *Window) AddCPU(v float64)         { w.CPU += v }
