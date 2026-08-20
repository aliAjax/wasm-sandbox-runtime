package artifact

import (
	"encoding/json"
	"strings"
)

func RedactJSON(data []byte, keys []string) []byte {
	var v any
	if json.Unmarshal(data, &v) != nil {
		return []byte("[REDACTED]")
	}
	redact(v, keys)
	out, _ := json.Marshal(v)
	return out
}
func redact(v any, keys []string) {
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			for _, target := range keys {
				if strings.EqualFold(k, target) {
					x[k] = "[REDACTED]"
					break
				}
			}
			redact(val, keys)
		}
	case []any:
		for _, item := range x {
			redact(item, keys)
		}
	}
}
