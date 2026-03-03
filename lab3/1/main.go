package main

import (
	"fmt"
	"math"
)

func factorial(n float64) float64 {
	fact := 1.0
	for i := 2.0; i <= n; i++ {
		fact *= i
	}
	return fact
}

func erfTaylor(x float64) float64 {
	n := 200.0
	sum := 0.0

	for i := 0.0; i <= n; i++ {
		sign := 1.0
		if int(i)%2 == 1 {
			sign = -1.0
		}

		term := sign * math.Pow(x, 2*i+1) / (factorial(i) * (2*i + 1))
		sum += term
	}

	return (2.0 / math.Sqrt(math.Pi)) * sum
}

func matrixNorm(A [][]float64) float64 {
	rows := len(A)
	cols := len(A[0])

	maxColSum := 0.0

	for j := 0; j < cols; j++ {
		colSum := 0.0
		for i := 0; i < rows; i++ {
			colSum += math.Abs(A[i][j])
		}
		if colSum > maxColSum {
			maxColSum = colSum
		}
	}

	return maxColSum
}

func inverseMatrix(A [][]float64) ([][]float64, error) {
	n := 3
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, 2*n)
		for j := 0; j < n; j++ {
			aug[i][j] = A[i][j]
		}
		aug[i][n+i] = 1.0
	}

	for i := 0; i < n; i++ {
		maxRow := i
		for k := i + 1; k < n; k++ {
			if math.Abs(aug[k][i]) > math.Abs(aug[maxRow][i]) {
				maxRow = k
			}
		}

		aug[i], aug[maxRow] = aug[maxRow], aug[i]

		if math.Abs(aug[i][i]) < 1e-15 {
			return nil, fmt.Errorf("")
		}

		pivot := aug[i][i]
		for j := 0; j < 2*n; j++ {
			aug[i][j] /= pivot
		}

		for k := i + 1; k < n; k++ {
			factor := aug[k][i]
			for j := 0; j < 2*n; j++ {
				aug[k][j] -= factor * aug[i][j]
			}
		}
	}

	for i := n - 1; i >= 0; i-- {
		for k := i - 1; k >= 0; k-- {
			factor := aug[k][i]
			for j := 0; j < 2*n; j++ {
				aug[k][j] -= factor * aug[i][j]
			}
		}
	}

	inv := make([][]float64, n)
	for i := 0; i < n; i++ {
		inv[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			inv[i][j] = aug[i][n+j]
		}
	}

	return inv, nil
}

func solveLinearSystem(A [][]float64, b []float64) ([]float64, error) {
	n := len(A)
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, n+1)
		for j := 0; j < n; j++ {
			aug[i][j] = A[i][j]
		}
		aug[i][n] = b[i]
	}

	for i := 0; i < n; i++ {
		maxRow := i
		for k := i + 1; k < n; k++ {
			if math.Abs(aug[k][i]) > math.Abs(aug[maxRow][i]) {
				maxRow = k
			}
		}

		aug[i], aug[maxRow] = aug[maxRow], aug[i]

		if math.Abs(aug[i][i]) < 1e-15 {
			return nil, fmt.Errorf("")
		}

		pivot := aug[i][i]
		for j := i; j < n+1; j++ {
			aug[i][j] /= pivot
		}

		for k := i + 1; k < n; k++ {
			factor := aug[k][i]
			for j := i; j < n+1; j++ {
				aug[k][j] -= factor * aug[i][j]
			}
		}
	}

	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = aug[i][n]
		for j := i + 1; j < n; j++ {
			x[i] -= aug[i][j] * x[j]
		}
	}

	return x, nil
}

func computeResidual(A [][]float64, x []float64, b []float64) float64 {
	n := len(A)
	residual := 0.0

	for i := 0; i < n; i++ {
		axi := 0.0
		for j := 0; j < n; j++ {
			axi += A[i][j] * x[j]
		}
		diff := math.Abs(axi - b[i])
		if diff > residual {
			residual = diff
		}
	}

	return residual
}

func conditionNumber(A [][]float64) (float64, error) {
	normA := matrixNorm(A)

	invA, err := inverseMatrix(A)
	if err != nil {
		return 0, err
	}

	normInvA := matrixNorm(invA)

	return normA * normInvA, nil
}

func task1() {
	A := [][]float64{
		{1.00, 0.80, 0.64},
		{1.00, 0.90, 0.81},
		{1.00, 1.10, 1.21},
	}

	b := []float64{
		erfTaylor(0.80),
		erfTaylor(0.90),
		erfTaylor(1.10),
	}

	cond, err := conditionNumber(A)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Обусловленность")
		fmt.Println(cond)
	}

	x, err := solveLinearSystem(A, b)
	if err != nil {
		fmt.Printf("Ошибка при решении: %v\n", err)
	} else {
		fmt.Printf("Решение")
		fmt.Println(x[0])
		fmt.Println(x[1])
		fmt.Println(x[2])
	}

	fmt.Println("Невязка")
	neviaz := computeResidual(A, x, b)
	fmt.Println(neviaz)

	fmt.Println("Сравнение")
	sumX := x[0] + x[1] + x[2]
	erf1 := erfTaylor(1.0)

	fmt.Println(sumX)
	fmt.Println(erf1)
	fmt.Println(math.Abs(sumX - erf1))
}

func task2() {
	A := [][]float64{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
		{0.7, 0.8, 0.9},
	}

	b := []float64{0.1, 0.3, 0.5}

	det := A[0][0]*A[1][1]*A[2][2] + A[0][1]*A[1][2]*A[2][0] + A[0][2]*A[1][0]*A[2][1] -
		A[0][2]*A[1][1]*A[2][0] - A[0][0]*A[1][2]*A[2][1] - A[0][1]*A[1][0]*A[2][2]

	c2 := (A[2][0]*A[0][1] - A[0][0]*A[2][1]) / (A[1][0]*A[0][1] - A[0][0]*A[1][1])
	c1 := (A[2][0] - c2*A[1][0]) / A[0][0]

	b_sum := c1*b[0] + c2*b[1]

	if math.Abs(det) < 1e-10 {
		if math.Abs(b[2]-b_sum) < 1e-10 {
			fmt.Println("Множество решений")
		} else {
			fmt.Println("Одно решение")
		}
	}
}

func main() {
	task1()
	task2()
}
