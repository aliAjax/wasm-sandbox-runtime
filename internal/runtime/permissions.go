package runtime

import "errors"

type PermissionSet map[string]struct{}

func NewPermissions(values []string) PermissionSet {
	p := PermissionSet{}
	for _, v := range values {
		p[v] = struct{}{}
	}
	return p
}
func (p PermissionSet) Require(name string) error {
	if _, ok := p[name]; !ok {
		return errors.New("permission denied: " + name)
	}
	return nil
}
func (p PermissionSet) Has(name string) bool { _, ok := p[name]; return ok }
