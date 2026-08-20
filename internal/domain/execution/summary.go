package execution

import "time"

type Summary struct {
	ID           string
	TenantID     string
	State        State
	Attempts     int
	Duration     time.Duration
	ResultDigest string
	Error        string
}

func Summarize(e Execution) Summary {
	return Summary{ID: e.ID, TenantID: e.TenantID, State: e.State, Attempts: len(e.Attempts), Duration: e.Duration(), ResultDigest: e.ResultDigest, Error: e.Error}
}
func Summaries(items []Execution) []Summary {
	out := make([]Summary, 0, len(items))
	for _, e := range items {
		out = append(out, Summarize(e))
	}
	return out
}
