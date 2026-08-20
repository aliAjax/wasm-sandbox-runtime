package inspector

import (
	"encoding/binary"
	"errors"
)

func Header(version uint32) []byte {
	b := make([]byte, 8)
	copy(b, "\x00asm")
	binary.LittleEndian.PutUint32(b[4:], version)
	return b
}
func IsCanonicalHeader(b []byte) bool {
	return len(b) >= 8 && string(b[:4]) == "\x00asm" && binary.LittleEndian.Uint32(b[4:]) == 1
}
func RequireCanonical(b []byte) error {
	if !IsCanonicalHeader(b) {
		return errors.New("wasm version 1 header required")
	}
	return nil
}
func SectionBounds(data []byte) []int {
	out := []int{}
	if len(data) < 8 {
		return out
	}
	s := scanner{data: data, pos: 8}
	for s.pos < len(data) {
		start := s.pos
		if _, e := s.byte(); e != nil {
			break
		}
		n, e := s.u32()
		if e != nil || n > uint32(len(data)-s.pos) {
			break
		}
		s.pos += int(n)
		out = append(out, start, s.pos)
	}
	return out
}
