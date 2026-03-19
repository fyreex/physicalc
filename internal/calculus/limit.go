package calculus

import "math"

func NumericalLimit(f MathFunc, target float64) (left, right float64) {
	h := 1e-7
	left = f(target - h)
	right = f(target + h)
	return left, right
}

func Converges(left, right, epsilon float64) bool {
	return math.Abs(left-right) < epsilon
}
