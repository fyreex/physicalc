package algebra

import (
	"errors"
	"math"
)

type Vector struct {
	Components []float64
	Dim        int
}

func NewVector(components []float64) Vector {
	return Vector{Components: components, Dim: len(components)}
}

func checkSameDim(a, b Vector) error {
	if a.Dim != b.Dim {
		return errors.New("vetores devem ter a mesma dimensão")
	}
	return nil
}

func (v Vector) Add(other Vector) (Vector, error) {
	if err := checkSameDim(v, other); err != nil {
		return Vector{}, err
	}
	result := make([]float64, v.Dim)
	for i := range result {
		result[i] = v.Components[i] + other.Components[i]
	}
	return NewVector(result), nil
}

func (v Vector) Sub(other Vector) (Vector, error) {
	if err := checkSameDim(v, other); err != nil {
		return Vector{}, err
	}
	result := make([]float64, v.Dim)
	for i := range result {
		result[i] = v.Components[i] - other.Components[i]
	}
	return NewVector(result), nil
}

func (v Vector) Scale(s float64) Vector {
	result := make([]float64, v.Dim)
	for i := range result {
		result[i] = v.Components[i] * s
	}
	return NewVector(result)
}

func (v Vector) Dot(other Vector) (float64, error) {
	if err := checkSameDim(v, other); err != nil {
		return 0, err
	}
	sum := 0.0
	for i := range v.Components {
		sum += v.Components[i] * other.Components[i]
	}
	return sum, nil
}

func (v Vector) Cross(other Vector) (Vector, error) {
	if v.Dim != 3 || other.Dim != 3 {
		return Vector{}, errors.New("produto vetorial requer vetores 3D")
	}
	a, b := v.Components, other.Components
	return NewVector([]float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}), nil
}

func (v Vector) Norm() float64 {
	sum := 0.0
	for _, c := range v.Components {
		sum += c * c
	}
	return math.Sqrt(sum)
}

func (v Vector) Normalize() (Vector, error) {
	n := v.Norm()
	if n == 0 {
		return Vector{}, errors.New("não é possível normalizar o vetor zero")
	}
	return v.Scale(1 / n), nil
}

func (v Vector) AngleDeg(other Vector) (float64, error) {
	dot, err := v.Dot(other)
	if err != nil {
		return 0, err
	}
	denom := v.Norm() * other.Norm()
	if denom == 0 {
		return 0, errors.New("ângulo indefinido com vetor zero")
	}
	cosTheta := dot / denom
	if cosTheta > 1 {
		cosTheta = 1
	}
	if cosTheta < -1 {
		cosTheta = -1
	}
	return math.Acos(cosTheta) * (180 / math.Pi), nil
}
