package algebra

import (
	"errors"
	"math"
)

type Matrix struct {
	Data [][]float64
	Rows int
	Cols int
}

func NewMatrix(data [][]float64) Matrix {
	rows := len(data)
	cols := 0
	if rows > 0 {
		cols = len(data[0])
	}
	return Matrix{Data: data, Rows: rows, Cols: cols}
}

func (m Matrix) Multiply(other Matrix) (Matrix, error) {
	if m.Cols != other.Rows {
		return Matrix{}, errors.New("dimensões incompatíveis para multiplicação")
	}
	result := make([][]float64, m.Rows)
	for i := range result {
		result[i] = make([]float64, other.Cols)
		for j := range result[i] {
			for k := 0; k < m.Cols; k++ {
				result[i][j] += m.Data[i][k] * other.Data[k][j]
			}
		}
	}
	return NewMatrix(result), nil
}

func (m Matrix) Transpose() Matrix {
	result := make([][]float64, m.Cols)
	for i := range result {
		result[i] = make([]float64, m.Rows)
		for j := 0; j < m.Rows; j++ {
			result[i][j] = m.Data[j][i]
		}
	}
	return NewMatrix(result)
}

func (m Matrix) Trace() (float64, error) {
	if m.Rows != m.Cols {
		return 0, errors.New("traço requer matriz quadrada")
	}
	sum := 0.0
	for i := 0; i < m.Rows; i++ {
		sum += m.Data[i][i]
	}
	return sum, nil
}

func (m Matrix) Determinant() (float64, error) {
	if m.Rows != m.Cols {
		return 0, errors.New("determinante requer matriz quadrada")
	}
	return det(m.Data, m.Rows), nil
}

func det(data [][]float64, n int) float64 {
	if n == 1 {
		return data[0][0]
	}
	if n == 2 {
		return data[0][0]*data[1][1] - data[0][1]*data[1][0]
	}
	result := 0.0
	for col := 0; col < n; col++ {
		sign := math.Pow(-1, float64(col))
		result += sign * data[0][col] * det(minor(data, 0, col, n), n-1)
	}
	return result
}

func minor(data [][]float64, row, col, n int) [][]float64 {
	result := make([][]float64, n-1)
	ri := 0
	for i := 0; i < n; i++ {
		if i == row {
			continue
		}
		result[ri] = make([]float64, n-1)
		ci := 0
		for j := 0; j < n; j++ {
			if j == col {
				continue
			}
			result[ri][ci] = data[i][j]
			ci++
		}
		ri++
	}
	return result
}

func (m Matrix) Inverse() (Matrix, error) {
	if m.Rows != m.Cols {
		return Matrix{}, errors.New("inversa requer matriz quadrada")
	}
	n := m.Rows
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, 2*n)
		copy(aug[i], m.Data[i])
		aug[i][n+i] = 1
	}
	for col := 0; col < n; col++ {
		maxRow := col
		for row := col + 1; row < n; row++ {
			if math.Abs(aug[row][col]) > math.Abs(aug[maxRow][col]) {
				maxRow = row
			}
		}
		aug[col], aug[maxRow] = aug[maxRow], aug[col]
		pivot := aug[col][col]
		if math.Abs(pivot) < 1e-12 {
			return Matrix{}, errors.New("matriz singular")
		}
		for j := 0; j < 2*n; j++ {
			aug[col][j] /= pivot
		}
		for row := 0; row < n; row++ {
			if row == col {
				continue
			}
			factor := aug[row][col]
			for j := 0; j < 2*n; j++ {
				aug[row][j] -= factor * aug[col][j]
			}
		}
	}
	result := make([][]float64, n)
	for i := 0; i < n; i++ {
		result[i] = aug[i][n:]
	}
	return NewMatrix(result), nil
}
