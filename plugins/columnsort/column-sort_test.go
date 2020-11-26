package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseInteger(t *testing.T) {
	res, err := ParseInteger("5,6,7,001,010,016,00")
	assert.NoError(t, err)
	assert.Len(t, res, 7)
	assert.True(t, res[3] == 1)
	assert.True(t, res[4] == 10)
	assert.True(t, res[5] == 16)
	assert.True(t, res[6] == 0)

	res, err = ParseInteger("5,,7")
	assert.Error(t, err)
	assert.Len(t, res, 1)

	res, err = ParseInteger("f5,,7")
	assert.Error(t, err)
	assert.Len(t, res, 0)
}

func TestSorting(t *testing.T) {
	res := Sorting([]int{3, 0, 16}, 0)
	assert.Len(t, res, 3)
	assert.True(t, res[0] == "0")
	assert.True(t, res[1] == "3")
	assert.True(t, res[2] == "16")

	res = Sorting([]int{3, 0, 16}, 1)
	assert.Len(t, res, 3)
	assert.True(t, res[0] == "16")
	assert.True(t, res[1] == "3")
	assert.True(t, res[2] == "0")
}
