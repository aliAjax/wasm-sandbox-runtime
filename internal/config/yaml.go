package config

import (
	"errors"
	"os"
	"strings"
)

func LoadDotEnv(path string) error {
	b, e := os.ReadFile(path)
	if e != nil {
		return e
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return errors.New("invalid env line")
		}
		if parts[0] == "" {
			return nil
		}
		if os.Getenv(parts[0]) == "" {
			_ = os.Setenv(parts[0], strings.Trim(parts[1], "\""))
		}
	}
	return nil
}
func RequiredSecret(ref SecretRef, env map[string]string) error {
	if v, ok := ref.ResolveOK(env); ok && v != "" {
		return nil
	}
	src := ref.From
	if src == "" {
		src = ref.Name
	}
	return errors.New("required secret is missing or empty: " + src)
}
func Required(c Config) error {
	if c.Address == "" || c.Workers < 1 {
		return ErrInvalid
	}
	return nil
}
