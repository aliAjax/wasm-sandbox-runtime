package execution

import (
	"errors"
	"strings"
)

type Result struct {
	Digest      string `json:"digest"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Redacted    bool   `json:"redacted"`
}

func (r Result) Validate() error {
	if !strings.HasPrefix(r.Digest, "sha256:") {
		return errors.New("result digest required")
	}
	if r.Size < 0 {
		return errors.New("result size invalid")
	}
	return nil
}
func (e Execution) HasResult() bool { return e.ResultDigest != "" && e.State == Succeeded }
