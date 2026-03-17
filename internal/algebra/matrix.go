package algebra

import "errors"

// Matrix representa uma matriz m×n (linhas × colunas).
type Matrix struct {
	Rows, Cols int
	Data       [][]float64
}

// NewMatrix cria uma matriz m×n preenchida com zeros.
func NewMatrix(rows, cols int) *Matrix {
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return &Matrix{Rows: rows, Cols: cols, Data: data}
}

// NewMatrixFrom cria uma Matrix a partir de um slice 2D.
func NewMatrixFrom(data [][]float64) (*Matrix, error) {
	if len(data) == 0 {
		return nil, errors.New("matriz não pode ser vazia")
	}
	cols := len(data[0])
	for _, row := range data {
		if len(row) != cols {
			return nil, errors.New("todas as linhas devem ter o mesmo número de colunas")
		}
	}
	return &Matrix{Rows: len(data), Cols: cols, Data: data}, nil
}

// Get retorna o elemento na posição (i, j).
func (m *Matrix) Get(i, j int) float64 { return m.Data[i][j] }

// Set define o elemento na posição (i, j).
func (m *Matrix) Set(i, j int, v float64) { m.Data[i][j] = v }

// Multiply multiplica duas matrizes: C = A × B
// Requer que o número de colunas de A seja igual ao número de linhas de B.
func Multiply(a, b *Matrix) (*Matrix, error) {
	if a.Cols != b.Rows {
		return nil, errors.New("dimensões incompatíveis para multiplicação")
	}
	result := NewMatrix(a.Rows, b.Cols)
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			sum := 0.0
			for k := 0; k < a.Cols; k++ {
				sum += a.Data[i][k] * b.Data[k][j]
			}
			result.Data[i][j] = sum
		}
	}
	return result, nil
}

// Transpose retorna a transposta da matriz: A^T onde A^T[j][i] = A[i][j]
func Transpose(m *Matrix) *Matrix {
	result := NewMatrix(m.Cols, m.Rows)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.Data[j][i] = m.Data[i][j]
		}
	}
	return result
}

// Determinant calcula o determinante de uma matriz quadrada n×n
// usando expansão de Laplace (recursivo). Funciona bem até ~6×6.
func Determinant(m *Matrix) (float64, error) {
	if m.Rows != m.Cols {
		return 0, errors.New("determinante é definido apenas para matrizes quadradas")
	}
	return det(m.Data), nil
}

func det(data [][]float64) float64 {
	n := len(data)
	if n == 1 {
		return data[0][0]
	}
	if n == 2 {
		return data[0][0]*data[1][1] - data[0][1]*data[1][0]
	}
	result := 0.0
	for j := 0; j < n; j++ {
		result += data[0][j] * cofactor(data, 0, j)
	}
	return result
}

func cofactor(data [][]float64, row, col int) float64 {
	sign := 1.0
	if (row+col)%2 != 0 {
		sign = -1.0
	}
	return sign * det(minor(data, row, col))
}

func minor(data [][]float64, row, col int) [][]float64 {
	n := len(data)
	result := make([][]float64, 0, n-1)
	for i := range data {
		if i == row {
			continue
		}
		newRow := make([]float64, 0, n-1)
		for j := range data[i] {
			if j != col {
				newRow = append(newRow, data[i][j])
			}
		}
		result = append(result, newRow)
	}
	return result
}
