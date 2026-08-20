package policy

import (
	"encoding/json"
	"errors"
)

func Marshal(p Policy) ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(p)
}
func Unmarshal(data []byte) (Policy, error) {
	var p Policy
	if len(data) == 0 {
		return p, errors.New("empty policy")
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}
	return p, p.Validate()
}
func Clone(p Policy) Policy {
	b, _ := json.Marshal(p)
	var out Policy
	_ = json.Unmarshal(b, &out)
	return out
}
