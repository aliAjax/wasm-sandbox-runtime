package policy

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

type Compilation struct {
	Compiled Compiled
	Digest   string
	Warnings []string
}

func Compile(p Policy) Compilation {
	p.Network.Domains = NormalizeDomains(p.Network.Domains)
	sort.Strings(p.EnvAllowlist)
	c := p.Compile()
	sum := sha256.Sum256([]byte(p.ID + string(rune(p.Revision))))
	return Compilation{Compiled: c, Digest: "sha256:" + hex.EncodeToString(sum[:])}
}
func (c Compilation) Stable() bool { return c.Digest != "" && c.Compiled.ID != "" }
