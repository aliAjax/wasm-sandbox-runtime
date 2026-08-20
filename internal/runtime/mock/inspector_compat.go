package mock

import (
	"errors"
	"github.com/example/wasm-sandbox-runtime/internal/wasm/inspector"
)

func Entrypoint(info inspector.ModuleInfo, name string) error {
	for _, e := range info.Exports {
		if e.Name == name && e.Kind == 0 {
			return nil
		}
	}
	return errors.New("entrypoint not exported")
}
