package calculus_test

import (
	"math"
	"testing"

	"github.com/seuuser/physicalc/internal/calculus"
	"github.com/stretchr/testify/assert"
)

// TestSimpson verifica que a integral de x² de 0 a 1 ≈ 1/3
func TestSimpson(t *testing.T) {
	f := func(x float64) float64 { return x * x }
	result := calculus.Simpson(f, 0, 1, 1000)
	assert.InDelta(t, 1.0/3.0, result, 1e-6, "∫x² de 0 a 1 deve ser ≈ 0.3333")
}

// TestTrapezoid verifica que a integral de x² de 0 a 1 ≈ 1/3
func TestTrapezoid(t *testing.T) {
	f := func(x float64) float64 { return x * x }
	result := calculus.Trapezoid(f, 0, 1, 10000)
	assert.InDelta(t, 1.0/3.0, result, 1e-4, "∫x² de 0 a 1 deve ser ≈ 0.3333")
}

// TestIntegralSin verifica que ∫sin(x) de 0 a π ≈ 2
func TestIntegralSin(t *testing.T) {
	result := calculus.Simpson(math.Sin, 0, math.Pi, 1000)
	assert.InDelta(t, 2.0, result, 1e-8, "∫sin(x) de 0 a π deve ser ≈ 2")
}

// TestCentralDifference verifica f'(x²) em x=3 ≈ 6
func TestCentralDifference(t *testing.T) {
	f := func(x float64) float64 { return x * x }
	result := calculus.CentralDifference(f, 3, 1e-5)
	assert.InDelta(t, 6.0, result, 1e-8, "derivada de x² em x=3 deve ser 6")
}

// TestSecondDerivative verifica f”(x²) = 2
func TestSecondDerivative(t *testing.T) {
	f := func(x float64) float64 { return x * x }
	result := calculus.SecondDerivative(f, 5, 1e-4)
	assert.InDelta(t, 2.0, result, 1e-4, "segunda derivada de x² deve ser 2")
}

// TestLimit verifica o limite de sin(x)/x quando x→0 ≈ 1
func TestLimit(t *testing.T) {
	f := func(x float64) float64 { return math.Sin(x) / x }
	left, right := calculus.NumericalLimit(f, 0)
	assert.InDelta(t, 1.0, left, 1e-5)
	assert.InDelta(t, 1.0, right, 1e-5)
	assert.True(t, calculus.Converges(left, right, 1e-6))
}
