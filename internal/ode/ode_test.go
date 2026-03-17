package ode_test

import (
	"math"
	"testing"

	"github.com/fyreex/physicalc/internal/ode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEulerDecay verifica o decaimento exponencial: y' = -y, y(0) = 1
// Solução exata: y(t) = e^(-t)
func TestEulerDecay(t *testing.T) {
	f, err := ode.GetODEFunction("decay")
	require.NoError(t, err)

	ts, ys := ode.Euler(f, 1.0, 0, 1.0, 0.001)
	exact := math.Exp(-ts[len(ts)-1])

	// Euler é menos preciso, tolerância maior
	assert.InDelta(t, exact, ys[len(ys)-1], 1e-2)
}

// TestRK4Decay verifica RK4 com decaimento: muito mais preciso que Euler
func TestRK4Decay(t *testing.T) {
	f, err := ode.GetODEFunction("decay")
	require.NoError(t, err)

	ts, ys := ode.RungeKutta4(f, 1.0, 0, 1.0, 0.01)
	exact := math.Exp(-ts[len(ts)-1])

	// RK4 é muito preciso com h=0.01
	assert.InDelta(t, exact, ys[len(ys)-1], 1e-8)
}

// TestRK4Growth verifica crescimento exponencial: y' = y, y(0) = 1
// Solução exata: y(t) = e^t
func TestRK4Growth(t *testing.T) {
	f, err := ode.GetODEFunction("growth")
	require.NoError(t, err)

	ts, ys := ode.RungeKutta4(f, 1.0, 0, 1.0, 0.01)
	exact := math.Exp(ts[len(ts)-1])

	assert.InDelta(t, exact, ys[len(ys)-1], 1e-8)
}

// TestRK4VsEulerAccuracy compara precisão dos dois métodos
func TestRK4VsEulerAccuracy(t *testing.T) {
	f, err := ode.GetODEFunction("decay")
	require.NoError(t, err)

	_, ysEuler := ode.Euler(f, 1.0, 0, 1.0, 0.1)
	_, ysRK4 := ode.RungeKutta4(f, 1.0, 0, 1.0, 0.1)
	exact := math.Exp(-1.0)

	errEuler := math.Abs(ysEuler[len(ysEuler)-1] - exact)
	errRK4 := math.Abs(ysRK4[len(ysRK4)-1] - exact)

	// RK4 deve ser muito mais preciso que Euler
	assert.True(t, errRK4 < errEuler, "RK4 deve ser mais preciso que Euler")
}

// TestOutputLength verifica que os slices de saída têm o tamanho correto
func TestOutputLength(t *testing.T) {
	f, _ := ode.GetODEFunction("decay")
	ts, ys := ode.RungeKutta4(f, 1.0, 0, 1.0, 0.1)
	assert.Equal(t, len(ts), len(ys), "ts e ys devem ter o mesmo tamanho")
	assert.Equal(t, 11, len(ts)) // 0, 0.1, 0.2, ..., 1.0 = 11 pontos
}

// TestGetODEFunctionInvalid verifica que função desconhecida retorna erro
func TestGetODEFunctionInvalid(t *testing.T) {
	_, err := ode.GetODEFunction("naoexiste")
	assert.Error(t, err)
}
