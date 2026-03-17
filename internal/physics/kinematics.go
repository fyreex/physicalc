package physics

import "math"

// KinematicsResult agrupa todos os resultados de um cálculo cinemático.
type KinematicsResult struct {
	// Entradas
	X0 float64 `json:"x0_m"`  // posição inicial (m)
	V0 float64 `json:"v0_ms"` // velocidade inicial (m/s)
	A  float64 `json:"a_ms2"` // aceleração (m/s²)
	T  float64 `json:"t_s"`   // tempo (s)

	// Resultados (MRUA — Movimento Retilíneo Uniformemente Acelerado)
	Position     float64 `json:"position_m"`      // x = x0 + v0*t + ½*a*t²
	Velocity     float64 `json:"velocity_ms"`     // v = v0 + a*t
	Displacement float64 `json:"displacement_m"`  // Δx = x - x0
	Distance     float64 `json:"distance_m"`      // distância percorrida (sempre positiva)
	VelocitySq   float64 `json:"velocity_sq_ms2"` // v² = v0² + 2*a*Δx

	// Grandezas extras
	KineticEnergy float64 `json:"kinetic_energy_j"` // Ec = ½mv² (massa = 1kg por padrão)
	Momentum      float64 `json:"momentum_kgms"`    // p = mv (massa = 1kg)

	// Análise do movimento
	MovementType string `json:"movement_type"` // MRU, MRUA, etc.
}

// Kinematics calcula todas as grandezas do movimento uniformemente acelerado.
//
// Equações fundamentais usadas:
//   - x = x0 + v0*t + ½*a*t²
//   - v = v0 + a*t
//   - v² = v0² + 2*a*(x - x0)
func Kinematics(x0, v0, a, t float64) KinematicsResult {
	position := x0 + v0*t + 0.5*a*t*t
	velocity := v0 + a*t
	displacement := position - x0
	velocitySq := v0*v0 + 2*a*displacement

	// Distância percorrida considera reversão de direção
	distance := math.Abs(displacement)
	if a != 0 && v0 != 0 {
		// Verifica se há reversão: o objeto para em t_stop = -v0/a
		tStop := -v0 / a
		if tStop > 0 && tStop < t {
			// Distância até parar + distância após parar
			xStop := x0 + v0*tStop + 0.5*a*tStop*tStop
			distance = math.Abs(xStop-x0) + math.Abs(position-xStop)
		}
	}

	movType := classifyMovement(v0, a)

	return KinematicsResult{
		X0:            x0,
		V0:            v0,
		A:             a,
		T:             t,
		Position:      position,
		Velocity:      velocity,
		Displacement:  displacement,
		Distance:      distance,
		VelocitySq:    velocitySq,
		KineticEnergy: 0.5 * velocity * velocity, // massa = 1 kg
		Momentum:      velocity,                  // massa = 1 kg
		MovementType:  movType,
	}
}

// classifyMovement classifica o tipo de movimento com base em v0 e a.
func classifyMovement(v0, a float64) string {
	switch {
	case a == 0 && v0 == 0:
		return "Repouso"
	case a == 0:
		return "MRU (Movimento Retilíneo Uniforme)"
	case v0 == 0:
		return "Movimento a partir do repouso"
	case (v0 > 0 && a > 0) || (v0 < 0 && a < 0):
		return "MRUA Acelerado"
	default:
		return "MRUA Retardado"
	}
}

// ProjectileMotion calcula o movimento de um projétil lançado obliquamente.
type ProjectileResult struct {
	MaxHeight    float64 `json:"max_height_m"`     // altura máxima
	Range        float64 `json:"range_m"`          // alcance horizontal
	TimeOfFlight float64 `json:"time_of_flight_s"` // tempo de voo total
	AngleDeg     float64 `json:"angle_deg"`        // ângulo de lançamento
}

// Projectile calcula as grandezas de um lançamento oblíquo.
// v0 = velocidade inicial (m/s), angleDeg = ângulo em graus, g = gravidade (m/s²)
func Projectile(v0, angleDeg, g float64) ProjectileResult {
	angleRad := angleDeg * math.Pi / 180
	v0x := v0 * math.Cos(angleRad) // componente horizontal
	v0y := v0 * math.Sin(angleRad) // componente vertical

	// Tempo de voo: 0 = v0y*t - ½g*t² → t = 2*v0y/g
	tFlight := 2 * v0y / g

	// Altura máxima: h = v0y²/(2g)
	hMax := (v0y * v0y) / (2 * g)

	// Alcance: R = v0x * tFlight
	rangeX := v0x * tFlight

	return ProjectileResult{
		MaxHeight:    hMax,
		Range:        rangeX,
		TimeOfFlight: tFlight,
		AngleDeg:     angleDeg,
	}
}
