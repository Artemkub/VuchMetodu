package main

import (
	"fmt"
	"math"
)

func sDirect(x float64, e float64) (float64, int) {
	sum1 := 0.0
	sum2 := 0.0
	k := 1
	count := 0

	for {
		k3 := float64(k) * float64(k) * float64(k)

		temp1 := 1.0 / math.Sqrt(k3+x)
		temp2 := 1.0 / math.Sqrt(k3-x)

		sum1 += temp1
		sum2 += temp2
		count += 2

		if temp1 < e/2 && temp2 < e/2 {
			break
		}

		k++
	}
	return sum1 - sum2, count
}

func sOptimized(x float64, e float64) (float64, int) {
	sum := 0.0
	k := 1
	count := 0

	for {
		k3 := float64(k) * float64(k) * float64(k)

		sqrt1 := math.Sqrt(k3 + x)
		sqrt2 := math.Sqrt(k3 - x)
		obshiu := sqrt1 * sqrt2

		temp := (sqrt2 - sqrt1) / obshiu
		sum += temp
		count += 1

		if math.Abs(temp) < e {
			break
		}

		k++
	}

	return sum, count
}

func main() {
	e := 3e-8

	x1 := 0.5
	result1, count1 := sDirect(x1, e)
	fmt.Println(result1, count1)

	x2 := 0.9999999999
	result2, count2 := sDirect(x2, e)
	fmt.Println(result2, count2)

	time1 := count1 * 500
	fmt.Println(count1, time1)
	time2 := count2 * 500
	fmt.Println(count2, time2)

	res1, optcount1 := sOptimized(x1, e)
	fmt.Println(res1, optcount1)
	res2, optcount2 := sOptimized(x2, e)
	fmt.Println(res2, optcount2)
}
