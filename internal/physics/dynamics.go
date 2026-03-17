package physics

import "math"

// DynamicsResult agrupa resultados de um cálculo dinâmico.
type DynamicsResult struct {
	// Grandezas básicas
	Force        float64 `json:"force_n"`          // F = m*a (N)
	Acceleration float64 `json:"acceleration_ms2"` // a = F/m (m/s²)
	Weight       float64 `json:"weight_n"`         // P = m*g (N)

	// Energia e trabalho
	Work            float64 `json:"work_j"`             // W = F*d (J)
	KineticEnergy   float64 `json:"kinetic_energy_j"`   // Ec = ½mv² (J)
	PotentialEnergy float64 `json:"potential_energy_j"` // Ep = mgh (J)
	TotalEnergy     float64 `json:"total_energy_j"`     // E = Ec + Ep (J)
	Power           float64 `json:"power_w"`            // P = W/t (W) — t=1s
	Momentum        float64 `json:"momentum_kgms"`      // p = mv (kg·m/s)
}

// Dynamics calcula grandezas dinâmicas para um objeto.
//
// Leis de Newton aplicadas:
//   - 1ª Lei: objeto em repouso ou MRU sem força resultante
//   - 2ª Lei: F = m*a
//   - 3ª Lei: implícita nas interações
func Dynamics(mass, force, displacement, velocity, height, gravity float64) DynamicsResult {
	acceleration := 0.0
	if mass != 0 {
		acceleration = force / mass // 2ª Lei de Newton
	}

	work := force * displacement                      // W = F*d (força paralela ao deslocamento)
	kineticEnergy := 0.5 * mass * velocity * velocity // Ec = ½mv²
	potentialEnergy := mass * gravity * height        // Ep = mgh
	totalEnergy := kineticEnergy + potentialEnergy
	momentum := mass * velocity // p = mv
	weight := mass * gravity    // P = mg

	return DynamicsResult{
		Force:           force,
		Acceleration:    acceleration,
		Weight:          weight,
		Work:            work,
		KineticEnergy:   kineticEnergy,
		PotentialEnergy: potentialEnergy,
		TotalEnergy:     totalEnergy,
		Power:           work, // P = W/t, t = 1s por padrão
		Momentum:        momentum,
	}
}

// CircularMotion calcula grandezas de movimento circular uniforme (MCU).
type CircularResult struct {
	Period           float64 `json:"period_s"`              // T = 2π/ω
	Frequency        float64 `json:"frequency_hz"`          // f = 1/T
	AngularVelocity  float64 `json:"angular_velocity_rads"` // ω = 2π*f
	LinearVelocity   float64 `json:"linear_velocity_ms"`    // v = ω*r
	CentripetalAccel float64 `json:"centripetal_accel_ms2"` // ac = v²/r = ω²r
	CentripetalForce float64 `json:"centripetal_force_n"`   // Fc = m*ac
}

// CircularMotion calcula grandezas do MCU.
// radius em metros, frequencyHz em Hz, mass em kg.
func CircularMotion(radius, frequencyHz, mass float64) CircularResult {
	omega := 2 * math.Pi * frequencyHz          // velocidade angular
	period := 1 / frequencyHz                   // período
	linearVel := omega * radius                 // velocidade linear
	centripetalAccel := omega * omega * radius  // aceleração centrípeta
	centripetalForce := mass * centripetalAccel // força centrípeta

	return CircularResult{
		Period:           period,
		Frequency:        frequencyHz,
		AngularVelocity:  omega,
		LinearVelocity:   linearVel,
		CentripetalAccel: centripetalAccel,
		CentripetalForce: centripetalForce,
	}
}

// SimpleHarmonicMotion calcula grandezas do MHS (massa em mola).
type SHMResult struct {
	AngularFrequency float64 `json:"angular_frequency_rads"` // ω = √(k/m)
	Period           float64 `json:"period_s"`               // T = 2π/ω
	Frequency        float64 `json:"frequency_hz"`           // f = 1/T
	MaxVelocity      float64 `json:"max_velocity_ms"`        // vmax = ω*A
	MaxAcceleration  float64 `json:"max_acceleration_ms2"`   // amax = ω²*A
}

// SimpleHarmonicMotion calcula grandezas do MHS (oscilador harmônico).
// k = constante da mola (N/m), mass em kg, amplitude em metros.
func SimpleHarmonicMotion(k, mass, amplitude float64) SHMResult {
	omega := math.Sqrt(k / mass)
	period := 2 * math.Pi / omega
	return SHMResult{
		AngularFrequency: omega,
		Period:           period,
		Frequency:        1 / period,
		MaxVelocity:      omega * amplitude,
		MaxAcceleration:  omega * omega * amplitude,
	}
}
