package main

import (
	"fmt"
	"math"
)

func jacobi(A [][]float64, b []float64, x0 []float64, eps float64, maxIt int) ([]float64, []float64, int) {
	n := len(b)
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = x0[i]
	}
	xNew := make([]float64, n)
	residuals := []float64{}

	for it := 1; it <= maxIt; it++ {
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < n; j++ {
				if j != i {
					sum += A[i][j] * x[j]
				}
			}
			xNew[i] = (b[i] - sum) / A[i][i]
		}
		for i := 0; i < n; i++ {
			x[i] = xNew[i]
		}

		maxR := 0.0
		for i := 0; i < n; i++ {
			axi := 0.0
			for j := 0; j < n; j++ {
				axi += A[i][j] * x[j]
			}
			r := math.Abs(axi - b[i])
			if r > maxR {
				maxR = r
			}
		}
		residuals = append(residuals, maxR)
		if maxR < eps {
			return x, residuals, it
		}
	}

	return x, residuals, maxIt
}

func seidel(A [][]float64, b []float64, x0 []float64, eps float64, maxIt int) ([]float64, []float64, int) {
	n := len(b)
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = x0[i]
	}
	residuals := []float64{}

	for it := 1; it <= maxIt; it++ {
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < n; j++ {
				if j != i {
					sum += A[i][j] * x[j]
				}
			}
			x[i] = (b[i] - sum) / A[i][i]
		}

		maxR := 0.0
		for i := 0; i < n; i++ {
			axi := 0.0
			for j := 0; j < n; j++ {
				axi += A[i][j] * x[j]
			}
			r := math.Abs(axi - b[i])
			if r > maxR {
				maxR = r
			}
		}
		residuals = append(residuals, maxR)
		if maxR < eps {
			return x, residuals, it
		}
	}

	return x, residuals, maxIt
}

func main() {
	A := [][]float64{
		{4, 1, 1},
		{2, 7, 1},
		{1, -3, 12},
	}

	b := []float64{7, 18, 17}

	x0_1 := []float64{0, 0, 0}
	x0_2 := []float64{10, 10, 10}

	eps := 1e-8
	maxIt := 100

	fmt.Println("Метод Якоби, начальное приближение 1")
	xj1, resj1, itj1 := jacobi(A, b, x0_1, eps, maxIt)
	fmt.Println("Итерации:", itj1)
	fmt.Println("Решение:")
	for i := 0; i < len(xj1); i++ {
		fmt.Println(i+1, xj1[i])
	}
	fmt.Println("Норма невязки на итерациях:")
	for i := 0; i < len(resj1); i++ {
		fmt.Println(i+1, resj1[i])
	}

	fmt.Println()
	fmt.Println("Метод Якоби, начальное приближение 2")
	xj2, resj2, itj2 := jacobi(A, b, x0_2, eps, maxIt)
	fmt.Println("Итерации:", itj2)
	fmt.Println("Решение:")
	for i := 0; i < len(xj2); i++ {
		fmt.Println(i+1, xj2[i])
	}
	fmt.Println("Норма невязки на итерациях:")
	for i := 0; i < len(resj2); i++ {
		fmt.Println(i+1, resj2[i])
	}

	fmt.Println()
	fmt.Println("Метод Зейделя, начальное приближение 1")
	xs1, ress1, its1 := seidel(A, b, x0_1, eps, maxIt)
	fmt.Println("Итерации:", its1)
	fmt.Println("Решение:")
	for i := 0; i < len(xs1); i++ {
		fmt.Println(i+1, xs1[i])
	}
	fmt.Println("Норма невязки на итерациях:")
	for i := 0; i < len(ress1); i++ {
		fmt.Println(i+1, ress1[i])
	}

	fmt.Println()
	fmt.Println("Метод Зейделя, начальное приближение 2")
	xs2, ress2, its2 := seidel(A, b, x0_2, eps, maxIt)
	fmt.Println("Итерации:", its2)
	fmt.Println("Решение:")
	for i := 0; i < len(xs2); i++ {
		fmt.Println(i+1, xs2[i])
	}
	fmt.Println("Норма невязки на итерациях:")
	for i := 0; i < len(ress2); i++ {
		fmt.Println(i+1, ress2[i])
	}
}
