package tproger

import (
	"os"
	"strconv"
	"unicode"
)

/*
Дано: два неупорядоченных среза.
а) a := []int{37, 5, 1, 2} и b := []int{6, 2, 4, 37}. => 2,37
б) a = []int{1, 1, 1} и b = []int{1, 1, 1, 1}. 1,1,1
Верните их пересечение.
*/
func Task7V1(a, b []int) []int {
	m := map[int]int{}
	for _, element := range a {
		m[element]++
	}

	result := make([]int, 0)
	for _, element := range b {
		if m[element] > 0 {
			result = append(result, element)
			m[element]--
		}
	}
	return result
}

func Task9V1(a string) string {
	nums := map[rune]struct{}{
		'0': {}, '1': {}, '2': {}, '3': {}, '4': {},
		'5': {}, '6': {}, '7': {}, '8': {}, '9': {},
	}
	var result []rune
	for _, v := range []rune(a) {
		if _, ok := nums[v]; ok {
			result = append(result, v)
		}
	}
	return string(result)
}

func Task9V2(a string) string {
	var result []rune
	for _, r := range []rune(a) {
		if unicode.IsDigit(r) {
			result = append(result, r)

		}
	}
	return string(result)
}

func Task11V1(a, b []int) []int {
	result := make([]int, len(a)+len(b))
	for i, v := range a {
		result[i] = v
	}
	for i, v := range b {
		result[i+len(a)] = v
	}
	return result
}

func Task11V2(a, b []int) []int {
	result := make([]int, len(a)+len(b))
	n := copy(result, a)
	copy(result[n:], b)
	return result
}

func Task11V3(a, b []int) []int {
	return append(a, b...)
}

func Task13(in string) *int {
	v, err := strconv.Atoi(in)
	if err != nil {
		return nil
	}
	return &v
}

func Task14(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()
	return true, nil
}

func Task15(numbers <-chan int) <-chan int {
	if numbers == nil {
		return nil
	}
	out := make(chan int, len(numbers))
	go func() {
		defer close(out)
		for {
			val, ok := <-numbers
			if !ok {
				return
			}
			out <- val * val
		}
	}()
	return out
}
