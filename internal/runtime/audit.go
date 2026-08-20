package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type AuditRecord struct {
	Sequence    uint64
	ExecutionID string
	Action      string
	Digest      string
	At          time.Time
	Previous    string
	Hash        string
}
type AuditChain struct {
	mu      sync.Mutex
	records []AuditRecord
	head    string
}

func (a *AuditChain) Append(execution, action, digest string, at time.Time) AuditRecord {
	a.mu.Lock()
	defer a.mu.Unlock()
	r := AuditRecord{Sequence: uint64(len(a.records) + 1), ExecutionID: execution, Action: action, Digest: digest, At: at, Previous: a.head}
	sum := sha256.Sum256([]byte(r.Previous + execution + action + digest + at.UTC().String()))
	r.Hash = "sha256:" + hex.EncodeToString(sum[:])
	a.head = r.Hash
	a.records = append(a.records, r)
	return r
}
func (a *AuditChain) Verify() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	prev := a.head
	for _, r := range a.records {
		sum := sha256.Sum256([]byte(prev + r.ExecutionID + r.Action + r.Digest + r.At.UTC().String()))
		if "sha256:"+hex.EncodeToString(sum[:]) != r.Hash {
			return false
		}
		prev = r.Hash
	}
	return true
}
func (a *AuditChain) Records() []AuditRecord {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]AuditRecord(nil), a.records...)
}
func (a *AuditChain) StableRecords() []AuditRecord { return a.records }
