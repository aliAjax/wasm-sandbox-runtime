package runtime

import "errors"

type HostLimits struct {
	MaxFunctions   int
	MaxImports     int
	MaxMemoryPages uint32
	MaxOutput      uint64
}

func (l HostLimits) Validate() error {
	if l.MaxFunctions < 1 || l.MaxImports < 0 || l.MaxMemoryPages < 1 || l.MaxOutput < 1 {
		return errors.New("host limits invalid")
	}
	return nil
}
