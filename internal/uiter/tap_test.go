package uiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTap(t *testing.T) {
	mockSeq := func(yield func(int, string) bool) {
		data := []struct {
			key   int
			value string
		}{
			{1, "one"},
			{2, "two"},
			{3, "three"},
		}
		for _, item := range data {
			if !yield(item.key, item.value) {
				return
			}
		}
	}

	// Mock function to track calls
	var called []struct {
		key   int
		value string
	}
	mockFn := func(key int, value string) {
		called = append(called, struct {
			key   int
			value string
		}{key, value})
	}

	// Apply Tap
	tappedSeq := Tap(mockSeq, mockFn)

	// Consume the tapped sequence
	var result []struct {
		key   int
		value string
	}
	for key, value := range tappedSeq {
		result = append(result, struct {
			key   int
			value string
		}{key, value})
	}

	// Assertions
	assert.Equal(t, called, result, "The tapped sequence should match the original sequence")
}
