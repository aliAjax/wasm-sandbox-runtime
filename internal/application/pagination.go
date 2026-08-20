package application

type Page[T any] struct {
	Items      []T  `json:"items"`
	Offset     int  `json:"offset"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	NextOffset *int `json:"next_offset,omitempty"`
}

func MakePage[T any](items []T, offset, limit, total int) Page[T] {
	p := Page[T]{Items: items, Offset: offset, Limit: limit, Total: total}
	n := offset + len(items)
	if n < total {
		p.NextOffset = &n
	}
	return p
}
func ClampPage(offset, limit, max int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit < 1 {
		limit = 20
	}
	if limit > max {
		limit = max
	}
	return offset, limit
}
