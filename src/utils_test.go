
package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInterpolate(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(3, interpolate(0, map[int]int{1: 3}))
	assert.Equal(3, interpolate(1, map[int]int{1: 3}))
	assert.Equal(3, interpolate(2, map[int]int{1: 3}))

	assert.Equal(3, interpolate(0, map[int]int{1: 3, 3: 5}))
	assert.Equal(3, interpolate(1, map[int]int{1: 3, 3: 5}))
	assert.Equal(4, interpolate(2, map[int]int{1: 3, 3: 5}))
	assert.Equal(5, interpolate(3, map[int]int{1: 3, 3: 5}))
	assert.Equal(5, interpolate(4, map[int]int{1: 3, 3: 5}))

	assert.Equal(3, interpolate(0, map[int]int{1: 3, 3: 5, 5: 0}))
	assert.Equal(3, interpolate(1, map[int]int{1: 3, 3: 5, 5: 0}))
	assert.Equal(4, interpolate(2, map[int]int{1: 3, 3: 5, 5: 0}))
	assert.Equal(5, interpolate(3, map[int]int{1: 3, 3: 5, 5: 0}))
	assert.Equal(2, interpolate(4, map[int]int{1: 3, 3: 5, 5: 0}))
	assert.Equal(0, interpolate(5, map[int]int{1: 3, 3: 5, 5: 0}))
	assert.Equal(0, interpolate(6, map[int]int{1: 3, 3: 5, 5: 0}))
}

func TestMinMax(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, minMax(1, 0, 3))
	assert.Equal(1, minMax(1, 1, 3))
	assert.Equal(2, minMax(1, 2, 3))
	assert.Equal(3, minMax(1, 3, 3))
	assert.Equal(3, minMax(1, 4, 3))

	assert.Equal(1, minMax(1, 0, 1))
	assert.Equal(1, minMax(1, 1, 1))
	assert.Equal(1, minMax(1, 2, 1))
}

func TestGetSortedKeys(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]int {}, getSortedKeys(&map[int]int {}))
	assert.Equal([]int {1}, getSortedKeys(&map[int]int {1: 2}))
	assert.Equal([]int {1, 3, 5}, getSortedKeys(&map[int]int {1: 2, 3: 4, 5: 6}))
	assert.Equal([]int {1, 2, 6}, getSortedKeys(&map[int]int {2: 4, 1: 3, 6: -1}))
}
