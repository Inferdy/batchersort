package main

import (
	"fmt"

	batcher "github.com/Inferdy/batchersort"
)

const (
	maxSize int = 40
)

var values []int = make([]int, maxSize)

func checkSolution(size int) bool {
	if size <= 1 {
		return true
	}
	var limit int = size - 2
	var i int = 0
	for {
		if values[i] > values[i+1] {
			return false
		}

		if i == limit {
			return true
		}

		i++
	}
}

func build(a int, b int, c int, d int) {
	var i int = 0

	for a > 0 {
		values[i] = 0
		a--
		i++
	}

	for b > 0 {
		values[i] = 1
		b--
		i++
	}

	for c > 0 {
		values[i] = 0
		c--
		i++
	}

	for d > 0 {
		values[i] = 1
		d--
		i++
	}
}

func test(a int, b int, c int, d int, first int, second int, total int) bool {
	build(a, b, c, d)
	batcher.Sort2(values, first, second, total, 4, 2)
	return checkSolution(total)
}

func main() {
	var second int

	for total := 0; total <= maxSize; total++ {
		for first := 0; first <= total; first++ {
			second = total - first

			for a := 0; a <= first; a++ {
				b := first - a
				for c := 0; c <= second; c++ {
					d := second - c

					if !test(a, b, c, d, first, second, total) {
						fmt.Printf("a = %d, b = %d, c = %d, d = %d, total = %d\n", a, b, c, d, total)
					}
				}
			}
		}
	}
}
