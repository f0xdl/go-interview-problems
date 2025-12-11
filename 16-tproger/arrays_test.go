package tproger_test

import (
	"testing"

	. "github.com/blindlobstar/go-interview-problems/16-tproger"
	"github.com/stretchr/testify/assert"
)

func Test_task1V1(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{
			name: "empty",
			a:    []int{},
			b:    []int{},
			want: []int{},
		},
		{
			name: "second empty",
			a:    []int{1, 2, 4},
			b:    []int{},
			want: []int{},
		},
		{
			name: "first empty",
			a:    []int{},
			b:    []int{1, 2, 3, 4},
			want: []int{},
		},
		{
			name: "identity",
			a:    []int{1},
			b:    []int{1},
			want: []int{1},
		},
		{
			name: "case 1",
			a:    []int{37, 5, 1, 2},
			b:    []int{6, 2, 4, 37},
			want: []int{2, 37},
		},
		{
			name: "multiple",
			a:    []int{1, 1, 1, 1, 2, 5},
			b:    []int{10, 45, 1, 1, 1},
			want: []int{1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Task7V1(tt.a, tt.b)
			assert.ElementsMatch(t, tt.want, got, tt.name)
		})
	}
}

func TestTask9V1(t *testing.T) {
	tests := []struct {
		name    string
		in, out string
	}{
		{
			name: "check empty",
			in:   "",
			out:  "",
		},
		{
			name: "without numbers",
			in:   "aa",
			out:  "",
		},
		{
			name: "default case",
			in:   "nuba9h8q8g7vwu8",
			out:  "98878",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Task9V1(test.in)
			assert.Equal(t, test.out, got, test.name)
		})
	}
}

func TestTask9V2(t *testing.T) {
	tests := []struct {
		name    string
		in, out string
	}{
		{
			name: "check empty",
			in:   "",
			out:  "",
		},
		{
			name: "without numbers",
			in:   "aa",
			out:  "",
		},
		{
			name: "default case",
			in:   "nuba9h8q8g7vwu8",
			out:  "98878",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Task9V2(test.in)
			assert.Equal(t, test.out, got, test.name)
		})
	}
}

func TestTask11V1(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "basic merge",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "empty first",
			a:        []int{},
			b:        []int{1, 2},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Task11V1(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTask11V2(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "basic merge",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "advanced merge",
			a:        []int{1, 2, 3, 4},
			b:        []int{5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "empty first",
			a:        []int{},
			b:        []int{1, 2},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Task11V2(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTask15(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty channel",
			input:    []int{},
			expected: nil,
		},
		{
			name:     "nil channel",
			input:    nil,
			expected: nil,
		},
		{
			name:     "single element",
			input:    []int{2},
			expected: []int{4},
		},
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 4, 9, 16},
		},
		{
			name:     "zeros",
			input:    []int{0, 0},
			expected: []int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan int, len(tt.input))
			for _, v := range tt.input {
				in <- v
			}
			close(in)

			out := Task15(in)
			var got []int
			for v := range out {
				got = append(got, v)
			}

			assert.Equal(t, tt.expected, got, tt.name)
		})
	}
}
