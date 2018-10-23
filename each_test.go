package workerpool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapEach(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := New(ctx, 5)

	o := wp.Each(m, func(k string, v string) string {
		assert.Equal(t, k, v)
		return v
	}).(map[string]string)

	for k, v := range m {
		assert.Equal(t, v, o[k])
	}
}

func TestMapEachVoid(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := New(ctx, 5)

	o := wp.Each(m, func(k string, v string) {
		assert.Equal(t, k, v)
	})
	assert.Nil(t, o)
}

func TestArrayEach(t *testing.T) {
	m := []int{0, 1, 2}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := New(ctx, 5)

	o := wp.Each(m, func(k int, v int) int {
		assert.Equal(t, k, v)
		return v
	}).([]int)

	for k, v := range m {
		assert.Equal(t, v, o[k])
	}
}

func TestArrayEachVoid(t *testing.T) {
	m := []int{0, 1, 2}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := New(ctx, 5)

	o := wp.Each(m, func(k int, v int) {
		assert.Equal(t, k, v)
	})
	assert.Nil(t, o)
}
