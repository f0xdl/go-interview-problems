package main

import (
	"reflect"
	"sort"
	"testing"

	"golang.org/x/tour/tree"
)

func TestWalk(t *testing.T) {
	tests := []struct {
		name     string
		tree     *tree.Tree
		expected []int
	}{
		{
			name:     "walk tree 1",
			tree:     tree.New(1),
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "walk tree 2",
			tree:     tree.New(2),
			expected: []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
		},
		{
			name:     "walk empty",
			tree:     nil,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan int)

			var result []int
			done := make(chan struct{})
			go func() {
				for v := range ch {
					result = append(result, v)
				}
				done <- struct{}{}
			}()

			Walk(tt.tree, ch)

			<-done

			sort.Ints(result)
			sort.Ints(tt.expected)

			if len(tt.expected) == 0 && len(result) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Walk() got = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSame(t *testing.T) {
	tests := []struct {
		name string
		t1   *tree.Tree
		t2   *tree.Tree
		want bool
	}{
		{
			name: "same values different structure",
			t1:   tree.New(1),
			t2:   tree.New(1),
			want: true,
		},
		{
			name: "different trees",
			t1:   tree.New(1),
			t2:   tree.New(2),
			want: false,
		},
		{
			name: "nil first tree",
			t1:   nil,
			t2:   tree.New(1),
			want: false,
		},
		{
			name: "nil second tree",
			t1:   tree.New(1),
			t2:   nil,
			want: false,
		},
		{
			name: "both nil trees",
			t1:   nil,
			t2:   nil,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Same(tt.t1, tt.t2)
			if got != tt.want {
				t.Errorf("Same() = %v, want %v", got, tt.want)
			}
		})
	}
}
