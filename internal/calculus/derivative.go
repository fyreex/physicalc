package calculus

func CentralDifference(f MathFunc, x, h float64) float64 {
	return (f(x+h) - f(x-h)) / (2 * h)
}

func ForwardDifference(f MathFunc, x, h float64) float64 {
	return (f(x+h) - f(x)) / h
}

func SecondDerivative(f MathFunc, x, h float64) float64 {
	return (f(x+h) - 2*f(x) + f(x-h)) / (h * h)
}

func PartialX(f func(x, y float64) float64, x, y, h float64) float64 {
	return (f(x+h, y) - f(x-h, y)) / (2 * h)
}

func PartialY(f func(x, y float64) float64, x, y, h float64) float64 {
	return (f(x, y+h) - f(x, y-h)) / (2 * h)
}
