package runtime

import (
	"context"
	"errors"
	"io"
)

type OutputStream struct {
	writer  io.Writer
	max     int64
	written int64
}

func NewOutputStream(w io.Writer, max int64) *OutputStream { return &OutputStream{writer: w, max: max} }
func (s *OutputStream) Write(p []byte) (int, error) {
	if s.written+int64(len(p)) > s.max {
		return 0, errors.New("output stream limit exceeded")
	}
	n, e := s.writer.Write(p)
	s.written += int64(n)
	return n, e
}
func (s *OutputStream) Written() int64 { return s.written }
func (s *OutputStream) WriteContext(ctx context.Context, p []byte) (int, error) {
	return s.Write(p)
}
