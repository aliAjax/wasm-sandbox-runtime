package module

func (m Module) Tenant() string      { return m.TenantID }
func (m Module) VersionCount() int   { return len(m.Versions) }
func (m Module) IsQuarantined() bool { return m.State == Quarantined }
func (v Version) ImportCount() int   { return len(v.Imports) }
func (v Version) ExportCount() int   { return len(v.Exports) }
