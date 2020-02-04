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

/*
 See https://en.wikipedia.org/wiki/Orbital_elements

 To reduce code verbosity in physics code, we use this naming:
 E = Eccentricity
 S = Semimajor Axis
 I = Inclination
 L = Longitude of the ascending node
 A = Argument of periapsis
 T = True anomaly

*/
type OE struct {
	E, S, I, L, A, T float64

	// TODO: verify
	EA float64
}

func (o *OE) Debug() {
	fmt.Printf("Orbital Params:\n")
	fmt.Printf("o: %v\n", o)
	fmt.Printf("a: %.3f km\n", o.S/1000)
	fmt.Printf("e: %.6f\n", o.E)
	fmt.Printf("i: %.1f deg\n", RadToDeg(o.I))
	fmt.Printf("Ω: %.1f deg\n", RadToDeg(o.L))
	fmt.Printf("ω: %.1f deg\n", RadToDeg(o.A))
	fmt.Printf("f: %.1f deg \n", RadToDeg(o.T))
	//log.Info("OE: ", "E", o.E, "S", o.S, "I", o.I, "L", o.L, "A", o.A, "T", o.T)
}

func OEToCart(o *OE, μ float64) (*V3, *V3) {
	a, e, i, Ω, ω, EA := o.S, o.E, o.I, o.L, o.A, o.EA

	ν := 2 * math.Atan(math.Sqrt((1+e)/(1-e))*math.Tan(EA/2))
	r := a * (1 - e*math.Cos(EA))
	h := math.Sqrt(μ * a * (1 - e*e))

	x := r * (math.Cos(Ω)*math.Cos(ω+ν) - math.Sin(Ω)*math.Sin(ω+ν)*math.Cos(i))
	y := r * (math.Sin(Ω)*math.Cos(ω+ν) + math.Cos(Ω)*math.Sin(ω+ν)*math.Cos(i))
	z := r * (math.Sin(i) * math.Sin(ω+ν))

	p := a * (1 - e*e)

	vx := (x*h*e/(r*p))*math.Sin(ν) - (h/r)*(math.Cos(Ω)*math.Sin(ω+ν)+math.Sin(Ω)*math.Cos(ω+ν)*math.Cos(i))
	vy := (y*h*e/(r*p))*math.Sin(ν) - (h/r)*(math.Sin(Ω)*math.Sin(ω+ν)-math.Cos(Ω)*math.Cos(ω+ν)*math.Cos(i))
	vz := (z*h*e/(r*p))*math.Sin(ν) - (h/r)*(math.Cos(ω+ν)*math.Sin(i))

	return &V3{x, y, z}, &V3{vx, vy, vz}
}

func CartToOE(pos, vel *V3, μ float64) *OE {
	hBar := new(V3).VectorProduct(pos, vel)
	h := hBar.Magnitude()

	r := pos.Magnitude()
	v := vel.Magnitude()

	E := 0.5*(v*v) - μ/r

	a := -μ / (2 * E)

	e := math.Sqrt(1 - (h*h)/(a*μ))

	i := math.Acos(hBar.Z / h)

	lan := math.Atan2(hBar.X, -hBar.Y)

	x0 := pos.Z / math.Sin(i)
	x1 := pos.X*math.Cos(lan) + pos.Y*math.Sin(lan)
	lat := math.Atan2(x0, x1)

	p := a * (1 - e*e)
	ν := math.Atan2(math.Sqrt(p/μ)*pos.ScalarProduct(vel), p-r)

	ap := lat - ν

	EA := 2 * math.Atan(math.Sqrt((1-e)/(1+e))*math.Tan(ν/2))

	n := math.Sqrt(μ / (a * a * a))
	// TODO: correct epoch/time calc
	t := 0.0
	T := t - (1/n)*(EA-e*math.Sin(EA))

	return &OE{S: a, E: e, I: i, A: ap, L: lan, T: T, EA: EA}
}

/*
// TODO: simplify/merge if statements and math
func CartToOE(pos, vel *V3, gm float64) *OE {
	// https://space.stackexchange.com/questions/1904/how-to-programmatically-calculate-orbital-elements-using-position-velocity-vecto
	// https://github.com/RazerM/orbital/blob/master/orbital/utilities.py#L252
	h := new(V3).VectorProduct(pos, vel)
	nHat := new(V3).VectorProduct(&KHat, h)

	// TODO: support parabolic orbits
	ev := eccentricity(pos, vel, gm)
	e := ev.Magnitude()
	if e == 1.0 {
		panic("orbital eccentricity == 1.0 (parabolic orbit)")
	}

	energy := specificOrbitalEnergy(pos, vel, gm)

	s := -gm / (2.0 * energy)
	i := math.Acos(h.Z / h.Magnitude())

	lim := 1e-15
	l, ap, f := 0.0, 0.0, 0.0
	if math.Abs(i) < lim {
		if math.Abs(e) >= lim {
			ap = math.Acos(ev.X / ev.Magnitude())
		}
	} else {
		l = math.Acos(nHat.X / nHat.Magnitude())
		if nHat.Y < 0 {
			l = 2*math.Pi - l
		}

		ap = math.Acos(nHat.ScalarProduct(ev) / (nHat.Magnitude() * e))
	}

	if math.Abs(e) < lim {
		if math.Abs(i) < lim {
			f = math.Acos(pos.X / pos.Magnitude())
			if vel.X > 0 {
				f = 2*math.Pi - f
			}
		} else {
			x0 := nHat.ScalarProduct(pos)
			x1 := nHat.Magnitude() * pos.Magnitude()
			f = math.Acos(x0 / x1)
			if nHat.ScalarProduct(vel) > 0 {
				f = 2*math.Pi - f
			}
		}
	} else {
		if ev.Z < 0 {
			ap = 2*math.Pi - ap
		}
		x0 := ev.ScalarProduct(pos)
		x1 := e * pos.Magnitude()
		f = math.Acos(x0 / x1)
		if pos.ScalarProduct(vel) < 0 {
			f = 2*math.Pi - f
		}
	}

	return &OE{E: e, S: s, I: i, L: l, A: ap, T: f}
}

func specificOrbitalEnergy(pos, vel *V3, gm float64) float64 {
	posMag := pos.Magnitude()
	velMag := vel.Magnitude()
	energy := (velMag * velMag) / 2.0
	energy -= gm / posMag
	return energy
}

func eccentricity(pos, vel *V3, gm float64) *V3 {
	// TODO: compare Magnitude vs Normalise
	posMag := pos.Magnitude()
	velMag := vel.Magnitude()
	// x0, x1, ... hold intermediate values (for easy debugging)
	x0 := velMag * velMag
	x0 -= gm / posMag

	x1 := new(V3).MulScalar(pos, x0)

	x2 := pos.ScalarProduct(vel)

	x3 := new(V3).MulScalar(vel, x2)

	x4 := new(V3).Sub(x1, x3)

	return new(V3).MulScalar(x4, 1/gm)
}
*/
