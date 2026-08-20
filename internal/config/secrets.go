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
func IsSensitive(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, "password") || strings.Contains(n, "secret") || strings.Contains(n, "token") || strings.Contains(n, "private")
}
func Redact(values map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range values {
		if IsSensitive(k) {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
