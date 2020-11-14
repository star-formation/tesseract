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
	"math"
	"testing"
)

type bsTest struct {
	Name string
	// A holds X,Y,Z 3D vector position and radius of sphere one, sphere two
	// and the third sphere representing the bounding volume encompassing
	// sphere one and two
	A [12]float64
}

func TestTwoBoundingSpheres(t *testing.T) {
	cases := []bsTest{
		bsTest{"one axis non-overlapping", [12]float64{
			// X,Y,Z position and radius of first sphere
			1.0, 1.0, 1.0, 1.0,
			// X,Y,Z position and radius of second sphere
			2.0, 1.0, 1.0, 1.0,
			// X,Y,Z position and radius of bounding volume
			// exactly covering sphere one and two
			1.5, 1.0, 1.0, 1.5}},
		bsTest{"two axis non-overlapping", [12]float64{
			1.0, 1.0, 1.0, 1.0,
			2.0, 2.0, 1.0, 1.0,
			1.5, 1.5, 1.0, math.Sqrt(1*1+1*1+0*0)/2 + 1}},
		bsTest{"three axis non-overlapping", [12]float64{
			1.0, 1.0, 1.0, 1.0,
			2.0, 2.0, 2.0, 1.0,
			1.5, 1.5, 1.5, math.Sqrt(1*1+1*1+1*1)/2 + 1}},
		bsTest{"three axis non-overlapping, larger dist", [12]float64{
			1.0, 1.0, 1.0, 1.0,
			3.0, 4.0, 5.0, 1.0,
			2.0, 2.5, 3.0, math.Sqrt(2*2+3*3+4*4)/2 + 1.0}},
		bsTest{"three axis non-overlapping, larger vol", [12]float64{
			1.0, 1.0, 1.0, 2.0,
			2.0, 2.0, 2.0, 2.0,
			1.5, 1.5, 1.5, math.Sqrt(1*1+1*1+1*1)/2 + 2.0}},
		bsTest{"one axis touching", [12]float64{
			1.0, 1.0, 1.0, 1.0,
			1.5, 1.0, 1.0, 1.0,
			1.25, 1.0, 1.0, 1.25}},
		bsTest{"one axis overlapping", [12]float64{
			1.0, 1.0, 1.0, 1.0,
			1.4, 1.0, 1.0, 1.0,
			1.2, 1.0, 1.0, 1.2}},
		bsTest{"two axis overlapping, larger vol", [12]float64{
			1.0, 1.0, 1.0, 40.0,
			5.0, 5.0, 1.0, 40.0,
			3.0, 3.0, 1.0, math.Sqrt(2*2+2*2+0*0) + 40}},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			s1 := &BoundingSphere{&V3{tc.A[0], tc.A[1], tc.A[2]}, tc.A[3]}
			s2 := &BoundingSphere{&V3{tc.A[4], tc.A[5], tc.A[6]}, tc.A[7]}
			bv := s1.NewBoundingVolume(s2)
			s3 := bv.(*BoundingSphere)
			if math.Abs(s3.P.X-tc.A[8]) > tolerance ||
				math.Abs(s3.P.Y-tc.A[9]) > tolerance ||
				math.Abs(s3.P.Z-tc.A[10]) > tolerance {
				t.Errorf("position mismatch got %v want %v", *s3.P, V3{tc.A[8], tc.A[9], tc.A[10]})
			}
			if math.Abs(s3.R-tc.A[11]) > tolerance {
				t.Errorf("radius mismatch got %v want %v", s3.R, tc.A[11])
			}
		})
	}
}

/*
func TestBVHTree(t *testing.T) {
	s1 := &BoundingSphere{&V3{1.0, 1.0, 1.0}, 1.0}
	s2 := &BoundingSphere{&V3{2.0, 1.0, 1.0}, 1.0}

	root := &BVHNode{}
	root.Insert(1, s1)
	root.Insert(2, s2)
}
*/
