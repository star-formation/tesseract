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
	"math"
)

type FComp map[*RefFrame]map[Id]*float64 // Generic float64 component
type V3Comp map[*RefFrame]map[Id]*V3     // Generic 3x1 vector component
type M3Comp map[*RefFrame]map[Id]*M3     // Generic 3x3 matrix component
type QComp map[*RefFrame]map[Id]*Q       // Generic quaternion component

const (
	linearDamping  = float64(0.99)
	angularDamping = float64(0.99)
)

type RadiusComp map[*RefFrame]map[Id]*float64

// TODO ...
type HotEnts struct {
	Frames []*RefFrame
	In     map[*RefFrame]*[]Id
}

// Physics Engine
type Physics struct {
	ents map[Id]bool

	// Mass component holds float64 values
	mc FComp

	// Components for Position, Velocity, Acceleration, Rotation
	// and Force/Torque Accumulators hold 3x1 vectors
	pc, vc, ac, rc, fc, tc V3Comp

	// Orientation Component holds quaternions
	oc QComp

	// Inertia component holds 3x3 matrices
	ic M3Comp
}

// System interface
func (p *Physics) Init() error {
	return nil
}
func (p *Physics) Update(elapsed float64, f *RefFrame, hotEnts *[]Id) error {
	// TODO: split into functions and loop over all ents _for each_ function

	for _, e := range *hotEnts {
		// 1. update linear acceleration from forces
		inverseMass := float64(1) / *p.mc[f][e]
		lastAcc := p.ac[f][e]
		lastAcc.AddScaledVector(p.fc[f][e], inverseMass)

		// 2. update angular acceleration from torques
		angularAcc := p.ic[f][e].Transform(p.tc[f][e])

		// 3. update linear and angular velocity
		p.vc[f][e].AddScaledVector(lastAcc, elapsed)
		p.rc[f][e].AddScaledVector(angularAcc, elapsed)

		// 4. impose drag
		p.vc[f][e].MulScalar(p.vc[f][e], math.Pow(linearDamping, elapsed))
		p.rc[f][e].MulScalar(p.rc[f][e], math.Pow(angularDamping, elapsed))

		// 5. update linear position (V3.AddScaledVector)
		p.pc[f][e].AddScaledVector(p.vc[f][e], elapsed)

		// 6. update angular position (Q.AddScaledVector)
		p.oc[f][e].AddScaledVector(p.rc[f][e], elapsed)

		// 7. normalize orientation
		p.oc[f][e].Normalise()

		// TODO: 8. update derived data
	}

	return nil
}

func (p *Physics) RegisterEntity(id Id) {
	p.ents[id] = true
}
func (p *Physics) DeregisterEntity(id Id) {
	p.ents[id] = false
}

func magnitude(x, y, z float64) float64 {
	return math.Sqrt(x*x + y*y + z*z)
}
