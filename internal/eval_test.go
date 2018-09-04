package internal

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestArgumentsToMeasuresAndDims1(t *testing.T) {
	measures, dims := argumentsToMeasuresAndDims([]string{"measure1", "measure2", "by", "dim1", "dim2"})
	assert.Equal(t, measures[0], "measure1")
	assert.Equal(t, measures[1], "measure2")
	assert.Equal(t, dims[0], "dim1")
	assert.Equal(t, dims[1], "dim2")
}

func TestArgumentsToMeasuresAndDims2(t *testing.T) {
	measures, dims := argumentsToMeasuresAndDims([]string{"dim1", "dim2"})
	assert.Equal(t, len(measures), 0)
	assert.Equal(t, dims[0], "dim1")
	assert.Equal(t, dims[1], "dim2")
}

func TestArgumentsToMeasuresAndDimsStar(t *testing.T) {
	measures, dims := argumentsToMeasuresAndDims([]string{"measure1", "measure2", "by", "*"})
	assert.Equal(t, measures[0], "measure1")
	assert.Equal(t, measures[1], "measure2")
	assert.Equal(t, dims, []string{})
}
