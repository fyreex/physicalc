// Package algebra implementa operações de álgebra linear.
package algebra

import (
	"errors"
	"math"
)

// Vector representa um vetor n-dimensional.
type Vector []float64

// Len retorna a dimensão do vetor.
func (v Vector) Len() int { return len(v) }

// Add soma dois vetores componente a componente.
// Retorna erro se os vetores tiverem dimensões diferentes.
func Add(a, b Vector) (Vector, error) {
	if len(a) != len(b) {
		return nil, errors.New("vetores devem ter a mesma dimensão")
	}
	result := make(Vector, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result, nil
}

// Scale multiplica o vetor por um escalar.
func Scale(v Vector, scalar float64) Vector {
	result := make(Vector, len(v))
	for i, val := range v {
		result[i] = val * scalar
	}
	return result
}

// Dot calcula o produto escalar (dot product) de dois vetores.
// Resultado é um escalar: a·b = Σ aᵢ·bᵢ
func Dot(a, b Vector) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("vetores devem ter a mesma dimensão")
	}
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum, nil
}

// Norm calcula o módulo (magnitude) do vetor: ||v|| = √(Σ vᵢ²)
func Norm(v Vector) float64 {
	sum := 0.0
	for _, val := range v {
		sum += val * val
	}
	return math.Sqrt(sum)
}

// Normalize retorna o vetor unitário (módulo = 1) na mesma direção.
func Normalize(v Vector) (Vector, error) {
	n := Norm(v)
	if n == 0 {
		return nil, errors.New("não é possível normalizar o vetor nulo")
	}
	return Scale(v, 1/n), nil
}

// Cross calcula o produto vetorial (cross product) de dois vetores 3D.
// Resultado é um vetor perpendicular a ambos.
func Cross(a, b Vector) (Vector, error) {
	if len(a) != 3 || len(b) != 3 {
		return nil, errors.New("produto vetorial é definido apenas para vetores 3D")
	}
	return Vector{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}, nil
}

// Angle calcula o ângulo em radianos entre dois vetores.
// Usa: cos(θ) = (a·b) / (||a|| · ||b||)
func Angle(a, b Vector) (float64, error) {
	dot, err := Dot(a, b)
	if err != nil {
		return 0, err
	}
	na, nb := Norm(a), Norm(b)
	if na == 0 || nb == 0 {
		return 0, errors.New("ângulo indefinido para vetor nulo")
	}
	// Clamp para evitar erro de arredondamento no acos
	cos := dot / (na * nb)
	if cos > 1 {
		cos = 1
	}
	if cos < -1 {
		cos = -1
	}
	return math.Acos(cos), nil
}
