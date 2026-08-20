package module

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func SignaturePayload(v Version) []byte {
	sum := sha256.Sum256([]byte(v.ModuleID + v.Digest + v.Entrypoint))
	return []byte(hex.EncodeToString(sum[:]))
}
func VerifySignature(v Version, trusted map[string]string) error {
	if v.Signature.KeyID == "" || v.Signature.Value == "" {
		return errors.New("signature required")
	}
	expected, ok := trusted[v.Signature.KeyID]
	if !ok || expected != v.Signature.Value {
		return errors.New("signature verification failed")
	}
	return nil
}
func (s Signature) Valid() bool {
	return s.Algorithm != "" && s.KeyID != "" && s.Value != "" && s.Verified
}
