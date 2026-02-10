package main

import (
	"fmt"
	"math"
)

func summa1(e float64) (float64, int) {
	sum := 0.0
	n := 1

	for {
		temp := 1.0 / (float64(n)*float64(n) + 1.0)
		sum += temp

		nexttemp := 1.0 / (float64(n+1)*float64(n+1) + 1.0)

		if nexttemp < e {
			break
		}

		n++
	}

	return sum, n
}

func summa2(e float64) (float64, int) {

	Part1 := (math.Pi * math.Pi / 6.0) - (math.Pow(math.Pi, 4) / 90.0)

	Part2 := 0.0
	n := 1

	for {
		n2 := float64(n) * float64(n)
		n4 := n2 * n2

		temp := 1.0 / (n4 * (n2 + 1.0))
		Part2 += temp

		nexttemp := 1.0 / (math.Pow(float64(n+1), 4) * (math.Pow(float64(n+1), 2) + 1.0))

		if nexttemp < e {
			break
		}

		n++

		if n > 10000000 {
			break
		}
	}

	total := Part1 + Part2

	return total, n
}

func main() {
	e := 1e-10
	sum1, n1 := summa1(e)
	fmt.Println(sum1, n1)
	sum2, n2 := summa2(e)
	fmt.Println(sum2, n2)
}
