package ode

import (
	"fmt"
	"math"
)

// ODEFunc é o tipo de uma equação diferencial de 1ª ordem: dy/dt = f(t, y)
type ODEFunc func(t, y float64) float64

// GetODEFunction retorna uma EDO predefinida pelo nome.
// Isso permite o usuário escolher EDOs comuns sem enviar código.
func GetODEFunction(name string) (ODEFunc, error) {
	functions := map[string]ODEFunc{
		// Decaimento exponencial: dy/dt = -y → solução: y = y0 * e^(-t)
		"decay": func(t, y float64) float64 { return -y },

		// Crescimento exponencial: dy/dt = y → solução: y = y0 * e^t
		"growth": func(t, y float64) float64 { return y },

		// Crescimento logístico: dy/dt = r*y*(1 - y/K), r=1, K=1
		"logistic": func(t, y float64) float64 { return y * (1 - y) },

		// Oscilação (pêndulo linearizado): dy/dt = cos(t)
		"oscillation": func(t, y float64) float64 { return math.Cos(t) },

		// Resfriamento de Newton: dy/dt = -(y - 25) (temperatura ambiente = 25°C)
		"cooling": func(t, y float64) float64 { return -(y - 25) },

		// Queda com resistência do ar: dy/dt = g - k*y, g=9.81, k=0.1
		"freefall": func(t, y float64) float64 { return 9.81 - 0.1*y },
	}

	fn, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("EDO '%s' não encontrada. Disponíveis: decay, growth, logistic, oscillation, cooling, freefall", name)
	}
	return fn, nil
}

// Euler resolve uma EDO dy/dt = f(t, y) usando o Método de Euler (1ª ordem).
//
// O método mais simples para EDOs. Menos preciso, mas fácil de entender.
// Fórmula: y(t+h) = y(t) + h * f(t, y(t))
//
// Retorna dois slices paralelos: os valores de t e os valores de y.
func Euler(f ODEFunc, y0, t0, tEnd, h float64) ([]float64, []float64) {
	n := int((tEnd-t0)/h) + 1

	ts := make([]float64, n)
	ys := make([]float64, n)

	ts[0] = t0
	ys[0] = y0

	for i := 1; i < n; i++ {
		t := ts[i-1]
		y := ys[i-1]

		// Passo de Euler: avança a solução em h usando a derivada atual
		ys[i] = y + h*f(t, y)
		ts[i] = t + h
	}

	return ts, ys
}

// RungeKutta4 resolve uma EDO usando Runge-Kutta de 4ª Ordem (RK4).
//
// Muito mais preciso que Euler. Calcula 4 "inclinações" (k1, k2, k3, k4)
// e combina com pesos para avançar a solução.
//
// Erro de truncamento local: O(h⁵) vs O(h²) do Euler.
//
// Fórmulas:
//
//	k1 = f(t, y)
//	k2 = f(t + h/2, y + h*k1/2)
//	k3 = f(t + h/2, y + h*k2/2)
//	k4 = f(t + h, y + h*k3)
//	y(t+h) = y + (h/6) * (k1 + 2k2 + 2k3 + k4)
func RungeKutta4(f ODEFunc, y0, t0, tEnd, h float64) ([]float64, []float64) {
	n := int((tEnd-t0)/h) + 1

	ts := make([]float64, n)
	ys := make([]float64, n)

	ts[0] = t0
	ys[0] = y0

	for i := 1; i < n; i++ {
		t := ts[i-1]
		y := ys[i-1]

		// Quatro estimativas de inclinação
		k1 := f(t, y)
		k2 := f(t+h/2, y+h*k1/2)
		k3 := f(t+h/2, y+h*k2/2)
		k4 := f(t+h, y+h*k3)

		// Média ponderada: k1 e k4 têm peso 1, k2 e k3 têm peso 2
		ys[i] = y + (h/6)*(k1+2*k2+2*k3+k4)
		ts[i] = t + h
	}

	return ts, ys
}
