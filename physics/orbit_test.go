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

import (
	"testing"
)

// Exact values in these tests often differ slightly from the
// calculations in [1] as the book rounds intermediate and final values to
// only a few decimals whereas we retain 64 bit precision in all calculations.

// TODO: add test data from Section 4.4

func TestStateVectorConv(t *testing.T) {
	// Example 4.3.
	r := &V3{-6045, -3490, 2500}
	v := &V3{-3.457, 6.618, 2.533}
	μ := 398600.0

	/*
		ex := &OE{
			58310,           // km^2/s (magnitude of specific angular momentum)
			DegToRad(153.2), // degrees (inclination)
			DegToRad(255.3), // degrees (longitude of the ascending node)
			0.1712,          // dimensionless (eccentricity)
			DegToRad(20.07), // degrees (argument of periapsis)
			DegToRad(28.45), // degrees (true anomaly)
			398600.0,
		}
	*/

	calc := StateVectorToOrbital(r, v, μ)
	calc.Debug()
	// TODO: add comparison within threshold
}

func TestOrbElemConv(t *testing.T) {
	// Example 4.7.
	o := &OE{80000, DegToRad(30), DegToRad(40), 1.4, DegToRad(60), DegToRad(30), 398600}
	r, v := o.OrbitalToStateVector()

	re := &V3{-4039.895923201738, 4814.560480182377, 3628.6247021718837}
	ve := &V3{-10.385987618194685, -4.771921637340853, 1.7438750000000005}

	if r.X != re.X || r.Y != re.Y || r.Z != re.Z {
		t.Errorf("pos: got: \n%v, expected: \n%v", r, re)
	}

	if v.X != ve.X || v.Y != ve.Y || v.Z != ve.Z {
		t.Errorf("vel: got: \n%v, expected: \n%v", v, ve)
	}

	o2 := StateVectorToOrbital(r, v, o.μ)
	o2.Debug()
	// TODO: check o2 vs o
}

func TestFromTimeElliptic(t *testing.T) {
	// This test uses intermediate values in Example 3.1 and 3.2.
	o := &OE{}
	o.h = 72472
	o.e = 0.37255
	o.μ = 398600.0

	// In example 3.1, step 3, the final value is 193.2
	exθ := DegToRad(193.1540909884592)
	θ := o.TrueAnomalyFromTime(0, 10800)
	if θ != exθ {
		t.Errorf("θ: got: \n%v, expected: \n%v", RadToDeg(θ), RadToDeg(exθ))
	}
}

func TestFromTimeParabolic(t *testing.T) {
	// This test uses values in Example 3.4.
	o := &OE{}
	o.h = 79720
	o.e = 1.0
	o.μ = 398600.0

	exθ := 144.75444965830107
	θ := o.TrueAnomalyFromTime(0, 6*3600)
	if RadToDeg(θ) != exθ {
		t.Errorf("θ: got: \n%v, expected: \n%v", RadToDeg(θ), exθ)
	}

	o.θ = θ

	exAlt := 86976.62246749947
	alt := o.Altitude()
	if alt != exAlt {
		t.Errorf("altitude: got: \n%v, expected: \n%v", alt, exAlt)
	}
}

func TestFromTimeHyperbolic(t *testing.T) {
	// This test uses values in Example 3.5.
	o := &OE{}
	o.h = 100170
	o.e = 2.7696
	o.μ = 398600.0

	// Example 3.5 step 5 yields 107.78 degrees
	exθ := 1.8811167388351486 //DegToRad(107.78)
	θ := o.TrueAnomalyFromTime(0, 4141.4+3*3600)
	if θ != exθ {
		t.Errorf("θ: got: \n%v, expected: \n%v", θ, exθ)
	}
	o.θ = θ

	exAlt := 163181.8946911754
	alt := o.Altitude()
	if alt != exAlt {
		t.Errorf("altitude: got: \n%v, expected: \n%v", alt, exAlt)
	}

}

func TestPointsApprox(t *testing.T) {
	
}

//
// Benchmarks
//
func BenchmarkOrbitConv(b *testing.B) {
	rnd, _ := NewRand(42)
	Rand = rnd

	μ := 398600.0
	posScale := 10000.0
	velScale := 10000.0

	r := &V3{}
	v := &V3{}
	oe := &OE{}

	randRV := func() {
		r.X = Rand.Float64() * posScale
		r.Y = Rand.Float64() * posScale
		r.Z = Rand.Float64() * posScale
		v.X = Rand.Float64() * velScale
		v.Y = Rand.Float64() * velScale
		v.Z = Rand.Float64() * velScale
	}

	randOE := func() {
		oe.h = 10000 + Rand.Float64()*100000
		oe.i = Rand.Float64() * DegToRad(180)
		oe.Ω = Rand.Float64() * DegToRad(360)
		oe.e = Rand.Float64() * 2.0
		oe.ω = Rand.Float64() * DegToRad(180)
		oe.θ = Rand.Float64() * DegToRad(360)
		oe.μ = Rand.Float64()*μ*1.5 + μ/3
	}

	b.ResetTimer()
	b.Run("SVToOE", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			randRV()
			b.StartTimer()
			oe = StateVectorToOrbital(r, v, μ)
		}
	})

	b.Run("OEToSV", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			randOE()
			b.StartTimer()
			r, v = oe.OrbitalToStateVector()
		}
	})

	b.Run("SVToOEToSV", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			randRV()
			b.StartTimer()
			oe = StateVectorToOrbital(r, v, μ)
			r, v = oe.OrbitalToStateVector()
			oe = StateVectorToOrbital(r, v, μ)
		}
	})
}
