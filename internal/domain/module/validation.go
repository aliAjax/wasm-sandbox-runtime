package module

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9-]{2,62}$`)

func ValidateName(name string) error {
	if !namePattern.MatchString(name) {
		return errors.New("name must be lowercase DNS label")
	}
	return nil
}
func ValidateTenant(tenant string) error {
	if strings.TrimSpace(tenant) == "" || len(tenant) > 128 {
		return errors.New("tenant identifier invalid")
	}
	return nil
}
func (v Version) ValidateDigest() error {
	if len(v.Digest) != 71 || !strings.HasPrefix(v.Digest, "sha256:") {
		return errors.New("version digest must use sha256")
	}
	if _, e := hex.DecodeString(v.Digest[7:]); e != nil {
		return errors.New("version digest is not hex")
	}
	return nil
}
func (v Version) HasPermission(p Permission) bool {
	for _, x := range v.Permissions {
		if x == p {
			return true
		}
	}
	return false
}
func (v Version) EffectivePermissions() []Permission {
	seen := map[Permission]bool{}
	out := make([]Permission, 0, len(v.Permissions))
	for _, p := range v.Permissions {
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}
func (m Module) Version(number int) (Version, bool) {
	for _, v := range m.Versions {
		if v.Number == number {
			return v, true
		}
	}
	return Version{}, false
}
func (m Module) Latest() (Version, error) {
	if len(m.Versions) == 0 {
		return Version{}, errors.New("no versions")
	}
	v := m.Versions[0]
	for _, x := range m.Versions[1:] {
		if x.Number > v.Number {
			v = x
		}
	}
	return v, nil
}
