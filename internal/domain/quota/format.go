package quota

import "fmt"

func FormatCPU(v float64) string { return fmt.Sprintf("%.3f CPU-seconds", v) }
func FormatBytes(v uint64) string {
	const m = 1024 * 1024
	if v >= m {
		return fmt.Sprintf("%.1f MiB", float64(v)/m)
	}
	return fmt.Sprintf("%d bytes", v)
}
func Summary(q Quota) string {
	c, m, s := q.Utilization()
	return fmt.Sprintf("tenant=%s concurrency=%d/%d cpu=%.2f memory=%.2f storage=%.2f", q.TenantID, q.UsedConcurrent, q.MaxConcurrent, c, m, s)
}
