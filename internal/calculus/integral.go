// Package calculus implementa métodos numéricos de cálculo diferencial e integral.
package calculus

import "errors"

// Func é o tipo de uma função matemática f(x).
type Func func(x float64) float64

// IntegrateSimpsons calcula a integral definida de f no intervalo [a, b]
// usando a Regra de Simpson 1/3.
//
// Parâmetros:
//   - f: função a integrar
//   - a, b: limites do intervalo
//   - n: número de subintervalos (deve ser par e > 0)
//
// Erro é retornado se n for ímpar ou <= 0.
func IntegrateSimpsons(f Func, a, b float64, n int) (float64, error) {
	if n <= 0 || n%2 != 0 {
		return 0, errors.New("n deve ser um número par positivo")
	}

	h := (b - a) / float64(n) // largura de cada subintervalo
	sum := f(a) + f(b)        // f(x0) + f(xn)

	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		if i%2 == 0 {
			sum += 2 * f(x) // coeficiente 2 para pontos pares
		} else {
			sum += 4 * f(x) // coeficiente 4 para pontos ímpares
		}
	}

	return (h / 3) * sum, nil
}

// IntegrateTrapezoid calcula a integral definida usando a Regra dos Trapézios.
// Menos preciso que Simpson, mas funciona com qualquer n.
func IntegrateTrapezoid(f Func, a, b float64, n int) (float64, error) {
	if n <= 0 {
		return 0, errors.New("n deve ser positivo")
	}

	h := (b - a) / float64(n)
	sum := (f(a) + f(b)) / 2

	for i := 1; i < n; i++ {
		sum += f(a + float64(i)*h)
	}

	return h * sum, nil
}
