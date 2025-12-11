package tproger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandGeneratorReturnsChannel(t *testing.T) {
	ch := randGenerator(1)
	assert.NotNil(t, ch)
}

func TestRandGeneratorChannelClosedAfterIteration(t *testing.T) {
	ch := randGenerator(3)
	for range ch {
	}

	_, ok := <-ch
	assert.False(t, ok, "channel should be closed")
}

func TestRandGeneratorReturnsIntValues(t *testing.T) {
	ch := randGenerator(5)
	count := 0
	for val := range ch {
		assert.IsType(t, 0, val)
		count++
	}
	assert.Equal(t, 5, count)
}

func TestRandGeneratorNegativeInput(t *testing.T) {
	ch := randGenerator(-1)
	count := 0
	for range ch {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestRandGeneratorOneValue(t *testing.T) {
	ch := randGenerator(1)
	count := 0
	for range ch {
		count++
	}
	assert.Equal(t, 1, count)
}

func TestRandGeneratorConcurrentReads(t *testing.T) {
	ch := randGenerator(100)
	done := make(chan bool)

	go func() {
		count := 0
		for range ch {
			count++
		}
		assert.Equal(t, 100, count)
		done <- true
	}()

	<-done
}
