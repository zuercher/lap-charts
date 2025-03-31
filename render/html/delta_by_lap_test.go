package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"zuercher.us/lapcharts/util/timex"
)

func TestComputeTickSize(t *testing.T) {
	assert := assert.New(t)

	tickSize := computeTickSize(14*timex.Minute, 5)
	assert.Equal(2*timex.Minute, tickSize)

	tickSize = computeTickSize(11*timex.Minute, 5)
	assert.Equal(2*timex.Minute, tickSize)

	tickSize = computeTickSize(31*timex.Minute, 5)
	assert.Equal(10*timex.Minute, tickSize)

	tickSize = computeTickSize(110*timex.Second, 5)
	assert.Equal(30*timex.Second, tickSize)

	tickSize = computeTickSize(1*timex.Minute, 5)
	assert.Equal(20*timex.Second, tickSize)

	tickSize = computeTickSize(1*timex.Second, 5)
	assert.Equal(200*timex.Millisecond, tickSize)
}
