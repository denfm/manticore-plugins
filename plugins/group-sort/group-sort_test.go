package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseInteger(t *testing.T) {
	res, err := ParseInteger("5,6,7")
	assert.NoError(t, err)
	assert.Len(t, res, 3)

	res, err = ParseInteger("5,,7")
	assert.Error(t, err)
	assert.Len(t, res, 1)

	res, err = ParseInteger("f5,,7")
	assert.Error(t, err)
	assert.Len(t, res, 0)
}

func TestParseSimilarPosition(t *testing.T) {
	res, err := ParseInteger("005,010,561")
	assert.NoError(t, err)
	assert.Len(t, res, 3)
	assert.True(t, res[0] == 5)
	assert.True(t, res[1] == 10)
	assert.True(t, res[2] == 561)

	res, err = ParseInteger("a005,010,561")
	assert.Error(t, err)
	assert.Len(t, res, 0)

	res, err = ParseInteger("005,010,f561")
	assert.Error(t, err)
	assert.Len(t, res, 2)

	res, err = ParseInteger("005,010,,561")
	assert.Error(t, err)
	assert.Len(t, res, 2)
}

func TestSorting(t *testing.T) {
	productIds := []int{1589, 1590, 1591, 1592, 1593}
	statuses := []int{1, 1, 2, 3, 2}
	similar := []int{20, 30, 2, 0, 1}

	assert.Equal(t, strings.Join(Sorting(productIds, statuses, similar), ","), "1589,1590,1593,1591,1592")
}
