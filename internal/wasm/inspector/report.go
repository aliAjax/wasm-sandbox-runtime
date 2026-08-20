package inspector

import "time"

type Report struct {
	Info          ModuleInfo
	Digest        string
	Allowed       bool
	DeniedImports []Import
	Warnings      []string
	GeneratedAt   time.Time
}

func BuildReport(info ModuleInfo, digest string, p CapabilityPolicy) Report {
	r := Report{Info: info, Digest: digest, GeneratedAt: time.Now().UTC()}
	r.DeniedImports = p.Denied(info)
	r.Allowed = len(r.DeniedImports) == 0
	if info.MemoryPages > 65536 {
		r.Warnings = append(r.Warnings, "memory exceeds common host limit")
	}
	if len(info.CustomSections) == 0 {
		r.Warnings = append(r.Warnings, "no build metadata custom section")
	}
	return r
}
