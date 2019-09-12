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

/* Vector, Matrix and Quaternion Math.

A good introduction is chapters 2 and 9 in [1].

References:

[1] Millington, Ian. Game physics engine development (Second Edition). CRC Press, 2010.

[2] https://github.com/idmillington/cyclone-physics/blob/master/include/cyclone/core.h

*/

const (
	DBL_EPSILON = 2.2204460492503131E-16
)

type V3 struct {
	X, Y, Z float64
}

func (v *V3) Add(a, b *V3) *V3 {
	v.X = a.X + b.X
	v.Y = a.Y + b.Y
	v.Z = a.Z + b.Z
	return v
}

func (v *V3) Sub(a, b *V3) *V3 {
	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z
	return v
}

func (v *V3) MulScalar(a *V3, s float64) *V3 {
	v.X = a.X * s
	v.Y = a.Y * s
	v.Z = a.Z * s
	return v
}

func (v *V3) AddScaledVector(a *V3, s float64) *V3 {
	v.X += a.X * s
	v.Y += a.Y * s
	v.Z += a.Z * s
	return v
}

func (v *V3) ComponentProduct(a, b *V3) *V3 {
	v.X = a.X * b.X
	v.Y = a.Y * b.Y
	v.Z = a.Z * b.Z
	return v
}

func (v *V3) VectorProduct(a, b *V3) *V3 {
	v.X = a.Y*b.Z - a.Z*b.Y
	v.Y = a.Z*b.X - a.X*b.Z
	v.Z = a.X*b.Y - a.Y*b.X
	return v
}

func (v *V3) ScalarProduct(a *V3) float64 {
	return v.X*a.X + v.Y*a.Y + v.Z*a.Z
}

func (v *V3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *V3) SquareMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *V3) Invert() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

// 3x3 Matrix
type M3 [9]float64

func (m *M3) Transform(v *V3) *V3 {
	return &V3{
		v.X*m[0] + v.Y*m[1] + v.Z*m[2],
		v.X*m[3] + v.Y*m[4] + v.Z*m[5],
		v.X*m[6] + v.Y*m[7] + v.Z*m[8],
	}
}

func (m *M3) TransformTranspose(v *V3) *V3 {
	return &V3{
		v.X*m[0] + v.Y*m[3] + v.Z*m[6],
		v.X*m[1] + v.Y*m[4] + v.Z*m[7],
		v.X*m[2] + v.Y*m[5] + v.Z*m[8],
	}
}

func (m *M3) Mul(a *M3) *M3 {
	t1 := m[0]*a[0] + m[1]*a[3] + m[2]*a[6]
	t2 := m[0]*a[1] + m[1]*a[4] + m[2]*a[7]
	t3 := m[0]*a[2] + m[1]*a[5] + m[2]*a[8]
	m[0] = t1
	m[1] = t2
	m[2] = t3

	t1 = m[3]*a[0] + m[4]*a[3] + m[5]*a[6]
	t2 = m[3]*a[1] + m[4]*a[4] + m[5]*a[7]
	t3 = m[3]*a[2] + m[4]*a[5] + m[5]*a[8]
	m[3] = t1
	m[4] = t2
	m[5] = t3

	t1 = m[6]*a[0] + m[7]*a[3] + m[8]*a[6]
	t2 = m[6]*a[1] + m[7]*a[4] + m[8]*a[7]
	t3 = m[6]*a[2] + m[7]*a[5] + m[8]*a[8]
	m[6] = t1
	m[7] = t2
	m[8] = t3

	return m
}

func (m *M3) Inverse() {
	x4 := m[0] * m[4]
	x6 := m[0] * m[5]
	x8 := m[1] * m[3]
	x10 := m[2] * m[3]
	x12 := m[1] * m[6]
	x14 := m[2] * m[6]

	det := x4*m[8] - x6*m[7] - x8*m[8] + x10*m[7] + x12*m[5] - x14*m[4]

	// TODO: check whether to handle error or safe to ignore
	if det == 0 {
		panic("zero matrix determinant")
	}
	x17 := 1 / det

	m[0] = (m[4]*m[8] - m[5]*m[7]) * x17
	m[1] = -(m[1]*m[8] - m[2]*m[7]) * x17
	m[2] = (m[1]*m[5] - m[2]*m[4]) * x17
	m[3] = -(m[3]*m[8] - m[5]*m[6]) * x17
	m[4] = (m[0]*m[8] - x14) * x17
	m[5] = -(x6 - x10) * x17
	m[6] = (m[3]*m[7] - m[4]*m[6]) * x17
	m[7] = -(m[0]*m[7] - x12) * x17
	m[8] = (x4 - x8) * x17

}

// 3x4 Matrix
type M4 [12]float64

// Quaternion
type Q struct {
	R, I, J, K float64
}

func (q *Q) Normalise() {
	x := q.R*q.R + q.I*q.I + q.J*q.J + q.K*q.K

	// return no-rotation quaternion if zero length
	// TODO: check if this is correct usage of DBL_EPSILON
	// - scaled correctly for magnitude of x?
	// - use math.Nextafter instead?
	if x < DBL_EPSILON {
		q.R = 1
		return
	}

	x = 1.0 / math.Sqrt(x)
	q.R *= x
	q.I *= x
	q.J *= x
	q.K *= x
}

func (q *Q) Mul(a *Q) *Q {
	q.R = q.R*a.R - q.I*a.I - q.J*a.J - q.K*a.K
	q.I = q.R*a.I + q.I*a.R + q.J*a.K - q.K*a.J
	q.J = q.R*a.J + q.J*a.R + q.K*a.I - q.I*a.K
	q.K = q.R*a.K + q.K*a.R + q.I*a.J - q.J*a.I
	return q

}

func (q *Q) AddScaledVector(v *V3, s float64) *Q {
	q2 := Q{0, v.X * s, v.Y * s, v.Z * s}
	q2.Mul(q)
	q.R += q2.R * 0.5
	q.I += q2.I * 0.5
	q.J += q2.J * 0.5
	q.K += q2.K * 0.5
	return q
}
