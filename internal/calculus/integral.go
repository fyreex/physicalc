package calculus

import (
	"fmt"
	"math"
)

type MathFunc func(x float64) float64

func GetFunction(name string) (MathFunc, error) {
	functions := map[string]MathFunc{
		"x2":    func(x float64) float64 { return x * x },
		"x3":    func(x float64) float64 { return x * x * x },
		"sin":   math.Sin,
		"cos":   math.Cos,
		"tan":   math.Tan,
		"exp":   math.Exp,
		"ln":    math.Log,
		"sqrt":  math.Sqrt,
		"1/x":   func(x float64) float64 { return 1.0 / x },
		"x2+1":  func(x float64) float64 { return x*x + 1 },
		"sinx2": func(x float64) float64 { return math.Sin(x * x) },
	}
	fn, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("funcao '%s' nao encontrada", name)
	}
	return fn, nil
}

func Simpson(f MathFunc, a, b float64, n int) float64 {
	if n%2 != 0 {
		n++
	}
	h := (b - a) / float64(n)
	sum := f(a) + f(b)
	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		if i%2 == 0 {
			sum += 2 * f(x)
		} else {
			sum += 4 * f(x)
		}
	}
	return (h / 3) * sum
}

func Trapezoid(f MathFunc, a, b float64, n int) float64 {
	h := (b - a) / float64(n)
	sum := f(a) + f(b)
	for i := 1; i < n; i++ {
		sum += 2 * f(a+float64(i)*h)
	}
	return (h / 2) * sum
}
