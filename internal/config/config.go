package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address        string
	Workers        int
	QueueSize      int
	DefaultTimeout time.Duration
	MaxBody        int64
	Lease          time.Duration
	Secrets        map[string]string
}

func Load() Config {
	return Config{Address: env("HTTP_ADDR", ":8087"), Workers: envInt("WORKERS", 4), QueueSize: envInt("QUEUE_SIZE", 256), DefaultTimeout: envDuration("DEFAULT_TIMEOUT", 30*time.Second), MaxBody: int64(envInt("MAX_BODY_BYTES", 4<<20)), Lease: envDuration("LEASE_DURATION", 20*time.Second), Secrets: nil}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func envInt(k string, d int) int {
	v, e := strconv.Atoi(env(k, ""))
	if e != nil {
		return d
	}
	return v
}
func envDuration(k string, d time.Duration) time.Duration {
	v, e := time.ParseDuration(env(k, ""))
	if e != nil {
		return d
	}
	return v
}
func (c Config) Validate() error {
	if c.Workers < 1 || c.QueueSize < c.Workers || c.MaxBody < 1024 {
		return ErrInvalid
	}
	return nil
}

var ErrInvalid = &configError{}

type configError struct{}

func (*configError) Error() string { return "invalid configuration" }
