package runtime

import (
	"crypto/sha256"
	"encoding/hex"
)

type Determinism struct {
	Enabled        bool
	RuntimeVersion string
	ModuleDigest   string
	InputDigest    string
}

func (d Determinism) Key() string {
	sum := sha256.Sum256([]byte(d.RuntimeVersion + "|" + d.ModuleDigest + "|" + d.InputDigest))
	return "det:" + hex.EncodeToString(sum[:])
}
func (d Determinism) Compatible(other Determinism) bool {
	return d.Enabled && other.Enabled && d.RuntimeVersion == other.RuntimeVersion && d.ModuleDigest == other.ModuleDigest && d.InputDigest == other.InputDigest
}
