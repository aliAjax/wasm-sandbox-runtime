package runtime

import (
	"crypto/sha256"
	"encoding/hex"
)

type Snapshot struct {
	ID       string
	Sequence uint64
	Digest   string
	State    []byte
}

func NewSnapshot(id string, seq uint64, state []byte) Snapshot {
	sum := sha256.Sum256(state)
	return Snapshot{ID: id, Sequence: seq, Digest: "sha256:" + hex.EncodeToString(sum[:]), State: append([]byte(nil), state...)}
}
func (s Snapshot) Verify() bool {
	sum := sha256.Sum256(s.State)
	return s.Digest == "sha256:"+hex.EncodeToString(sum[:])
}
