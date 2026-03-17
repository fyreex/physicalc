package algebra_test

import (
	"math"
	"testing"

	"github.com/fyreex/physicalc/internal/algebra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVectorAdd(t *testing.T) {
	a := algebra.NewVector([]float64{1, 2, 3})
	b := algebra.NewVector([]float64{4, 5, 6})
	result, err := a.Add(b)
	require.NoError(t, err)
	assert.Equal(t, []float64{5, 7, 9}, result.Components)
}

func TestVectorDot(t *testing.T) {
	a := algebra.NewVector([]float64{1, 2, 3})
	b := algebra.NewVector([]float64{4, 5, 6})
	result, err := a.Dot(b)
	require.NoError(t, err)
	// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	assert.Equal(t, 32.0, result)
}

func TestVectorCross(t *testing.T) {
	a := algebra.NewVector([]float64{1, 0, 0})
	b := algebra.NewVector([]float64{0, 1, 0})
	result, err := a.Cross(b)
	require.NoError(t, err)
	// i × j = k = [0,0,1]
	assert.InDelta(t, 0.0, result.Components[0], 1e-10)
	assert.InDelta(t, 0.0, result.Components[1], 1e-10)
	assert.InDelta(t, 1.0, result.Components[2], 1e-10)
}

func TestVectorNorm(t *testing.T) {
	v := algebra.NewVector([]float64{3, 4})
	// ||[3,4]|| = 5 (triângulo 3-4-5)
	assert.InDelta(t, 5.0, v.Norm(), 1e-10)
}

func TestVectorAngle(t *testing.T) {
	a := algebra.NewVector([]float64{1, 0})
	b := algebra.NewVector([]float64{0, 1})
	angle, err := a.AngleDeg(b)
	require.NoError(t, err)
	// vetores perpendiculares = 90°
	assert.InDelta(t, 90.0, angle, 1e-10)
}

func TestMatrixMultiply(t *testing.T) {
	a := algebra.NewMatrix([][]float64{{1, 2}, {3, 4}})
	b := algebra.NewMatrix([][]float64{{5, 6}, {7, 8}})
	result, err := a.Multiply(b)
	require.NoError(t, err)
	assert.Equal(t, [][]float64{{19, 22}, {43, 50}}, result.Data)
}

func TestMatrixDeterminant(t *testing.T) {
	m := algebra.NewMatrix([][]float64{{1, 2}, {3, 4}})
	det, err := m.Determinant()
	require.NoError(t, err)
	// det = 1*4 - 2*3 = -2
	assert.InDelta(t, -2.0, det, 1e-10)
}

func TestMatrixInverse(t *testing.T) {
	m := algebra.NewMatrix([][]float64{{1, 2}, {3, 4}})
	inv, err := m.Inverse()
	require.NoError(t, err)

	// Verifica A * A⁻¹ ≈ I
	identity, err := m.Multiply(inv)
	require.NoError(t, err)
	assert.InDelta(t, 1.0, identity.Data[0][0], 1e-10)
	assert.InDelta(t, 0.0, identity.Data[0][1], 1e-10)
	assert.InDelta(t, 0.0, identity.Data[1][0], 1e-10)
	assert.InDelta(t, 1.0, identity.Data[1][1], 1e-10)
}

func TestMatrixTranspose(t *testing.T) {
	m := algebra.NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}})
	t2 := m.Transpose()
	assert.Equal(t, 3, t2.Rows)
	assert.Equal(t, 2, t2.Cols)
	assert.Equal(t, 1.0, t2.Data[0][0])
	assert.Equal(t, 4.0, t2.Data[0][1])
}

func TestVectorNormalize(t *testing.T) {
	v := algebra.NewVector([]float64{3, 4})
	u, err := v.Normalize()
	require.NoError(t, err)
	// Norma do vetor unitário deve ser 1
	assert.InDelta(t, 1.0, u.Norm(), 1e-10)
	assert.InDelta(t, 3.0/5.0, u.Components[0], 1e-10)
	assert.InDelta(t, 4.0/5.0, u.Components[1], 1e-10)
}

func TestParallelVectorsAngle(t *testing.T) {
	a := algebra.NewVector([]float64{1, 0, 0})
	b := algebra.NewVector([]float64{2, 0, 0})
	angle, err := a.AngleDeg(b)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, angle, 1e-10)
}

func TestMatrixSingular(t *testing.T) {
	// Matriz singular (det = 0) não tem inversa
	m := algebra.NewMatrix([][]float64{{1, 2}, {2, 4}})
	_, err := m.Inverse()
	assert.Error(t, err)
}

func TestVectorDimensionMismatch(t *testing.T) {
	a := algebra.NewVector([]float64{1, 2})
	b := algebra.NewVector([]float64{1, 2, 3})
	_, err := a.Add(b)
	assert.Error(t, err)
}

func TestDeterminant3x3(t *testing.T) {
	m := algebra.NewMatrix([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
	det, err := m.Determinant()
	require.NoError(t, err)
	// Matriz singular, det = 0
	assert.InDelta(t, 0.0, math.Abs(det), 1e-8)
}
