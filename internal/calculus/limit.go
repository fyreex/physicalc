package calculus

import "math"

// LimitApproach calcula o limite de f quando x → target vindo dos dois lados.
// Retorna o valor do limite e se ele existe (lim esquerdo ≈ lim direito).
//
// Usa uma sequência de passos decrescentes para verificar convergência.
func LimitApproach(f Func, target, tolerance float64) (value float64, exists bool) {
	steps := []float64{1e-3, 1e-4, 1e-5, 1e-6}

	var leftVal, rightVal float64

	for _, h := range steps {
		leftVal = f(target - h)
		rightVal = f(target + h)
	}

	// Verifica se ambos os lados convergem para o mesmo valor
	if math.IsNaN(leftVal) || math.IsInf(leftVal, 0) ||
		math.IsNaN(rightVal) || math.IsInf(rightVal, 0) {
		return 0, false
	}

	exists = math.Abs(leftVal-rightVal) < tolerance
	value = (leftVal + rightVal) / 2
	return
}
