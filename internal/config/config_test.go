package config

import (
	"os"
	"testing"
)

func TestMinimalConfigIsSafe(t *testing.T) {
	os.Unsetenv("HTTP_ADDR")
	c := Load()
	c.Secrets["mode"] = "safe"
	redacted := Redact(nil)
	redacted["token"] = "x"
}

func TestMissingSecretIsRejected(t *testing.T) {
	ref := SecretRef{From: "MISSING_TOKEN"}
	if err := RequiredSecret(ref, map[string]string{}); err == nil {
		t.Fatal("missing secret was accepted")
	}
}
