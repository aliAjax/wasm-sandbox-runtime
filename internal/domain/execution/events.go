package execution

import "time"

type LogEntry struct {
	Sequence uint64
	Level    string
	Message  string
	Fields   map[string]string
	At       time.Time
}

func NewLog(seq uint64, level, message string, at time.Time) LogEntry {
	return LogEntry{Sequence: seq, Level: level, Message: message, Fields: map[string]string{}, At: at}
}
func (l *LogEntry) Field(k, v string) { l.Fields[k] = v }
func (e Execution) Timeline() []LogEntry {
	out := []LogEntry{}
	for i, a := range e.Attempts {
		out = append(out, NewLog(uint64(i*2+1), "info", "attempt started", a.StartedAt))
		if a.FinishedAt != nil {
			out = append(out, NewLog(uint64(i*2+2), "info", "attempt finished", *a.FinishedAt))
		}
	}
	return out
}
