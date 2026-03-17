package api

import (
	"encoding/json"
	"net/http"

	"github.com/seuuser/physicalc/internal/algebra"
	"github.com/seuuser/physicalc/internal/calculus"
	"github.com/seuuser/physicalc/internal/ode"
	"github.com/seuuser/physicalc/internal/physics"
)

// --- Helpers ---

// writeJSON escreve uma resposta JSON com o status code informado.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError escreve uma resposta de erro padronizada.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// HealthHandler retorna o status da API.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": "1.0.0"})
}

// =============================================================================
// CÁLCULO
// =============================================================================

// IntegrateRequest é o body esperado pelo endpoint de integração.
type IntegrateRequest struct {
	// A e B são os limites de integração
	A float64 `json:"a"`
	B float64 `json:"b"`
	// N é o número de subintervalos (mais alto = mais preciso)
	N int `json:"n"`
	// Function é o nome da função predefinida: "x2", "sin", "cos", "exp"
	Function string `json:"function"`
}

// IntegrateHandler calcula a integral numérica de uma função no intervalo [a, b].
//
// @Summary     Integração numérica
// @Tags        calculus
// @Accept      json
// @Produce     json
// @Param       body body IntegrateRequest true "Parâmetros"
// @Success     200 {object} map[string]any
// @Router      /calculus/integrate [post]
func IntegrateHandler(w http.ResponseWriter, r *http.Request) {
	var req IntegrateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if req.N <= 0 {
		req.N = 1000 // valor padrão
	}

	// Seleciona a função matemática pelo nome
	fn, err := calculus.GetFunction(req.Function)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	simpson := calculus.Simpson(fn, req.A, req.B, req.N)
	trapezoid := calculus.Trapezoid(fn, req.A, req.B, req.N)

	writeJSON(w, http.StatusOK, map[string]any{
		"function":  req.Function,
		"a":         req.A,
		"b":         req.B,
		"n":         req.N,
		"simpson":   simpson,
		"trapezoid": trapezoid,
	})
}

// DerivativeRequest é o body esperado pelo endpoint de derivada.
type DerivativeRequest struct {
	Function string  `json:"function"` // nome da função
	X        float64 `json:"x"`        // ponto onde calcular a derivada
	H        float64 `json:"h"`        // tamanho do passo (padrão: 1e-5)
}

// DerivativeHandler calcula a derivada numérica f'(x) usando diferenças centrais.
func DerivativeHandler(w http.ResponseWriter, r *http.Request) {
	var req DerivativeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if req.H == 0 {
		req.H = 1e-5
	}

	fn, err := calculus.GetFunction(req.Function)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := calculus.CentralDifference(fn, req.X, req.H)
	writeJSON(w, http.StatusOK, map[string]any{
		"function":   req.Function,
		"x":          req.X,
		"h":          req.H,
		"derivative": result,
	})
}

// LimitHandler calcula o limite numérico de uma função num ponto.
func LimitHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Function  string  `json:"function"`
		X         float64 `json:"x"`
		Direction string  `json:"direction"` // "left", "right", "both"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if req.Direction == "" {
		req.Direction = "both"
	}

	fn, err := calculus.GetFunction(req.Function)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	left, right := calculus.NumericalLimit(fn, req.X)
	writeJSON(w, http.StatusOK, map[string]any{
		"function":    req.Function,
		"x":           req.X,
		"limit_left":  left,
		"limit_right": right,
		"converges":   calculus.Converges(left, right, 1e-6),
	})
}

// =============================================================================
// ÁLGEBRA LINEAR
// =============================================================================

// VectorRequest é o body para operações com vetores.
type VectorRequest struct {
	Operation string    `json:"operation"` // "add", "dot", "cross", "norm", "angle"
	A         []float64 `json:"a"`
	B         []float64 `json:"b"`
}

// VectorHandler realiza operações vetoriais.
func VectorHandler(w http.ResponseWriter, r *http.Request) {
	var req VectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	va := algebra.NewVector(req.A)
	vb := algebra.NewVector(req.B)

	var result any
	var opErr error

	switch req.Operation {
	case "add":
		res, err := va.Add(vb)
		result, opErr = res, err
	case "sub":
		res, err := va.Sub(vb)
		result, opErr = res, err
	case "dot":
		res, err := va.Dot(vb)
		result, opErr = res, err
	case "cross":
		res, err := va.Cross(vb)
		result, opErr = res, err
	case "norm":
		result = va.Norm()
	case "angle":
		res, err := va.AngleDeg(vb)
		result, opErr = res, err
	case "normalize":
		res, err := va.Normalize()
		result, opErr = res, err
	default:
		writeError(w, http.StatusBadRequest, "operação inválida: "+req.Operation)
		return
	}

	if opErr != nil {
		writeError(w, http.StatusBadRequest, opErr.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"operation": req.Operation,
		"a":         req.A,
		"b":         req.B,
		"result":    result,
	})
}

// MatrixRequest é o body para operações com matrizes.
type MatrixRequest struct {
	Operation string      `json:"operation"` // "mul", "det", "transpose", "inverse"
	A         [][]float64 `json:"a"`
	B         [][]float64 `json:"b"`
}

// MatrixHandler realiza operações matriciais.
func MatrixHandler(w http.ResponseWriter, r *http.Request) {
	var req MatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	ma := algebra.NewMatrix(req.A)

	var result any
	var opErr error

	switch req.Operation {
	case "mul":
		mb := algebra.NewMatrix(req.B)
		res, err := ma.Multiply(mb)
		result, opErr = res, err
	case "det":
		res, err := ma.Determinant()
		result, opErr = res, err
	case "transpose":
		result = ma.Transpose()
	case "inverse":
		res, err := ma.Inverse()
		result, opErr = res, err
	case "trace":
		res, err := ma.Trace()
		result, opErr = res, err
	default:
		writeError(w, http.StatusBadRequest, "operação inválida: "+req.Operation)
		return
	}

	if opErr != nil {
		writeError(w, http.StatusBadRequest, opErr.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"operation": req.Operation,
		"result":    result,
	})
}

// =============================================================================
// FÍSICA
// =============================================================================

// KinematicsHandler calcula grandezas cinemáticas.
func KinematicsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		V0 float64 `json:"v0"` // velocidade inicial (m/s)
		A  float64 `json:"a"`  // aceleração (m/s²)
		T  float64 `json:"t"`  // tempo (s)
		X0 float64 `json:"x0"` // posição inicial (m)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	result := physics.Kinematics(req.X0, req.V0, req.A, req.T)
	writeJSON(w, http.StatusOK, result)
}

// DynamicsHandler calcula grandezas dinâmicas.
func DynamicsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Mass         float64 `json:"mass"`         // kg
		Force        float64 `json:"force"`        // N
		Displacement float64 `json:"displacement"` // m
		Velocity     float64 `json:"velocity"`     // m/s
		Height       float64 `json:"height"`       // m
		Gravity      float64 `json:"gravity"`      // m/s² (padrão: 9.81)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if req.Gravity == 0 {
		req.Gravity = 9.81
	}

	result := physics.Dynamics(req.Mass, req.Force, req.Displacement, req.Velocity, req.Height, req.Gravity)
	writeJSON(w, http.StatusOK, result)
}

// =============================================================================
// EQUAÇÕES DIFERENCIAIS
// =============================================================================

// ODESolveHandler resolve uma EDO de primeira ordem numericamente.
func ODESolveHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Method string  `json:"method"` // "euler" ou "rk4"
		Func   string  `json:"func"`   // "decay", "growth", "oscillation"
		Y0     float64 `json:"y0"`     // condição inicial
		T0     float64 `json:"t0"`     // tempo inicial
		TEnd   float64 `json:"t_end"`  // tempo final
		H      float64 `json:"h"`      // tamanho do passo
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if req.H == 0 {
		req.H = 0.01
	}
	if req.Method == "" {
		req.Method = "rk4"
	}

	fn, err := ode.GetODEFunction(req.Func)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var ts, ys []float64
	switch req.Method {
	case "euler":
		ts, ys = ode.Euler(fn, req.Y0, req.T0, req.TEnd, req.H)
	case "rk4":
		ts, ys = ode.RungeKutta4(fn, req.Y0, req.T0, req.TEnd, req.H)
	default:
		writeError(w, http.StatusBadRequest, "método inválido: use 'euler' ou 'rk4'")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"method": req.Method,
		"func":   req.Func,
		"y0":     req.Y0,
		"t":      ts,
		"y":      ys,
		"steps":  len(ts),
	})
}
