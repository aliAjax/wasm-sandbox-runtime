package inspector

import (
	"crypto/sha256"
	"encoding/hex"
)

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}
func DigestEqual(data []byte, want string) bool { return Digest(data) == want }
func DigestParts(parts ...[]byte) string {
	h := sha256.New()
	for _, p := range parts {
		_, _ = h.Write(p)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
