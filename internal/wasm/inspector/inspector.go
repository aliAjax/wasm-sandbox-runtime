package inspector

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrMalformed = errors.New("malformed wasm binary")

type ModuleInfo struct {
	Version               uint32   `json:"version"`
	Imports               []Import `json:"imports"`
	Exports               []Export `json:"exports"`
	MemoryPages           uint32   `json:"memory_pages"`
	Functions             uint32   `json:"functions"`
	EstimatedInstructions uint64   `json:"estimated_instructions"`
	CustomSections        []string `json:"custom_sections"`
}
type Import struct {
	Module string `json:"module"`
	Name   string `json:"name"`
	Kind   byte   `json:"kind"`
}
type Export struct {
	Name  string `json:"name"`
	Kind  byte   `json:"kind"`
	Index uint32 `json:"index"`
}
type scanner struct {
	data []byte
	pos  int
}

func Inspect(data []byte) (ModuleInfo, error) {
	if len(data) < 8 || string(data[:4]) != "\x00asm" {
		return ModuleInfo{}, ErrMalformed
	}
	s := scanner{data: data, pos: 8}
	info := ModuleInfo{Version: binary.LittleEndian.Uint32(data[4:8])}
	for s.pos < len(data) {
		id, e := s.byte()
		if e != nil {
			return info, e
		}
		n, e := s.u32()
		if e != nil || n > uint32(len(data)-s.pos) {
			return info, ErrMalformed
		}
		section := s.data[s.pos : s.pos+int(n)]
		s.pos += int(n)
		if id == 0 {
			if name := customName(section); name != "" {
				info.CustomSections = append(info.CustomSections, name)
			}
		} else if e := parseSection(id, section, &info); e != nil {
			return info, e
		}
	}
	if info.EstimatedInstructions > 1e12 {
		return info, errors.New("instruction estimate exceeds safety bound")
	}
	return info, nil
}
func (s *scanner) byte() (byte, error) {
	if s.pos >= len(s.data) {
		return 0, io.ErrUnexpectedEOF
	}
	v := s.data[s.pos]
	s.pos++
	return v, nil
}
func (s *scanner) u32() (uint32, error) {
	var v uint32
	var shift uint
	for i := 0; i < 5; i++ {
		b, e := s.byte()
		if e != nil {
			return 0, e
		}
		v |= uint32(b&0x7f) << shift
		if b&0x80 == 0 {
			return v, nil
		}
		shift += 7
	}
	return 0, ErrMalformed
}
func (s *scanner) name() (string, error) {
	n, e := s.u32()
	if e != nil || n > uint32(len(s.data)-s.pos) {
		return "", ErrMalformed
	}
	v := string(s.data[s.pos : s.pos+int(n)])
	s.pos += int(n)
	return v, nil
}
func customName(b []byte) string {
	s := scanner{data: b}
	n, e := s.name()
	if e != nil {
		return ""
	}
	return n
}
func parseSection(id byte, b []byte, i *ModuleInfo) error {
	s := scanner{data: b}
	switch id {
	case 0:
	case 1:
		n, e := s.u32()
		if e != nil {
			return e
		}
		for x := uint32(0); x < n; x++ {
			if _, e = s.byte(); e != nil {
				return e
			}
			if _, e = s.byte(); e != nil {
				return e
			}
		}
	case 2:
		n, e := s.u32()
		if e != nil {
			return e
		}
		for x := uint32(0); x < n; x++ {
			m, e := s.name()
			if e != nil {
				return e
			}
			name, e := s.name()
			if e != nil {
				return e
			}
			k, e := s.byte()
			if e != nil {
				return e
			}
			i.Imports = append(i.Imports, Import{Module: m, Name: name, Kind: k})
			if e = skipImportType(&s, k); e != nil {
				return e
			}
		}
	case 3:
		n, e := s.u32()
		if e != nil {
			return e
		}
		i.Functions = n
		for x := uint32(0); x < n; x++ {
			if _, e = s.u32(); e != nil {
				return e
			}
		}
	case 5:
		n, e := s.u32()
		if e != nil {
			return e
		}
		if n > 0 {
			flags, e := s.byte()
			if e != nil {
				return e
			}
			min, e := s.u32()
			if e != nil {
				return e
			}
			i.MemoryPages = min
			if flags&1 != 0 {
				if _, e = s.u32(); e != nil {
					return e
				}
			}
		}
	case 7:
		n, e := s.u32()
		if e != nil {
			return e
		}
		for x := uint32(0); x < n; x++ {
			name, e := s.name()
			if e != nil {
				return e
			}
			k, e := s.byte()
			if e != nil {
				return e
			}
			idx, e := s.u32()
			if e != nil {
				return e
			}
			i.Exports = append(i.Exports, Export{Name: name, Kind: k, Index: idx})
		}
	case 10:
		n, e := s.u32()
		if e != nil {
			return e
		}
		for x := uint32(0); x < n; x++ {
			sz, e := s.u32()
			if e != nil || sz > uint32(len(s.data)-s.pos) {
				return ErrMalformed
			}
			body := s.data[s.pos : s.pos+int(sz)]
			s.pos += int(sz)
			if len(body) > 0 {
				i.EstimatedInstructions += uint64(len(body))
			}
		}
	}
	return nil
}
func skipImportType(s *scanner, k byte) error {
	switch k {
	case 0:
		return skipFunc(s)
	case 1:
		return skipTable(s)
	case 2:
		return skipMemory(s)
	case 3:
		_, e := s.byte()
		if e != nil {
			return e
		}
		_, e = s.byte()
		return e
	default:
		return fmt.Errorf("unknown import kind %d", k)
	}
}
func skipFunc(s *scanner) error { _, e := s.u32(); return e }
func skipTable(s *scanner) error {
	if _, e := s.byte(); e != nil {
		return e
	}
	return skipMemoryLimits(s)
}
func skipMemory(s *scanner) error { return skipMemoryLimits(s) }
func skipMemoryLimits(s *scanner) error {
	f, e := s.byte()
	if e != nil {
		return e
	}
	if _, e = s.u32(); e != nil {
		return e
	}
	if f&1 != 0 {
		_, e = s.u32()
	}
	return e
}
func ValidateCapabilities(info ModuleInfo, allowed map[string]bool) error {
	for _, im := range info.Imports {
		key := strings.ToLower(im.Module + "." + im.Name)
		if !allowed[key] && !allowed[im.Module] {
			return fmt.Errorf("import capability denied: %s", key)
		}
	}
	return nil
}
func Entrypoint(info ModuleInfo, name string) error {
	for _, e := range info.Exports {
		if e.Name == name && e.Kind == 0 {
			return nil
		}
	}
	return fmt.Errorf("function export %q not found", name)
}
