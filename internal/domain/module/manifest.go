package module

import (
	"sort"
	"strings"
)

type Manifest struct {
	ModuleID    string
	Version     int
	Entrypoint  string
	Imports     []string
	Exports     []string
	Permissions []Permission
}

func BuildManifest(v Version) Manifest {
	return Manifest{ModuleID: v.ModuleID, Version: v.Number, Entrypoint: v.Entrypoint, Imports: NormalizeImports(v.Imports), Exports: sorted(v.Exports), Permissions: v.EffectivePermissions()}
}
func sorted(in []string) []string { out := append([]string(nil), in...); sort.Strings(out); return out }
func (m Manifest) HasExport(name string) bool {
	for _, x := range m.Exports {
		if x == name {
			return true
		}
	}
	return false
}
func (m Manifest) CapabilityNames() []string {
	out := []string{}
	for _, p := range m.Permissions {
		out = append(out, strings.ToLower(string(p)))
	}
	return out
}
