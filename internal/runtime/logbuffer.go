package runtime

import "sync"

type LogBuffer struct {
	mu      sync.Mutex
	max     int
	lines   []string
	dropped uint64
}

func NewLogBuffer(max int) *LogBuffer {
	if max < 1 {
		max = 100
	}
	return &LogBuffer{max: max}
}
func (b *LogBuffer) Add(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.lines) >= b.max {
		b.lines = b.lines[1:]
		b.dropped++
	}
	b.lines = append(b.lines, line)
}
func (b *LogBuffer) Lines() ([]string, uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.lines...), b.dropped
}
