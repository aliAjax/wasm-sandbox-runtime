package policy

import (
	"net"
	"strings"
)

func (p Policy) AllowsIP(ip net.IP) bool {
	if !p.Network.AllowEgress {
		return false
	}
	for _, raw := range p.Network.CIDRs {
		_, n, e := net.ParseCIDR(raw)
		if e == nil && n.Contains(ip) {
			return true
		}
	}
	return false
}
func NormalizeDomains(in []string) []string {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, d := range in {
		d = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(d), "."))
		if d != "" {
			if _, ok := seen[d]; !ok {
				seen[d] = struct{}{}
				out = append(out, d)
			}
		}
	}
	return out
}
