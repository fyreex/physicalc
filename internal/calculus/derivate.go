package calculus

// Derivative calcula a derivada numérica de f no ponto x
// usando diferenças centradas: f'(x) ≈ [f(x+h) - f(x-h)] / 2h
//
// h é o passo (ex: 1e-5). Valores menores aumentam precisão mas podem
// introduzir erros de ponto flutuante.
func Derivative(f Func, x, h float64) float64 {
	return (f(x+h) - f(x-h)) / (2 * h)
}

// SecondDerivative calcula a segunda derivada f''(x) usando diferenças centradas.
// Fórmula: f''(x) ≈ [f(x+h) - 2f(x) + f(x-h)] / h²
func SecondDerivative(f Func, x, h float64) float64 {
	return (f(x+h) - 2*f(x) + f(x-h)) / (h * h)
}

// PartialDerivativeX calcula ∂f/∂x para uma função de duas variáveis f(x,y).
func PartialDerivativeX(f func(x, y float64) float64, x, y, h float64) float64 {
	return (f(x+h, y) - f(x-h, y)) / (2 * h)
}

// PartialDerivativeY calcula ∂f/∂y para uma função de duas variáveis f(x,y).
func PartialDerivativeY(f func(x, y float64) float64, x, y, h float64) float64 {
	return (f(x, y+h) - f(x, y-h)) / (2 * h)
}
