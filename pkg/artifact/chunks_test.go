package artifact

import (
	"bytes"
	"testing"
)

func TestManifestReadIsolated(t *testing.T) {
	s := NewStore()
	m := s.Put([]byte("abcdef"), 2)
	got, _ := s.Manifest(m.Digest)
	got.Chunks[0].Data[0] = 'X'
	again, _ := s.Manifest(m.Digest)
	if again.Chunks[0].Data[0] == 'X' {
		t.Fatal("manifest aliased")
	}
}
func TestCloneManifestIsolated(t *testing.T) {
	m := BuildManifest([]byte("abcdef"), 2)
	c := CloneManifest(m)
	c.Chunks[0].Data[0] = 'X'
	if m.Chunks[0].Data[0] == 'X' {
		t.Fatal("clone aliased")
	}
}
func TestReaderSnapshotIsolated(t *testing.T) {
	r := NewReader(bytes.NewBufferString("abc"), 10)
	_, _, _ = r.ReadAll()
	v := r.Snapshot()
	v[0] = 'X'
	if r.Bytes()[0] == 'X' {
		t.Fatal("reader snapshot aliased")
	}
}
func TestChunksCopyIsolated(t *testing.T) {
	m := BuildManifest([]byte("abcdef"), 2)
	v := m.ChunksCopy()
	v[0].Data[0] = 'X'
	if m.Chunks[0].Data[0] == 'X' {
		t.Fatal("chunk copy aliased")
	}
}
