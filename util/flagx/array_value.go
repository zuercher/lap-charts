package flagx

import (
	"flag"
	"fmt"
	"strings"
)

type ArrayValue[T any] struct {
	Values         *[]T
	Parse          func(string) (T, error)
	PreserveSpaces bool
}

func NewStringArrayValue(values *[]string, preserveSpaces bool) *ArrayValue[string] {
	return &ArrayValue[string]{
		Values:         values,
		Parse:          func(s string) (string, error) { return s, nil },
		PreserveSpaces: preserveSpaces,
	}
}

var _ flag.Getter = &ArrayValue[any]{}

func (v *ArrayValue[T]) String() string {
	if v == nil || v.Values == nil || len(*(v.Values)) == 0 {
		return "[]"
	}

	b := &strings.Builder{}
	b.WriteRune('[')
	for _, vv := range *(v.Values) {
		b.WriteString(fmt.Sprintf("%v", vv))
	}
	b.WriteRune(']')
	return b.String()
}

func (v *ArrayValue[T]) Set(s string) error {
	parts := strings.SplitSeq(s, ",")

	for part := range parts {
		if !v.PreserveSpaces {
			part = strings.TrimSpace(part)
		}
		vv, err := v.Parse(part)
		if err != nil {
			return err
		}
		*(v.Values) = append(*(v.Values), vv)

	}

	return nil
}

func (v *ArrayValue[T]) Get() any {
	return *(v.Values)
}
