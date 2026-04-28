package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type CubicSpline struct {
	xs         []float64
	a, b, c, d []float64
	n          int
}

func NewCubicSpline(xs, ys []float64) *CubicSpline {
	n := len(xs) - 1

	h := make([]float64, n)
	for i := 0; i < n; i++ {
		h[i] = xs[i+1] - xs[i]
	}

	c := make([]float64, n+1)

	m := n - 1
	if m > 0 {
		diag := make([]float64, m)
		lower := make([]float64, m)
		upper := make([]float64, m)
		rhs := make([]float64, m)

		for i := 0; i < m; i++ {
			diag[i] = 2 * (h[i] + h[i+1])
			rhs[i] = 3 * ((ys[i+2]-ys[i+1])/h[i+1] - (ys[i+1]-ys[i])/h[i])
			if i > 0 {
				lower[i] = h[i]
			}
			if i < m-1 {
				upper[i] = h[i+1]
			}
		}

		cInner := thomas(lower, diag, upper, rhs)
		for i := 0; i < m; i++ {
			c[i+1] = cInner[i]
		}
	}

	a := make([]float64, n)
	b := make([]float64, n)
	d := make([]float64, n)
	cArr := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = ys[i]
		cArr[i] = c[i]
		d[i] = (c[i+1] - c[i]) / (3 * h[i])
		b[i] = (ys[i+1]-ys[i])/h[i] - h[i]*(2*c[i]+c[i+1])/3
	}

	return &CubicSpline{xs: xs, a: a, b: b, c: cArr, d: d, n: n}
}

func thomas(lower, diag, upper, rhs []float64) []float64 {
	n := len(diag)
	w := make([]float64, n)
	g := make([]float64, n)
	x := make([]float64, n)

	w[0] = diag[0]
	g[0] = rhs[0]
	for i := 1; i < n; i++ {
		mu := lower[i] / w[i-1]
		w[i] = diag[i] - mu*upper[i-1]
		g[i] = rhs[i] - mu*g[i-1]
	}
	x[n-1] = g[n-1] / w[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = (g[i] - upper[i]*x[i+1]) / w[i]
	}
	return x
}

func (s *CubicSpline) Eval(x float64) float64 {
	if x <= s.xs[0] {
		dx := x - s.xs[0]
		return s.a[0] + s.b[0]*dx + s.c[0]*dx*dx + s.d[0]*dx*dx*dx
	}
	if x >= s.xs[s.n] {
		i := s.n - 1
		dx := x - s.xs[i]
		return s.a[i] + s.b[i]*dx + s.c[i]*dx*dx + s.d[i]*dx*dx*dx
	}
	lo, hi := 0, s.n-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if s.xs[mid] <= x {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	dx := x - s.xs[lo]
	return s.a[lo] + s.b[lo]*dx + s.c[lo]*dx*dx + s.d[lo]*dx*dx*dx
}

func f(x float64) float64 {
	return 1.0 / (1.0 + 25*x*x)
}

func task1() {
	fmt.Println("=== Задание 1: f(x) = 1/(1+25x²) на [-1, 1] ===")
	fmt.Println()

	nValues := []int{5, 10, 20}

	ptsErr := make(plotter.XYs, len(nValues))

	for idx, n := range nValues {
		xs := make([]float64, n+1)
		ys := make([]float64, n+1)
		for i := 0; i <= n; i++ {
			xs[i] = -1 + 2*float64(i)/float64(n)
			ys[i] = f(xs[i])
		}

		spline := NewCubicSpline(xs, ys)

		const testPts = 1000
		maxErr := 0.0
		for i := 0; i <= testPts; i++ {
			x := -1 + 2*float64(i)/float64(testPts)
			err := math.Abs(spline.Eval(x) - f(x))
			if err > maxErr {
				maxErr = err
			}
		}
		fmt.Printf("  n = %2d: max|S(x) − f(x)| = %.8f\n", n, maxErr)
		ptsErr[idx].X = float64(n)
		ptsErr[idx].Y = maxErr

		buildTask1SplinePlot(n, xs, ys, spline)
	}

	buildErrorPlot(ptsErr)
	fmt.Println()
}

func buildTask1SplinePlot(n int, xs, ys []float64, spline *CubicSpline) {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Кубический сплайн vs f(x) = 1/(1+25x²),  n = %d", n)
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	const pts = 600

	ptsF := make(plotter.XYs, pts+1)
	for i := 0; i <= pts; i++ {
		x := -1 + 2*float64(i)/float64(pts)
		ptsF[i] = plotter.XY{X: x, Y: f(x)}
	}
	lineF, _ := plotter.NewLine(ptsF)
	lineF.Color = color.RGBA{R: 50, G: 100, B: 220, A: 255}
	lineF.Width = vg.Points(2)

	ptsS := make(plotter.XYs, pts+1)
	for i := 0; i <= pts; i++ {
		x := -1 + 2*float64(i)/float64(pts)
		ptsS[i] = plotter.XY{X: x, Y: spline.Eval(x)}
	}
	lineS, _ := plotter.NewLine(ptsS)
	lineS.Color = color.RGBA{R: 220, G: 50, B: 50, A: 255}
	lineS.Width = vg.Points(2)
	lineS.Dashes = []vg.Length{vg.Points(6), vg.Points(2)}

	nodes := make(plotter.XYs, len(xs))
	for i, x := range xs {
		nodes[i] = plotter.XY{X: x, Y: ys[i]}
	}
	scatter, _ := plotter.NewScatter(nodes)
	scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 160, B: 0, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(4)

	p.Add(lineF, lineS, scatter)
	p.Legend.Add("f(x)", lineF)
	p.Legend.Add(fmt.Sprintf("Сплайн (n=%d)", n), lineS)
	p.Legend.Add("Узлы", scatter)
	p.Legend.Top = true

	fname := fmt.Sprintf("task1_n%d.png", n)
	if err := p.Save(10*vg.Inch, 6*vg.Inch, fname); err != nil {
		fmt.Printf("  [!] Ошибка сохранения %s: %v\n", fname, err)
	} else {
		fmt.Printf("  График: %s\n", fname)
	}
}

func buildErrorPlot(ptsErr plotter.XYs) {
	p := plot.New()
	p.Title.Text = "Максимальное отклонение сплайна от f(x)"
	p.X.Label.Text = "n (число узлов)"
	p.Y.Label.Text = "max|S(x) − f(x)|"

	line, _ := plotter.NewLine(ptsErr)
	line.Color = color.RGBA{R: 180, G: 0, B: 180, A: 255}
	line.Width = vg.Points(2)

	scatter, _ := plotter.NewScatter(ptsErr)
	scatter.GlyphStyle.Color = color.RGBA{R: 180, G: 0, B: 180, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(5)

	p.Add(line, scatter)
	p.Legend.Add("max|S−f|", line)
	p.Legend.Top = true

	fname := "task1_error.png"
	if err := p.Save(8*vg.Inch, 5*vg.Inch, fname); err != nil {
		fmt.Printf("  [!] Ошибка сохранения %s: %v\n", fname, err)
	} else {
		fmt.Printf("  График отклонения: %s\n", fname)
	}
}

func task2() {
	fmt.Println("=== Задание 2: Кубический сплайн по таблице ===")
	fmt.Println()

	xs := []float64{2, 3, 5, 7}
	ys := []float64{4, -2, 6, -3}

	fmt.Printf("  %-6s", "x:")
	for _, x := range xs {
		fmt.Printf("  %6.1f", x)
	}
	fmt.Printf("\n  %-6s", "f(x):")
	for _, y := range ys {
		fmt.Printf("  %6.1f", y)
	}
	fmt.Println("\n")

	spline := NewCubicSpline(xs, ys)

	fmt.Println("  Коэффициенты (S_i(x) = a + b·Δx + c·Δx² + d·Δx³, Δx = x − xᵢ):")
	fmt.Printf("  %-4s  %-10s  %-14s  %-14s  %-14s  %-14s  %s\n",
		"i", "xᵢ", "a", "b", "c", "d", "отрезок")
	fmt.Println("  " + strings.Repeat("─", 90))
	for i := 0; i < spline.n; i++ {
		fmt.Printf("  %-4d  %-10.4f  %-14.6f  %-14.6f  %-14.6f  %-14.6f  [%.1f, %.1f]\n",
			i+1, xs[i],
			spline.a[i], spline.b[i], spline.c[i], spline.d[i],
			xs[i], xs[i+1])
	}

	fmt.Println("\n  Проверка в узловых точках:")
	fmt.Printf("  %-8s  %-14s  %-14s  %-14s\n", "x", "f(x)", "S(x)", "|S(x)−f(x)|")
	fmt.Println("  " + strings.Repeat("─", 55))
	for i, x := range xs {
		sx := spline.Eval(x)
		err := math.Abs(sx - ys[i])
		fmt.Printf("  %-8.1f  %-14.6f  %-14.6f  %-14.2e\n", x, ys[i], sx, err)
	}

	buildTask2Plot(xs, ys, spline)
	fmt.Println()
}

func buildTask2Plot(xs, ys []float64, spline *CubicSpline) {
	p := plot.New()
	p.Title.Text = "Кубический сплайн (табличная функция)"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	xMin, xMax := xs[0], xs[len(xs)-1]
	const pts = 600

	ptsS := make(plotter.XYs, pts+1)
	for i := 0; i <= pts; i++ {
		x := xMin + (xMax-xMin)*float64(i)/float64(pts)
		ptsS[i] = plotter.XY{X: x, Y: spline.Eval(x)}
	}
	lineS, _ := plotter.NewLine(ptsS)
	lineS.Color = color.RGBA{R: 50, G: 100, B: 220, A: 255}
	lineS.Width = vg.Points(2.5)

	nodes := make(plotter.XYs, len(xs))
	for i, x := range xs {
		nodes[i] = plotter.XY{X: x, Y: ys[i]}
	}
	scatter, _ := plotter.NewScatter(nodes)
	scatter.GlyphStyle.Color = color.RGBA{R: 220, G: 50, B: 50, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(6)

	p.Add(lineS, scatter)
	p.Legend.Add("Кубический сплайн", lineS)
	p.Legend.Add("Узловые точки", scatter)
	p.Legend.Top = true

	fname := "task2_spline.png"
	if err := p.Save(10*vg.Inch, 6*vg.Inch, fname); err != nil {
		fmt.Printf("  [!] Ошибка сохранения %s: %v\n", fname, err)
	} else {
		fmt.Printf("  График: %s\n", fname)
	}
}

func main() {
	task1()
	task2()
}
