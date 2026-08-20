package inspector

import "strings"

type CapabilityPolicy struct {
	Modules   map[string]bool
	Functions map[string]bool
}

func (p CapabilityPolicy) Allows(module, name string) bool {
	if p.Modules[module] {
		return true
	}
	return p.Functions[strings.ToLower(module+"."+name)]
}
func (p CapabilityPolicy) Denied(info ModuleInfo) []Import {
	out := []Import{}
	for _, im := range info.Imports {
		if !p.Allows(im.Module, im.Name) {
			out = append(out, im)
		}
	}
	return out
}
