package config

import "time"

type Limits struct {
	MaxRequestBytes int64
	MaxExecution    time.Duration
	MaxModules      int
	MaxTenants      int
}

func DefaultLimits() Limits {
	return Limits{MaxRequestBytes: 4 << 20, MaxExecution: 30 * time.Second, MaxModules: 10000, MaxTenants: 1000}
}
func (l Limits) Valid() bool {
	return l.MaxRequestBytes > 0 && l.MaxExecution > 0 && l.MaxModules > 0 && l.MaxTenants > 0
}
