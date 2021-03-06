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

import "math"

type Planet struct {
	Entity Id
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
	μ := GravitationalConstant * p.Mass
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
