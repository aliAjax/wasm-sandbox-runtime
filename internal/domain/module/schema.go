package module

import (
	"encoding/json"
	"errors"
)

func ValidateJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		return errors.New("schema cannot be null")
	}
	return nil
}
func Compatible(oldSchema, newSchema Schema) error {
	if err := ValidateJSON([]byte(oldSchema.Input)); err != nil {
		return err
	}
	if err := ValidateJSON([]byte(newSchema.Input)); err != nil {
		return err
	}
	return nil
}
func NormalizeSchema(s Schema) Schema {
	return Schema{Input: normalize(s.Input), Output: normalize(s.Output)}
}
func normalize(s string) string {
	var v any
	if json.Unmarshal([]byte(s), &v) != nil {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}
