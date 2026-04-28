package main

import (
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	epsilon = 0.000001
	maxIter = 1000
	n       = 4
)

var A = [n][n]float64{
	{12.14, 1.32, -0.78, -2.75},
	{-0.89, 16.75, 1.88, -1.55},
	{2.65, -1.27, -15.64, -0.64},
	{2.44, 1.52, 1.93, -11.43},
}

var b = [n]float64{14.78, -12.14, -11.65, 4.26}

func checkDiagonalDominance() {
	fmt.Println("=== Проверка условия сходимости ===")
	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			if i != j {
				sum += math.Abs(A[i][j])
			}
		}
		diag := math.Abs(A[i][i])
		status := "✓"
		if diag <= sum {
			status = "✗"
		}
		fmt.Printf("  Строка %d: |%.2f| > %.4f  %s\n", i+1, A[i][i], sum, status)
	}
	fmt.Println()
}

func norm(xNew, xOld [n]float64) float64 {
	maxVal := 0.0
	for i := 0; i < n; i++ {
		if d := math.Abs(xNew[i] - xOld[i]); d > maxVal {
			maxVal = d
		}
	}
	return maxVal
}

func residualNorm(x [n]float64) float64 {
	maxVal := 0.0
	for i := 0; i < n; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			sum += A[i][j] * x[j]
		}
		if d := math.Abs(sum - b[i]); d > maxVal {
			maxVal = d
		}
	}
	return maxVal
}

func jacobiMethod(x0 [n]float64) ([n]float64, int, []float64) {
	x := x0
	var residuals []float64
	for iter := 1; iter <= maxIter; iter++ {
		var xNew [n]float64
		for i := 0; i < n; i++ {
			sum := b[i]
			for j := 0; j < n; j++ {
				if i != j {
					sum -= A[i][j] * x[j]
				}
			}
			xNew[i] = sum / A[i][i]
		}
		delta := norm(xNew, x)
		residuals = append(residuals, residualNorm(xNew))
		x = xNew
		if delta <= epsilon {
			return x, iter, residuals
		}
	}
	return x, maxIter, residuals
}

func seidelMethod(x0 [n]float64) ([n]float64, int, []float64) {
	x := x0
	var residuals []float64
	for iter := 1; iter <= maxIter; iter++ {
		xOld := x
		for i := 0; i < n; i++ {
			sum := b[i]
			for j := 0; j < n; j++ {
				if i != j {
					sum -= A[i][j] * x[j]
				}
			}
			x[i] = sum / A[i][i]
		}
		delta := norm(x, xOld)
		residuals = append(residuals, residualNorm(x))
		if delta <= epsilon {
			return x, iter, residuals
		}
	}
	return x, maxIter, residuals
}

func printResult(name string, x [n]float64, iter int, residuals []float64) {
	fmt.Printf("  [%s]\n", name)
	fmt.Printf("  Итераций: %d\n", iter)
	for i := 0; i < n; i++ {
		fmt.Printf("  x%d = %.6f\n", i+1, x[i])
	}
	fmt.Printf("  Норма невязки: %.8f\n\n", residuals[len(residuals)-1])
}

// buildChart строит PNG график норм невязки для одного начального приближения
func buildChart(initLabel string, jacobiRes, seidelRes []float64, filename string) error {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Норма невязки — начальное приближение %s", initLabel)
	p.X.Label.Text = "Номер итерации"
	p.Y.Label.Text = "Норма невязки"

	// Якоби
	ptsJ := make(plotter.XYs, len(jacobiRes))
	for i, v := range jacobiRes {
		ptsJ[i].X = float64(i + 1)
		ptsJ[i].Y = v
	}
	lineJ, err := plotter.NewLine(ptsJ)
	if err != nil {
		return err
	}
	lineJ.Color = color.RGBA{R: 220, G: 50, B: 50, A: 255}
	lineJ.Width = vg.Points(2)

	// Зейдель
	ptsS := make(plotter.XYs, len(seidelRes))
	for i, v := range seidelRes {
		ptsS[i].X = float64(i + 1)
		ptsS[i].Y = v
	}
	lineS, err := plotter.NewLine(ptsS)
	if err != nil {
		return err
	}
	lineS.Color = color.RGBA{R: 50, G: 100, B: 220, A: 255}
	lineS.Width = vg.Points(2)

	p.Add(lineJ, lineS)
	p.Legend.Add("Якоби", lineJ)
	p.Legend.Add("Зейдель", lineS)
	p.Legend.Top = true

	// Пунктирная линия точности epsilon
	epsLine := plotter.NewFunction(func(x float64) float64 { return epsilon })
	epsLine.Color = color.RGBA{R: 80, G: 180, B: 80, A: 200}
	epsLine.Dashes = []vg.Length{vg.Points(4), vg.Points(4)}
	epsLine.Width = vg.Points(1.5)
	p.Add(epsLine)
	p.Legend.Add(fmt.Sprintf("ε = %.3f", epsilon), epsLine)

	return p.Save(10*vg.Inch, 5*vg.Inch, filename)
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║   Решение СЛАУ методами Якоби и Зейделя             ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Printf("\nТочность: ε = %.3f\n\n", epsilon)

	checkDiagonalDominance()

	initials := [][n]float64{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{-1, 2, -1, 0},
	}
	labels := []string{"[0,0,0,0]", "[1,1,1,1]", "[-1,2,-1,0]"}

	for idx, x0 := range initials {
		fmt.Printf("════════════════════════════════════════\n")
		fmt.Printf("Начальное приближение #%d: %s\n\n", idx+1, labels[idx])

		xJ, iterJ, resJ := jacobiMethod(x0)
		printResult("Метод Якоби", xJ, iterJ, resJ)

		xS, iterS, resS := seidelMethod(x0)
		printResult("Метод Зейделя", xS, iterS, resS)

		fmt.Printf("  → Якоби: %d ит. | Зейдель: %d ит.\n\n", iterJ, iterS)

		fname := fmt.Sprintf("chart_%d.png", idx+1)
		if err := buildChart(labels[idx], resJ, resS, fname); err != nil {
			fmt.Printf("  Ошибка построения графика: %v\n", err)
		} else {
			fmt.Printf("  График сохранён: %s\n\n", fname)
		}
	}
}
