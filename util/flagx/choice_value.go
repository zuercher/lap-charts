package flagx

import (
	"flag"
	"fmt"
	"strings"
)

type ChoiceValue[T comparable] struct {
	Value   *T
	Choices []T
	Parse   func(string) (T, error)
}

func NewStringChoiceValue(value *string, choices ...string) *ChoiceValue[string] {
	return &ChoiceValue[string]{
		Value:   value,
		Choices: choices,
		Parse:   func(s string) (string, error) { return s, nil },
	}
}

var _ flag.Getter = &ChoiceValue[any]{}

func (v *ChoiceValue[T]) String() string {
	if v == nil || v.Value == nil {
		return ""
	}

	return fmt.Sprintf("%v", *v.Value)
}

func (v *ChoiceValue[T]) Set(s string) error {
	vv, err := v.Parse(s)
	if err != nil {
		return err
	}

	for _, choice := range v.Choices {
		if vv == choice {
			*v.Value = vv
			return nil
		}
	}

	return fmt.Errorf("%s is not a valid choice", s)
}

func (v *ChoiceValue[T]) Get() any {
	return *v.Value
}

func (v *ChoiceValue[T]) Describe() string {
	choices := []string{}
	for _, choice := range v.Choices {
		choices = append(choices, fmt.Sprintf("%v", choice))
	}
	switch len(choices) {
	case 0:
		return "n/a"
	case 1:
		return choices[0]
	case 2:
		return fmt.Sprintf("%s or %s", choices[0], choices[1])
	default:
		return fmt.Sprintf("%s, or %s", strings.Join(choices[:len(choices)-1], ", "), choices[len(choices)-1])
	}
}
