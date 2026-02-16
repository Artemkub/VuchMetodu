package main

import (
	"fmt"
	"math"
)

const eps = 1e-15

func StableMetod(a, b, c float64) (float64, float64, bool) {
	scale := math.Max(math.Abs(a), math.Max(math.Abs(b), math.Abs(c)))
	if scale == 0 {
		return 0, 0, false
	}

	aScaled := a / scale
	bScaled := b / scale
	cScaled := c / scale

	if math.Abs(aScaled) < eps {
		if math.Abs(bScaled) < eps {
			return 0, 0, false
		}
		x := -cScaled / bScaled
		return x, x, true
	}

	D := bScaled*bScaled - 4*aScaled*cScaled
	if D < 0 {
		return 0, 0, false
	}

	sqrtD := math.Sqrt(D)

	var q float64
	if bScaled >= 0 {
		q = -0.5 * (bScaled + sqrtD)
	} else {
		q = -0.5 * (bScaled - sqrtD)
	}

	x1 := q / aScaled

	var x2 float64
	if math.Abs(q) < eps {
		x2 = -bScaled / (2 * aScaled)
	} else {
		x2 = cScaled / q
	}

	return x1, x2, true
}

func DiscriminantMetod(a, b, c float64) (float64, float64, bool) {
	if math.Abs(a) < eps {
		if math.Abs(b) < eps {
			return 0, 0, false
		}
		x := -c / b
		return x, x, true
	}

	D := b*b - 4*a*c

	if D < 0 {
		return 0, 0, false
	}

	sqrtD := math.Sqrt(D)

	x2 := (-b + sqrtD) / (2 * a)
	x1 := (-b - sqrtD) / (2 * a)

	return x1, x2, true
}

func main() {
	var a float64 = 1
	var b float64 = 100000000
	var c float64 = 1

	fmt.Println("Stable metod")
	x1, x2, ok := StableMetod(a, b, c)
	if ok {
		fmt.Println("x1 =", x1)
		fmt.Println("x2 =", x2)

		res1 := a*x1*x1 + b*x1 + c
		res2 := a*x2*x2 + b*x2 + c

		if res1 != 1 {
			fmt.Println("res =", res1)
		}
		if res2 != 1 {
			fmt.Println("res =", res2)
		}
	} else {
		fmt.Println("Net kornei")
	}

	fmt.Println("Discriminant metod")
	x1d, x2d, okd := DiscriminantMetod(a, b, c)
	if okd {
		fmt.Println("x1 =", x1d)
		fmt.Println("x2 =", x2d)

		res1 := a*x1d*x1d + b*x1d + c
		res2 := a*x2d*x2d + b*x2d + c

		if res1 != 1 {
			fmt.Println("res =", res1)
		}
		if res2 != 1 {
			fmt.Println("res =", res2)
		}
	} else {
		fmt.Println("Net kornei")
	}
}
