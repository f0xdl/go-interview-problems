package checkurl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type result struct {
	workerId int
	value    string
	err      error
}

func CheckUrls(ctx context.Context, workerCount int, urls []string) ([]string, error) {
	inputCh := generator(ctx, urls)
	resultCh := make(chan result, len(urls))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := &sync.WaitGroup{}
	for id := range workerCount {
		wg.Add(1)
		go worker(ctx, id, wg, inputCh, resultCh)
	}
	wg.Wait()
	close(resultCh)

	correctUrls := make([]string, 0, len(urls))
	errs := make([]error, 0, len(urls))
	for res := range resultCh {
		if res.err == nil {
			correctUrls = append(correctUrls, res.value)
		} else {
			errs = append(errs, res.err)
		}
	}
	return correctUrls, errors.Join(errs...)
}

func generator[T string | int](ctx context.Context, urls []T) <-chan T {
	ch := make(chan T, len(urls))
	go func() {
		defer close(ch)

		for _, url := range urls {
			select {
			case ch <- url:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

func worker(ctx context.Context, id int, wg *sync.WaitGroup, urls <-chan string, out chan<- result) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("worker %d: panic: %v", id, r)
		}
	}()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-urls:
			if !ok {
				return
			}
			res := result{value: url}
			res.err = checkUrl(url)
			select {
			case <-ctx.Done():
				return
			case out <- res:
			}

		}
	}

}

func checkUrl(url string) error {
	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("return wrong status code %v", r.StatusCode)
	}
	return nil
}
