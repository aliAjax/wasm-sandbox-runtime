package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

type Reader struct {
	r    io.Reader
	max  int64
	read int64
	hash [32]byte
	data []byte
}

func NewReader(r io.Reader, max int64) *Reader { return &Reader{r: r, max: max} }
func (r *Reader) ReadAll() ([]byte, string, error) {
	if r.max < 1 {
		return nil, "", errors.New("max bytes must be positive")
	}
	b, err := io.ReadAll(io.LimitReader(r.r, r.max+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(b)) > r.max {
		return nil, "", errors.New("artifact exceeds maximum")
	}
	r.data = b
	r.hash = sha256.Sum256(b)
	return b, "sha256:" + hex.EncodeToString(r.hash[:]), nil
}
func (r *Reader) Bytes() []byte { return append([]byte(nil), r.data...) }
