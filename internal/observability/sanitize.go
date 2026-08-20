package observability

import "strings"

func SanitizeFields(fields map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range fields {
		lk := strings.ToLower(k)
		if strings.Contains(lk, "token") || strings.Contains(lk, "secret") || strings.Contains(lk, "password") {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
func SafeMessage(message string) string {
	if len(message) > 1024 {
		return message[:1024] + "..."
	}
	return message
}
