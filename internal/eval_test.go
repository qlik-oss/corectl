package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentsToMeasuresAndDims1(t *testing.T) {
	measures, dims := argumentsToMeasuresAndDims([]string{"measure1", "measure2", "by", "dim1", "dim2"})
	assert.Equal(t, "measure1", measures[0])
	assert.Equal(t, "measure2", measures[1])
	assert.Equal(t, "dim1", dims[0])
	assert.Equal(t, "dim2", dims[1])
}

func TestArgumentsToMeasuresAndDims2(t *testing.T) {
	measures, dims := argumentsToMeasuresAndDims([]string{"by", "dim1", "dim2"})
	assert.Equal(t, 0, len(measures))
	assert.Equal(t, "dim1", dims[0])
	assert.Equal(t, "dim2", dims[1])
}

func TestArgumentsToMeasuresAndDimsStar(t *testing.T) {
	measures, dims := argumentsToMeasuresAndDims([]string{"measure1", "measure2"})
	assert.Equal(t, "measure1", measures[0])
	assert.Equal(t, "measure2", measures[1])
	assert.Equal(t, []string{}, dims)
}
