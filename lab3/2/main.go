package main

import (
	"fmt"
)

const eps = 1e-20

func progonka(a, b, c, d []float64) (xVal, betaVal []float64) {
	n := len(b)

	alphaVal := make([]float64, n)
	betaVal = make([]float64, n)

	alphaVal[0] = -c[0] / b[0]
	betaVal[0] = d[0] / b[0]

	for i := 1; i < n; i++ {
		mulAv := a[i] * alphaVal[i-1]
		mulBv := a[i] * betaVal[i-1]
		denV := b[i] + mulAv
		numV := d[i] - mulBv
		alphaVal[i] = -c[i] / denV
		betaVal[i] = numV / denV
	}

	xVal = make([]float64, n)
	xVal[n-1] = betaVal[n-1]

	for i := n - 2; i >= 0; i-- {
		pV := alphaVal[i] * xVal[i+1]
		xVal[i] = pV + betaVal[i]
	}

	return xVal, betaVal
}

func main() {

	a := []float64{0, -2, 2, 1, 3}
	b := []float64{1, 4, -2, 1, -1}
	c := []float64{3, -1, 1, 1, 0}
	d := []float64{5, 1, 3, -2, -1}

	xVal, betaVal := progonka(a, b, c, d)

	fmt.Println("Коэффициенты бета:")
	for i := 0; i < len(betaVal); i++ {
		fmt.Println(i+1, betaVal[i])
	}

	fmt.Println("Решение системы:")
	for i := 0; i < len(xVal); i++ {
		fmt.Println(i+1, xVal[i])
	}
}
