package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func SHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}
func Verify(data []byte, digest string) bool { return SHA256(data) == digest }
func ValidateDigest(d string) error {
	if len(d) != 71 || d[:7] != "sha256:" {
		return errors.New("digest must be sha256:<64 hex>")
	}
	if _, e := hex.DecodeString(d[7:]); e != nil {
		return errors.New("digest contains invalid hex")
	}
	return nil
}

type Chunk struct {
	Index  int
	Digest string
	Size   int
	Data   []byte
}
type Manifest struct {
	Digest string
	Size   int64
	Chunks []Chunk
}

func BuildManifest(data []byte, chunkSize int) Manifest {
	if chunkSize < 1 {
		chunkSize = 1 << 20
	}
	m := Manifest{Digest: SHA256(data), Size: int64(len(data))}
	for i, off := 0, 0; off < len(data); i, off = i+1, off+chunkSize {
		end := off + chunkSize
		if end > len(data) {
			end = len(data)
		}
		part := append([]byte(nil), data[off:end]...)
		m.Chunks = append(m.Chunks, Chunk{Index: i, Digest: SHA256(part), Size: len(part), Data: part})
	}
	return m
}
func CloneManifest(m Manifest) Manifest { return m }
