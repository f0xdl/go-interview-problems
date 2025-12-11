package aggregator

import (
	"context"
	"sync"
	"sync/atomic"
)

func Aggregate(ctx context.Context, sources []<-chan int) <-chan int {
	result := make(chan int, 1)
	var wg sync.WaitGroup
	var sum int64

	for _, source := range sources {
		if source != nil {
			wg.Add(1)
			go func(source <-chan int) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case value, ok := <-source:
						if !ok {
							return
						}
						atomic.AddInt64(&sum, int64(value))
					}
				}
			}(source)
		}
	}
	go func() {
		wg.Wait()
		result <- int(sum)
		close(result)
	}()

	return result
}
