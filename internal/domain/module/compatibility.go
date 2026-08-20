package module

import "errors"

type Compatibility struct {
	Compatible bool
	Breaking   []string
	Warnings   []string
}

func Compare(a, b Version) Compatibility {
	c := Compatibility{Compatible: true}
	if a.Entrypoint != b.Entrypoint {
		c.Compatible = false
		c.Breaking = append(c.Breaking, "entrypoint changed")
	}
	if a.Schema.Output != b.Schema.Output {
		c.Warnings = append(c.Warnings, "output schema changed")
	}
	if len(b.Permissions) > len(a.Permissions) {
		c.Warnings = append(c.Warnings, "new capabilities requested")
	}
	return c
}
func RequireCompatible(a, b Version) error {
	c := Compare(a, b)
	if !c.Compatible {
		return errors.New("incompatible module version")
	}
	return nil
}
