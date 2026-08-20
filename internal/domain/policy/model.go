package policy

import (
	"errors"
	"net"
	"strings"
)

type Policy struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenant_id"`
	Name          string        `json:"name"`
	Limits        Limits        `json:"limits"`
	Network       NetworkPolicy `json:"network"`
	EnvAllowlist  []string      `json:"env_allowlist"`
	Mounts        []Mount       `json:"mounts"`
	Deterministic bool          `json:"deterministic"`
	Revision      int64         `json:"revision"`
}
type Limits struct {
	CPU            uint64 `json:"cpu"`
	Memory         uint64 `json:"memory"`
	Output         uint64 `json:"output"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	Retries        int    `json:"retries"`
}
type NetworkPolicy struct {
	AllowEgress bool     `json:"allow_egress"`
	CIDRs       []string `json:"cidrs"`
	Domains     []string `json:"domains"`
}
type Mount struct {
	Name         string `json:"name"`
	ReadOnly     bool   `json:"read_only"`
	ObjectPrefix string `json:"object_prefix"`
}

func (p Policy) Validate() error {
	if p.ID == "" || p.TenantID == "" {
		return errors.New("policy identity required")
	}
	if p.Limits.CPU == 0 || p.Limits.Memory < 64*1024 || p.Limits.Output == 0 {
		return errors.New("policy limits invalid")
	}
	if p.Limits.TimeoutSeconds < 1 || p.Limits.TimeoutSeconds > 86400 {
		return errors.New("timeout invalid")
	}
	if p.Limits.Retries < 0 || p.Limits.Retries > 10 {
		return errors.New("retries invalid")
	}
	for _, c := range p.Network.CIDRs {
		if _, _, e := net.ParseCIDR(c); e != nil {
			return errors.New("invalid network cidr")
		}
	}
	for _, m := range p.Mounts {
		if !m.ReadOnly {
			return errors.New("mounts must be read-only")
		}
	}
	return nil
}
func (p Policy) AllowsDomain(host string) bool {
	if !p.Network.AllowEgress {
		return false
	}
	host = strings.ToLower(host)
	for _, d := range p.Network.Domains {
		if host == strings.ToLower(d) || strings.HasSuffix(host, "."+strings.ToLower(d)) {
			return true
		}
	}
	return false
}
func (p Policy) AllowsEnv(k string) bool {
	for _, v := range p.EnvAllowlist {
		if v == k {
			return true
		}
	}
	return false
}
func (p Policy) Compile() Compiled {
	c := Compiled{ID: p.ID, Deterministic: p.Deterministic, Env: map[string]struct{}{}}
	for _, v := range p.EnvAllowlist {
		c.Env[v] = struct{}{}
	}
	for _, d := range p.Network.Domains {
		c.Domains = append(c.Domains, strings.ToLower(d))
	}
	return c
}

type Compiled struct {
	ID            string
	Deterministic bool
	Env           map[string]struct{}
	Domains       []string
}
