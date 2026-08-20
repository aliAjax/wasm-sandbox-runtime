package config

import "strings"

type SecretRef struct {
	Name  string
	Value string
	From  string
}

func (s SecretRef) Resolve(env map[string]string) string {
	if s.From != "" {
		return env[s.From]
	}
	return s.Value
}
func (s SecretRef) ResolveOK(env map[string]string) (string, bool) {
	if s.From == "" {
		return s.Value, s.Value != ""
	}
	v, ok := env[s.From]
	return v, ok && v != ""
}
func IsSensitive(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, "password") || strings.Contains(n, "secret") || strings.Contains(n, "token") || strings.Contains(n, "private")
}
func Redact(values map[string]string) map[string]string {
	out := make(map[string]string, len(values))
	for k, v := range values {
		if IsSensitive(k) {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
