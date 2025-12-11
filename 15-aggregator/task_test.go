package aggregator_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	. "github.com/blindlobstar/go-interview-problems/15-aggregator"
	"github.com/stretchr/testify/require"
)

func makeChan(vals ...int) chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func TestAggregate_EmptySources(t *testing.T) {
	ctx := context.Background()
	res := Aggregate(ctx, nil)

	v, ok := <-res
	require.True(t, ok)
	require.Equal(t, 0, v)
}

func TestAggregate_SingleSource(t *testing.T) {
	ctx := context.Background()
	src := makeChan(1, 2, 3)

	res := Aggregate(ctx, []<-chan int{src})

	v := <-res
	require.Equal(t, 6, v)
}

func TestAggregate_MultipleSources(t *testing.T) {
	ctx := context.Background()

	src1 := makeChan(1, 2, 3)
	src2 := makeChan(4, 5)
	src3 := makeChan(10)

	res := Aggregate(ctx, []<-chan int{src1, src2, src3})

	v := <-res
	require.Equal(t, 25, v)
}

func TestAggregate_NilChannelsIgnored(t *testing.T) {
	ctx := context.Background()

	src := makeChan(10, 20)

	res := Aggregate(ctx, []<-chan int{nil, src, nil})

	v := <-res
	require.Equal(t, 30, v)
}

func TestAggregate_ContextCancelStopsReading(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := make(chan int)

	go func() {
		src <- 100
		time.Sleep(time.Hour) // гарантированная блокировка, если cancel не работает
	}()

	cancel() // должны завершиться без deadlock

	res := Aggregate(ctx, []<-chan int{src})

	require.Eventually(t, func() bool {
		_, ok := <-res
		return !ok // должен закрыться
	}, time.Second, 10*time.Millisecond)
}

func TestAggregate_NoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	ctx := context.Background()
	src1 := makeChan(1)
	src2 := makeChan(2)
	src3 := makeChan(3)

	<-Aggregate(ctx, []<-chan int{src1, src2, src3})

	time.Sleep(50 * time.Millisecond) // дать горутинам завершиться

	after := runtime.NumGoroutine()

	require.InDelta(t, before, after, 2) // допускаем ±2
}

func TestAggregate_ConcurrencyAtomicity(t *testing.T) {
	ctx := context.Background()

	const N = 1000
	const Goroutines = 20

	srcs := make([]<-chan int, Goroutines)
	for i := range Goroutines {
		ch := make(chan int, N)
		for j := 0; j < N; j++ {
			ch <- 1
		}
		close(ch)
		srcs[i] = ch
	}

	res := Aggregate(ctx, srcs)

	v := <-res
	require.Equal(t, N*Goroutines, v)
}

func TestAggregate_OutputChanClosedExactlyOnce(t *testing.T) {
	ctx := context.Background()
	src := makeChan(1)

	res := Aggregate(ctx, []<-chan int{src})

	v, ok := <-res
	require.True(t, ok)
	require.Equal(t, 1, v)

	_, ok = <-res
	require.False(t, ok) // канал должен быть корректно закрыт
}
