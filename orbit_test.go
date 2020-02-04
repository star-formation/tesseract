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
	"testing"
)

type orbTest struct {
	Name     string
	pos, vel *V3
	gm       float64
	o        *OE
}

func TestOrbitalElementConv(t *testing.T) {
	case2 := orbTest{
		"earth1",
		&V3{REarth + 600.0*1000, 0.0, 50.0},
		&V3{0.0, 6.5 * 1000, 0.0},
		MUEarth,
		&OE{0.26035023023005477, 5.536635637306466e+06, 7.165289638066952e-06, -1.5707963267948966, -1.5707963267948966, -2049.9813051575525, 3.141592653589793}, // TODO
	}

	/*
		case3 := orbTest{
			"ISS",
			&V3{},
			&V3{},
			MUEarth,
			&OE{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		}
	*/

	cases := []orbTest{
		case2,
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			o := CartToOE(tc.pos, tc.vel, tc.gm)
			o.Debug()
			if o.E != tc.o.E || o.S != tc.o.S || o.I != tc.o.I || o.L != tc.o.L || o.A != tc.o.A || o.T != tc.o.T {

				t.Errorf("Cartesian to Orbital Elements mismatch, got:\n%v\nwant:\n%v", o, tc.o)
			}
			pos, vel := OEToCart(o, tc.gm)
			if math.Abs(pos.X)-math.Abs(tc.pos.X) > tolerance ||
				math.Abs(pos.Y)-math.Abs(tc.pos.Y) > tolerance ||
				math.Abs(pos.Z)-math.Abs(tc.pos.Z) > 0.001 { // TODO: lower
				t.Errorf("Orbital Elements to Cartesian mismatch, got:\n%v\nwant:\n%v", pos, tc.pos)
			}
			if math.Abs(vel.X)-math.Abs(tc.vel.X) > 0.000000000001 ||
				math.Abs(vel.Y)-math.Abs(tc.vel.Y) > 0.00000000001 ||
				math.Abs(vel.Z)-math.Abs(tc.vel.Z) > tolerance {
				t.Errorf("Orbital Elements to Cartesian mismatch, got:\n%v\nwant:\n%v", vel, tc.vel)
			}
		})
	}
}

func TestOrbitalElementConvISS(t *testing.T) {
	ISSOE := &OE{
		E: 0.0005040,          // dimensionless
		S: REarth + 422*1000,  // m
		I: DegToRad(51.6420),  // radians
		L: DegToRad(286.4848), // radians
		A: DegToRad(230.6842), // radians
		T: DegToRad(129.3862), // radians
	}

	pos, vel := OEToCart(ISSOE, MUEarth)
	fmt.Printf("pos: %v\nvel: %v\nvelMag: %v\n", pos, vel, vel.Magnitude())

	o := CartToOE(pos, vel, MUEarth)
	o.Debug()
}
