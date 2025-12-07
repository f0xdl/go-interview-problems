package main

import (
	"context"
	"errors"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	if len(addresses) == 0 {
		return "", nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resCh := make(chan string)
	errCh := make(chan error, len(addresses))

	for _, address := range addresses {
		go worker(ctx, getter, address, key, resCh, errCh)
	}

	errs := []error{}
	for {
		select {
		case r := <-resCh:
			cancel()
			return r, nil
		case err := <-errCh:
			errs = append(errs, err)
			if len(errs) == len(addresses) {
				return "", errors.Join(errs...)
			}
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func worker(ctx context.Context, getter Getter, address, key string, result chan<- string, errs chan<- error) {
	res, err := getter.Get(ctx, address, key)
	if err != nil {
		errs <- err
	} else {
		select {
		case result <- res:
		default:
		}
	}
}
