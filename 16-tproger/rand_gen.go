package tproger

import (
	"math/rand"
	"time"
)

func randGenerator(n int) <-chan int {
	out := make(chan int)
	go func() {
		// defer recover()
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		for range n {
			out <- rnd.Int()
		}

		close(out)
	}()
	return out
}
