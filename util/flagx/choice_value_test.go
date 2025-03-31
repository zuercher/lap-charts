package flagx

import (
	"flag"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoiceValue(t *testing.T) {
	assert := assert.New(t)

	{
		value := 0

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(&ChoiceValue[int]{Value: &value, Choices: []int{1, 2, 3}, Parse: strconv.Atoi}, "v", "usage")
		err := flags.Parse([]string{"-v", "1"})
		assert.NoError(err)
		assert.Equal(1, value)
	}

	{
		value := 0

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(&ChoiceValue[int]{Value: &value, Choices: []int{1, 2, 3}, Parse: strconv.Atoi}, "v", "usage")
		err := flags.Parse([]string{"-v", "10"})
		assert.Error(err)
	}

	{
		value := ""

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(NewStringChoiceValue(&value, "a", "b", "c"), "v", "usage")
		err := flags.Parse([]string{"-v", "a"})
		assert.NoError(err)
		assert.Equal("a", value)
	}

	{
		value := 0

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(&ChoiceValue[int]{Value: &value, Choices: []int{1, 2, 3}, Parse: strconv.Atoi}, "v", "usage")
		err := flags.Parse([]string{"-h"})
		assert.Equal(flag.ErrHelp, err)
		flags.Usage()
	}
}
