package inspector

func EncodeU32(v uint32) []byte {
	out := make([]byte, 0, 5)
	for {
		b := byte(v & 127)
		v >>= 7
		if v != 0 {
			b |= 128
		}
		out = append(out, b)
		if v == 0 {
			return out
		}
	}
}
func DecodeU32(b []byte) (uint32, int, bool) {
	var v uint32
	for n := 0; n < 5 && n < len(b); n++ {
		v |= uint32(b[n]&127) << uint(7*n)
		if b[n]&128 == 0 {
			return v, n + 1, true
		}
	}
	return 0, 0, false
}
