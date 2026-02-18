package main

import (
	"fmt"
	"math"
)

// Ваша функция вычисления факториала
func factorial(n float64) float64 {
	fact := 1.0
	for i := 2.0; i <= n; i++ {
		fact *= i
	}
	return fact
}

// Ваша функция вычисления erf через ряд Тейлора
func erfTaylor(x float64) float64 {
	// Используем достаточно большое число членов для точности
	n := 200.0
	sum := 0.0

	for i := 0.0; i <= n; i++ {
		// Вычисляем знак: (-1)^i
		sign := 1.0
		if int(i)%2 == 1 {
			sign = -1.0
		}

		// Член ряда: (-1)^i * x^(2i+1) / (i! * (2i+1))
		term := sign * math.Pow(x, 2*i+1) / (factorial(i) * (2*i + 1))
		sum += term
	}

	// Умножаем на коэффициент 2/√π
	return (2.0 / math.Sqrt(math.Pi)) * sum
}

// Норма матрицы: ||A|| = max_j sum_i |a_ij|
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

// Обращение матрицы 3x3
func inverseMatrix3x3(A [][]float64) ([][]float64, error) {
	n := 3
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, 2*n)
		for j := 0; j < n; j++ {
			aug[i][j] = A[i][j]
		}
		aug[i][n+i] = 1.0
	}

	// Прямой ход
	for i := 0; i < n; i++ {
		// Поиск главного элемента
		maxRow := i
		for k := i + 1; k < n; k++ {
			if math.Abs(aug[k][i]) > math.Abs(aug[maxRow][i]) {
				maxRow = k
			}
		}

		aug[i], aug[maxRow] = aug[maxRow], aug[i]

		if math.Abs(aug[i][i]) < 1e-15 {
			return nil, fmt.Errorf("матрица вырождена")
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

	// Обратный ход
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

// Решение СЛАУ методом Гаусса
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

	// Прямой ход
	for i := 0; i < n; i++ {
		maxRow := i
		for k := i + 1; k < n; k++ {
			if math.Abs(aug[k][i]) > math.Abs(aug[maxRow][i]) {
				maxRow = k
			}
		}

		aug[i], aug[maxRow] = aug[maxRow], aug[i]

		if math.Abs(aug[i][i]) < 1e-15 {
			return nil, fmt.Errorf("матрица вырождена")
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

	// Обратный ход
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = aug[i][n]
		for j := i + 1; j < n; j++ {
			x[i] -= aug[i][j] * x[j]
		}
	}

	return x, nil
}

// Вычисление невязки |Ax - b|
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

// Вычисление числа обусловленности cond(A) = ||A|| * ||A^(-1)||
func conditionNumber(A [][]float64) (float64, error) {
	normA := matrixNorm(A)

	invA, err := inverseMatrix3x3(A)
	if err != nil {
		return 0, err
	}

	normInvA := matrixNorm(invA)

	return normA * normInvA, nil
}

// Задача 1
func task1() {
	fmt.Println("ЗАДАЧА 1")
	fmt.Println("=======")

	// Матрица A
	A := [][]float64{
		{1.00, 0.80, 0.64},
		{1.00, 0.90, 0.81},
		{1.00, 1.10, 1.21},
	}

	// Правая часть b = erf(x)
	b := []float64{
		erfTaylor(0.80),
		erfTaylor(0.90),
		erfTaylor(1.10),
	}

	fmt.Println("\nМатрица A:")
	for i := 0; i < 3; i++ {
		fmt.Printf("  [%6.2f %6.2f %6.2f]\n", A[i][0], A[i][1], A[i][2])
	}

	fmt.Println("\nПравая часть b:")
	fmt.Printf("  b1 = erf(0.80) = %.10f\n", b[0])
	fmt.Printf("  b2 = erf(0.90) = %.10f\n", b[1])
	fmt.Printf("  b3 = erf(1.10) = %.10f\n", b[2])

	// а) Обусловленность матрицы и решение
	fmt.Println("\nа) Обусловленность матрицы и решение:")

	cond, err := conditionNumber(A)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("   cond(A) = ||A|| * ||A^(-1)|| = %.6f\n", cond)
	}

	x, err := solveLinearSystem(A, b)
	if err != nil {
		fmt.Printf("Ошибка при решении: %v\n", err)
	} else {
		fmt.Printf("\n   Решение:\n")
		fmt.Printf("   x1 = %.10f\n", x[0])
		fmt.Printf("   x2 = %.10f\n", x[1])
		fmt.Printf("   x3 = %.10f\n", x[2])
	}

	// б) Невязка
	fmt.Println("\nб) Невязка |Ax - b|:")
	residual := computeResidual(A, x, b)
	fmt.Printf("   |Ax - b| = %.2e\n", residual)

	// в) Сумма x1+x2+x3 и сравнение с erf(1.0)
	fmt.Println("\nв) Сумма x1+x2+x3 и сравнение с erf(1.0):")
	sumX := x[0] + x[1] + x[2]
	erf1 := erfTaylor(1.0)

	fmt.Printf("   x1 + x2 + x3 = %.10f\n", sumX)
	fmt.Printf("   erf(1.0) = %.10f\n", erf1)
	fmt.Printf("   Разница = %.2e\n", math.Abs(sumX-erf1))

	fmt.Println("\n   Почему эти числа близки?")
	fmt.Println("   Потому что x1 + x2 + x3 - это значение интерполяционного")
	fmt.Println("   многочлена в точке x=1, построенного по точкам")
	fmt.Println("   (0.8, erf(0.8)), (0.9, erf(0.9)), (1.1, erf(1.1)).")
}

// Задача 2
func task2() {
	fmt.Println("\n\nЗАДАЧА 2")
	fmt.Println("=======")

	A := [][]float64{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
		{0.7, 0.8, 0.9},
	}

	b := []float64{0.1, 0.3, 0.5}

	fmt.Println("\nМатрица A:")
	for i := 0; i < 3; i++ {
		fmt.Printf("  [%4.1f %4.1f %4.1f]\n", A[i][0], A[i][1], A[i][2])
	}

	fmt.Println("\nПравая часть b:")
	fmt.Printf("  b = [%4.1f %4.1f %4.1f]\n", b[0], b[1], b[2])

	// Вычисляем определитель
	det := A[0][0]*A[1][1]*A[2][2] + A[0][1]*A[1][2]*A[2][0] + A[0][2]*A[1][0]*A[2][1] -
		A[0][2]*A[1][1]*A[2][0] - A[0][0]*A[1][2]*A[2][1] - A[0][1]*A[1][0]*A[2][2]

	fmt.Printf("\nОпределитель матрицы: %.2e\n", det)

	// Проверка совместности - решаем первые два уравнения
	// 0.1*x1 + 0.2*x2 + 0.3*x3 = 0.1
	// 0.4*x1 + 0.5*x2 + 0.6*x3 = 0.3

	// Выражаем x1, x2 через x3
	// Из первого уравнения: 0.1*x1 = 0.1 - 0.2*x2 - 0.3*x3
	// Из второго: 0.4*x1 = 0.3 - 0.5*x2 - 0.6*x3

	// Умножаем первое на 4: 0.4*x1 = 0.4 - 0.8*x2 - 1.2*x3
	// Приравниваем: 0.4 - 0.8*x2 - 1.2*x3 = 0.3 - 0.5*x2 - 0.6*x3
	// 0.1 - 0.3*x2 - 0.6*x3 = 0
	// 0.3*x2 = 0.1 - 0.6*x3
	// x2 = (0.1 - 0.6*x3)/0.3 = 1/3 - 2*x3

	// Подставляем в первое: 0.1*x1 + 0.2*(1/3 - 2*x3) + 0.3*x3 = 0.1
	// 0.1*x1 + 0.2/3 - 0.4*x3 + 0.3*x3 = 0.1
	// 0.1*x1 + 0.2/3 - 0.1*x3 = 0.1
	// 0.1*x1 = 0.1 - 0.2/3 + 0.1*x3
	// x1 = 1 - 2/3 + x3 = 1/3 + x3

	fmt.Println("\nОбщее решение:")
	fmt.Println("  x1 = 1/3 + t")
	fmt.Println("  x2 = 1/3 - 2t")
	fmt.Println("  x3 = t")
	fmt.Println("где t - любое действительное число")

	// Проверка для t=0
	fmt.Println("\nПример решения (t=0):")
	fmt.Println("  x1 = 0.333333")
	fmt.Println("  x2 = 0.333333")
	fmt.Println("  x3 = 0")

	// Проверка для третьего уравнения
	x1 := 1.0 / 3.0
	x2 := 1.0 / 3.0
	x3 := 0.0
	check3 := 0.7*x1 + 0.8*x2 + 0.9*x3
	fmt.Printf("\nПроверка 3-го уравнения: 0.7*%.6f + 0.8*%.6f = %.6f, должно быть 0.5\n",
		x1, x2, check3)
	fmt.Printf("Разница: %.2e\n", math.Abs(check3-0.5))
}

func main() {
	task1()
	task2()
}
