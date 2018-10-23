package workerpool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := New(ctx, 5)

	a := false
	b := false
	c := false

	err := wp.All(ctx, func(ctx context.Context) error {
		a = true
		return nil
	}, func(ctx context.Context) error {
		b = true
		return nil
	}, func(ctx context.Context) error {
		c = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, a)
	assert.True(t, b)
	assert.True(t, c)
}

func TestAllError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp := New(ctx, 5)

	a := false
	b := false
	c := false
	expectedError := fmt.Errorf("expected error")

	err := wp.All(ctx, func(ctx context.Context) error {
		<-ctx.Done()
		a = true
		return nil
	}, func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		b = true
		return expectedError
	}, func(ctx context.Context) error {
		<-ctx.Done()
		c = true
		return nil
	})

	assert.Error(t, expectedError, err)
	assert.True(t, a)
	assert.True(t, b)
	assert.True(t, c)
}
