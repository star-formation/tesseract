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
package physics

import "math"

const (
	//
	// Physics and Astrophysics Constants (real world)
	//
	speedOfLight          = 299792458.0 // m/s
	GravitationalConstant = 6.674e-11
	stefanBoltzmann       = 5.670373e-8

	aum  = 149597870700.0    // meters per AU
	AULY = 63241.07708426628 // AU per light year
	AUPC = 2.06265e5         // AU per parsec
	lypc = 3.26156           // light years per parsec

	milkyWayDiscHeight = 2000.0 // ly
	radiusLocalBubble  = 150.0  // ly

	// https://en.wikipedia.org/wiki/Solar_mass
	// https://en.wikipedia.org/wiki/Solar_luminosity
	// https://en.wikipedia.org/wiki/Solar_radius
	solarMass   = 1.98855e30 // kg
	solarLum    = 3.828e26   // W
	solarRadius = 6.957e8    // m

	// https://en.wikipedia.org/wiki/Earth_mass
	earthMass             = 5.9722e24 // kg
	earthRadius           = 6.3781e6  // km
	earthMu               = 3.986004418e14
	earthSeaLevelPressure = 101325 // pascals
	g0                    = 9.80665
)

type Planet struct {
	Entity uint64
	Mass   float64
	Radius float64

	AxialTilt      float64
	RotationPeriod float64

	SurfaceGravity float64
	Atmosphere     *Atmosphere
}

// DefaultOrbit returns a circular, prograde orbit 100km above a planet's
// surface or above its atmosphere (if it has one).
func (p *Planet) DefaultOrbit() *OE {
	e, i, Ω, ω, θ := 0.0, 0.0, 0.0, 0.0, 0.0
	μ := GravitationalConstant * p.Mass * earthMass
	r := 100000.0
	if p.Atmosphere != nil {
		r += p.Atmosphere.Height
	}
	h := math.Sqrt((r + r*e*math.Cos(θ)) * μ)
	return &OE{h: h, μ: μ, e: e, i: i, Ω: Ω, ω: ω, θ: θ}
}

// https://en.wikipedia.org/wiki/Gravity_of_Earth#Altitude
// p.surface_gravity has been pre-calculated by world building scripts
func (p *Planet) GravityAtAltitude(alt float64) float64 {
	// TODO: for debugging...
	if alt < -p.Radius {
		panic("kraken")
	}

	// TODO: handle negative altitude (caves, canyons, etc)
	if alt < 0 {
		return p.SurfaceGravity
	}

	x := p.Radius / (p.Radius + alt)
	return p.SurfaceGravity * (x * x)
}

// Where we lock X,Y,Z coords does not matter, as we do not yet have
// surface features.  Simply select a orientation derived from the
// planet's orientation 3D vector.
func (p *Planet) GeodeticToCartesian(lat, lon float64) float64 {
	return 0
}
