package module

import "sort"

func SortVersions(v []Version) []Version {
	out := append([]Version(nil), v...)
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out
}
func NextVersion(v []Version) int {
	n := 0
	for _, x := range v {
		if x.Number > n {
			n = x.Number
		}
	}
	return n + 1
}
func VersionNumbers(v []Version) []int {
	out := make([]int, 0, len(v))
	for _, x := range SortVersions(v) {
		out = append(out, x.Number)
	}
	return out
}
