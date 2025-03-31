package flagx

import (
	"flag"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayValue(t *testing.T) {
	assert := assert.New(t)

	{
		values := []int{}

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(&ArrayValue[int]{Values: &values, Parse: strconv.Atoi}, "v", "usage")
		err := flags.Parse([]string{"-v", "1,2,3"})
		assert.NoError(err)
		assert.Equal([]int{1, 2, 3}, values)
	}

	{
		values := []string{}

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(NewStringArrayValue(&values, true), "v", "usage")
		err := flags.Parse([]string{"-v", " a , b , c "})
		assert.NoError(err)
		assert.Equal([]string{" a ", " b ", " c "}, values)
	}

	{
		values := []int{}

		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(&ArrayValue[int]{Values: &values, Parse: strconv.Atoi}, "v", "usage")
		err := flags.Parse([]string{"-h"})
		assert.Equal(flag.ErrHelp, err)
		flags.Usage()
	}
}
