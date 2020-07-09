/*  Copyright 2019 The tesseract Authors

    This file is part of tesseract.

    tesseract is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    tesseract is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package tesseract

import (
	"fmt"
	"math"
)

// Types and functions for Orbital Mechanics.
//
// NOTE: Unless otherwise noted, referenced equations, algorithms, tables,
//       chapters and examples are from reference [1].
//
// References:
//
// [1] Curtis, H.D., 2013. Orbital mechanics for engineering students.
// [2] https://en.wikipedia.org/wiki/Orbital_elements
// [3] https://en.wikipedia.org/wiki/Specific_relative_angular_momentum
//

// OE holds orbital elements [2] and related variables to uniquely represent
// a specific orbit in the simplified two-body model.
// The orbiter is assumed to have neglible mass compared to the primary (host)
// body and the primary is stationary in the applicable reference frame.
//
// Elements/Fields:
//
// h: Specific angular momentum (See [3] and chapter 2.4 in [1])
// i: Inclination
// Ω: Longitude of the ascending node
// e: Eccentricity
// ω: Argument of Periapsis
// θ: True Anomaly
// μ: Standard gravitational parameter of the primary
type OE struct {
	h, i, Ω, e, ω, θ, μ float64
}

func (o *OE) Clone() interface{} {
	return &OE{o.h, o.i, o.Ω, o.e, o.ω, o.θ, o.μ}
}

func (o *OE) Debug() {
	fmt.Printf("h: %f i: %f Ω: %f e: %f ω: %f θ: %f μ: %f\n",
		o.h, RadToDeg(o.i), RadToDeg(o.Ω), o.e, RadToDeg(o.ω), RadToDeg(o.θ), o.μ)
}

func (o *OE) Fmt() string {
	return fmt.Sprintf("h: %.14g i: %.12f Ω: %.12f e: %.12f ω: %.12f θ: %.12f μ: %.4g",
		o.h, RadToDeg(o.i), RadToDeg(o.Ω), o.e, RadToDeg(o.ω), RadToDeg(o.θ), o.μ)
}

// SemimajorAxis returns the semimajor axis of the orbit.
func (o *OE) SemimajorAxis() float64 {
	h, e, μ := o.h, o.e, o.μ
	// Eqn 2.71 and 3.47
	return ((h * h) / μ) * (1 / (1 - e*e))
}

// SemiminorAxis returns the semiminor axis of the orbit.
// TODO: generalize to all orbit types, see table 3.1
func (o *OE) SemiminorAxis() float64 {
	e := o.e
	// Eqn 2.76
	return o.SemimajorAxis() * math.Sqrt(1-e*e)
}

// Periapsis returns the nearest point of the orbit.
func (o *OE) Periapsis() float64 {
	h, e, μ := o.h, o.e, o.μ
	// Eqn 2.50
	return ((h * h) / μ) * (1 / (1 + e))
}

// Apoapsis returns the farthest point of the orbit.
// Positive infinity is returned if the orbit is parabolic or hyperbolic.
func (o *OE) Apoapsis() float64 {
	h, e, μ := o.h, o.e, o.μ
	if e >= 1.0 {
		return math.Inf(1)
	} else {
		// Eqn 2.70
		return ((h * h) / μ) * (1 / (1 - e))
	}
}

// Altitude returns the distance between the orbiter and the primary.
// Altitude works for parabolic and hyperbolic orbits, but returns
// positive infinity if 1 + e*math.Cos(θ) <= 0.
func (o *OE) Altitude() float64 {
	h, e, θ, μ := o.h, o.e, o.θ, o.μ
	// Eqn 2.45
	d := 1 + e*math.Cos(θ)
	if d <= 0 {
		return math.Inf(1)
	} else {
		return ((h * h) / μ) * (1 / d)
	}
}

// Speed returns the orbital speed of the orbiter relative to the primary.
// https://en.wikipedia.org/wiki/Vis-viva_equation
func (o *OE) Speed() float64 {
	μ := o.μ
	r := o.Altitude()
	a := o.SemimajorAxis()
	return math.Sqrt(μ * ((2 / r) - (1 / a)))
}

// Period returns the orbital period.
// Positive infinity is returned if the orbit is parabolic or hyperbolic.
// https://en.wikipedia.org/wiki/Orbital_period#Small_body_orbiting_a_central_body
func (o *OE) Period() float64 {
	e, μ := o.e, o.μ
	if e < 1 { // elliptical and circular
		return twoPi * math.Sqrt(math.Pow(o.SemimajorAxis(), 3)/μ)
	}

	// parabolic and hyperbolic
	return math.Inf(1) // positive infinity
}

// TrueAnomalyFromTime returns the orbit's true anomaly in radians at time t1
// using t0 as the time of the set true anomaly.
// TrueAnomalyFromTime panics if t1 is not greater than t0.
func (o *OE) TrueAnomalyFromTime(t0, t1 float64) float64 {
	if t1 <= t0 {
		panic("t1 must be greater than t0")
	}

	h, e, μ := o.h, o.e, o.μ
	var θ float64

	switch {
	case e == 0: // circular
		// Chapter 3.3
		θ = (twoPi / o.Period()) * t1

	case e < 1: // elliptical
		// Eqn 3.8
		Me := (twoPi / o.Period()) * t1
		// Algorithm 3.1
		E := eccentricAnomaly(e, Me)
		// Eqn 3.13a
		x0 := math.Sqrt((1.0+e)/(1.0-e)) * math.Tan(E/2.0)
		θ = 2 * math.Atan(x0)

	case e == 1: // parabolic
		// Eqn 3.31
		Mp := (μ * μ * t1) / (h * h * h)
		// Eqn 3.32
		x0 := (3*Mp + math.Sqrt(math.Pow(3*Mp, 2)+1))
		x1 := math.Pow(x0, 1.0/3.0) - math.Pow(x0, -1.0/3.0)
		θ = 2 * math.Atan(x1)

	case e > 1: // hyperbolic
		// Eqn 3.34
		x0 := (μ * μ) / (h * h * h)
		x1 := math.Pow(e*e-1, 3.0/2.0)
		Mh := x0 * x1 * t1
		// Algorithm 3.2
		F := hyperbolicEccentricAnomaly(e, Mh)
		// Eqn 3.44b
		x2 := math.Sqrt((e+1)/(e-1)) * math.Tanh(F/2)
		θ = 2 * math.Atan(x2)
	}

	return NormalizeAngle(θ)
}

// TimeFromTrueAnomaly returns the time for a given true anomaly.
func (o *OE) TimeFromTrueAnomaly(θ float64) float64 {
	// TODO: normalize / mod
	if θ < 0 || θ > twoPi {
		panic("ta must be in 0-2π radians")
	}
	h, e, μ := o.h, o.e, o.μ
	var t float64

	switch {
	case e == 0: // circular
		// Chapter 3.3
		t = ((h * h * h) / (μ * μ)) * θ
	case e < 1: // elliptical
		// Eqn 3.13b
		x0 := math.Sqrt((1 - e) / (1 + e))
		x1 := math.Tan(θ) / 2
		E := 2 * math.Atan(x0*x1)
		// Eqn 3.14
		Me := E - e*math.Sin(E)
		// Eqn 3.15
		t = (Me / twoPi) * o.Period()
	case e == 1: // parabolic
		// Eqn 3.30 (substitution for Mp)
		x0 := math.Tan(θ / 2)
		Mp := (1/2)*x0 + (1/6)*math.Pow(x0, 3)
		// Eqn 3.31 (substitution for t)
		t = (Mp * (h * h * h)) / (μ * μ)
	case e > 1: // hyperbolic
		// Eqn 3.33
		x0 := e * math.Sqrt(e*e-1) * math.Sin(θ)
		x1 := 1 + e*math.Cos(θ)
		x2 := math.Sqrt(e + 1)
		x3 := math.Sqrt(e-1) * math.Tan(θ/2)
		x4 := math.Log((x2 + x3) / (x2 - x3))
		Mh := (x0 / x1) - x4
		// Eqn 3.34 (substitution for t)
		t = Mh / (((μ * μ) / (h * h * h)) * math.Pow(e*e-1, 3.0/2.0))
	}

	return t
}

func eccentricAnomaly(e, Me float64) float64 {
	// Algorithm 3.1
	Ei := Me + e/2
	if Me > math.Pi {
		Ei = Me - e/2
	}

	var ratio float64

NewApproximation:
	ratio = (Ei - e*math.Sin(Ei) - Me) / (1 - e*math.Cos(Ei))
	if math.Abs(ratio) > eccentricAnomalyTolerance {
		Ei -= ratio
		goto NewApproximation
	}

	return Ei
}

func hyperbolicEccentricAnomaly(e, Mh float64) float64 {
	// Algorithm 3.2
	Fi := Mh

	var ratio float64

NewApproximation:
	ratio = (e*math.Sinh(Fi) - Fi - Mh) / (e*math.Cosh(Fi) - 1)
	if math.Abs(ratio) > hyperbolicEccentricAnomalyTolerance {
		Fi -= ratio
		goto NewApproximation
	}

	return Fi
}

// StateVectorToOrbital returns the orbital elements converted from the orbital
// state vector and standard gravitational parameter of the primary.
// See Algorithm 4.2 and https://en.wikipedia.org/wiki/Orbital_state_vectors
func StateVectorToOrbital(r, v *V3, μ float64) *OE {
	dist := r.Magnitude()
	speed := v.Magnitude()
	radialVel := r.ScalarProduct(v) / dist

	h := new(V3).VectorProduct(r, v)
	hMag := h.Magnitude()

	i := math.Acos(h.Z / hMag)

	nodeLine := new(V3).VectorProduct(KHat, h)
	nodeLineMag := nodeLine.Magnitude()

	Ω := 0.0
	if nodeLine.X != 0.0 {
		Ω = math.Acos(nodeLine.X / nodeLineMag)
	}
	if nodeLine.Y < 0 {
		Ω = 2*math.Pi - Ω
	}

	// x0, .., xn hold intermediate calculations for eq 4.10
	x0 := (speed * speed) - (μ / dist)
	x1 := new(V3).MulScalar(r, x0)
	x2 := new(V3).MulScalar(v, dist*radialVel)
	x3 := new(V3).Sub(x1, x2)
	eVec := new(V3).MulScalar(x3, 1/μ)
	e := eVec.Magnitude()

	ω := 0.0
	if nodeLineMag != 0.0 {
		ω = math.Acos(nodeLine.ScalarProduct(eVec) / (nodeLineMag * e))
	}
	if eVec.Z < 0 {
		ω = 2*math.Pi - ω
	}

	θ := math.Acos(eVec.ScalarProduct(r) / (e * dist))
	if radialVel < 0 {
		θ = 2*math.Pi - θ
	}

	return &OE{hMag, i, Ω, e, ω, θ, μ}
}

// OrbitalToStateVector returns the orbital state vector of the orbit.
// Perifocal coordinates are used in an intermediate step.
// See Algorithm 4.5 and
// https://en.wikipedia.org/wiki/Perifocal_coordinate_system
func (o *OE) OrbitalToStateVector() (*V3, *V3) {
	h, i, Ω, e, ω, θ, μ := o.h, o.i, o.Ω, o.e, o.ω, o.θ, o.μ
	// x0, ..., xn hold intermediate calculations
	x0 := ((h * h) / μ)
	x0 *= (1 / (1 + e*math.Cos(θ)))
	periPos := &V3{
		x0 * math.Cos(θ),
		x0 * math.Sin(θ),
		0.0}

	x1 := μ / h

	periVel := &V3{
		x1 * (-math.Sin(θ)),
		x1 * (e + math.Cos(θ)),
		0.0,
	}

	// Eqn 4.49
	perifocalToHostcentric := &M3{
		-math.Sin(Ω)*math.Cos(i)*math.Sin(ω) + math.Cos(Ω)*math.Cos(ω),
		-math.Sin(Ω)*math.Cos(i)*math.Cos(ω) - math.Cos(Ω)*math.Sin(ω),
		math.Sin(Ω) * math.Sin(i),

		math.Cos(Ω)*math.Cos(i)*math.Sin(ω) + math.Sin(Ω)*math.Cos(ω),
		math.Cos(Ω)*math.Cos(i)*math.Cos(ω) - math.Sin(Ω)*math.Sin(ω),
		-math.Cos(Ω) * math.Sin(i),

		math.Sin(i) * math.Sin(ω),
		math.Sin(i) * math.Cos(ω),
		math.Cos(i),
	}

	pos := perifocalToHostcentric.Transform(periPos)
	vel := perifocalToHostcentric.Transform(periVel)

	return pos, vel
}
