package physics_test

import (
	"testing"

	"github.com/seuuser/physicalc/internal/physics"
	"github.com/stretchr/testify/assert"
)

func TestKinematicsMRU(t *testing.T) {
	// MRU: v0=10m/s, a=0, t=5s → x = 50m
	result := physics.Kinematics(0, 10, 0, 5)
	assert.InDelta(t, 50.0, result.Position, 1e-10)
	assert.InDelta(t, 10.0, result.Velocity, 1e-10)
	assert.Equal(t, "MRU (Movimento Retilíneo Uniforme)", result.MovementType)
}

func TestKinematicsMRUA(t *testing.T) {
	// MRUA: v0=0, a=2m/s², t=3s → v=6m/s, x=9m
	result := physics.Kinematics(0, 0, 2, 3)
	assert.InDelta(t, 9.0, result.Position, 1e-10)
	assert.InDelta(t, 6.0, result.Velocity, 1e-10)
}

func TestKinematicsFreeFall(t *testing.T) {
	// Queda livre: v0=0, a=-9.81m/s², t=2s
	g := -9.81
	result := physics.Kinematics(0, 0, g, 2)
	// h = ½*g*t² = ½*9.81*4 = 19.62m
	assert.InDelta(t, 0.5*g*4, result.Position, 1e-6)
	assert.InDelta(t, g*2, result.Velocity, 1e-6)
}

func TestDynamicsNewton2(t *testing.T) {
	// F=10N, m=2kg → a=5m/s²
	result := physics.Dynamics(2, 10, 0, 0, 0, 9.81)
	assert.InDelta(t, 5.0, result.Acceleration, 1e-10)
	assert.InDelta(t, 2*9.81, result.Weight, 1e-6)
}

func TestDynamicsKineticEnergy(t *testing.T) {
	// m=2kg, v=3m/s → Ec = ½*2*9 = 9J
	result := physics.Dynamics(2, 0, 0, 3, 0, 9.81)
	assert.InDelta(t, 9.0, result.KineticEnergy, 1e-10)
}

func TestDynamicsPotentialEnergy(t *testing.T) {
	// m=2kg, h=5m, g=9.81 → Ep = 2*9.81*5 = 98.1J
	result := physics.Dynamics(2, 0, 0, 0, 5, 9.81)
	assert.InDelta(t, 98.1, result.PotentialEnergy, 1e-6)
}

func TestProjectile(t *testing.T) {
	// Lançamento a 45° com v0=20m/s, g=9.81
	result := physics.Projectile(20, 45, 9.81)
	// Alcance máximo ocorre a 45°
	assert.True(t, result.Range > 0)
	assert.True(t, result.MaxHeight > 0)
	assert.True(t, result.TimeOfFlight > 0)
}

func TestCircularMotion(t *testing.T) {
	// r=1m, f=1Hz, m=1kg
	result := physics.CircularMotion(1, 1, 1)
	assert.InDelta(t, 1.0, result.Period, 1e-10)
	assert.InDelta(t, 1.0, result.Frequency, 1e-10)
	// v = 2π ≈ 6.283
	assert.InDelta(t, 2*3.14159265358979, result.LinearVelocity, 1e-8)
}

func TestSHM(t *testing.T) {
	// k=100N/m, m=1kg, A=0.1m
	result := physics.SimpleHarmonicMotion(100, 1, 0.1)
	// ω = √(k/m) = 10 rad/s
	assert.InDelta(t, 10.0, result.AngularFrequency, 1e-6)
	// T = 2π/10 ≈ 0.628s
	assert.InDelta(t, 2*3.14159265358979/10, result.Period, 1e-8)
}
