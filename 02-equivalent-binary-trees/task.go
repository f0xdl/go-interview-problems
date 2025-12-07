package main

import (
	"sync"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go walk(t, ch, wg)
	wg.Wait()
	if ch != nil {
		close(ch) //incorrect pattern, only for done testcase
	}
}

func walk(t *tree.Tree, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	if t == nil {
		return
	}
	ch <- t.Value
	wg.Add(2)
	go walk(t.Left, ch, wg)
	go walk(t.Right, ch, wg)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool { //solo unique keys
	if t1 == nil || t2 == nil {
		return t1 == t2
	}
	t1Ch := make(chan int)
	t2Ch := make(chan int)
	go Walk(t1, t1Ch)
	go Walk(t2, t2Ch)

	values := map[int]int{}
	for {
		select {
		case v, ok := <-t1Ch:
			if !ok {
				t1Ch = nil
			}
			values[v]++
		case v, ok := <-t2Ch:
			if !ok {
				t2Ch = nil
			}
			values[v]++
		}
		if t1Ch == nil && t2Ch == nil {
			break
		}
	}

	for _, v := range values {
		if v != 2 {
			return false
		}
	}
	return true
}
