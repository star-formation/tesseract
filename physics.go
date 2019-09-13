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

	"github.com/ethereum/go-ethereum/log"
)

type FComp map[*RefFrame]map[Id]*float64 // Generic float64 component
type V3Comp map[*RefFrame]map[Id]*V3     // Generic 3x1 vector component
type M3Comp map[*RefFrame]map[Id]*M3     // Generic 3x3 matrix component
type QComp map[*RefFrame]map[Id]*Q       // Generic quaternion component

const (
	linearDamping  = float64(0.9999)
	angularDamping = float64(0.9999)
)

type RadiusComp map[*RefFrame]map[Id]*float64

// TODO ...
type HotEnts struct {
	Frames []*RefFrame
	In     map[*RefFrame]*[]Id
}

// Physics Engine
type Physics struct {
	Ents map[Id]bool

	// Mass component holds float64 values
	MC FComp

	// Components for Position, Velocity, Acceleration, Rotation
	// and Force/Torque Accumulators hold 3x1 vectors
	PC, VC, AC, RC, FC, TC V3Comp

	// Orientation Component holds quaternions
	OC QComp

	// Inertia component holds 3x3 matrices
	IC M3Comp
}

func NewPhysics() *Physics {
	p := new(Physics)

	p.Ents = make(map[Id]bool, 10)

	p.MC = make(FComp, 1)

	p.PC = make(V3Comp, 1)
	p.VC = make(V3Comp, 1)
	p.AC = make(V3Comp, 1)
	p.RC = make(V3Comp, 1)
	p.FC = make(V3Comp, 1)
	p.TC = make(V3Comp, 1)

	p.OC = make(QComp, 1)

	p.IC = make(M3Comp, 1)

	return p
}

func (p *Physics) NewFrame() *RefFrame {
	rf := new(RefFrame)

	p.MC[rf] = make(map[Id]*float64, 10)

	p.PC[rf] = make(map[Id]*V3, 10)
	p.VC[rf] = make(map[Id]*V3, 10)
	p.AC[rf] = make(map[Id]*V3, 10)
	p.RC[rf] = make(map[Id]*V3, 10)
	p.FC[rf] = make(map[Id]*V3, 10)
	p.TC[rf] = make(map[Id]*V3, 10)

	p.OC[rf] = make(map[Id]*Q, 10)

	p.IC[rf] = make(map[Id]*M3, 10)

	return rf
}

// TODO: auto-increment entity ID and decouple its components
func (p *Physics) NewEntity(e Id, rf *RefFrame) {
	p.AC[rf][e] = &V3{} // acceleration
	p.FC[rf][e] = &V3{} // force accumulator

	p.IC[rf][e] = &M3{} // inertia tensor
	p.TC[rf][e] = &V3{} // torque accumulator

	p.OC[rf][e] = &Q{}  // orientation
	p.RC[rf][e] = &V3{} // rotation
}

// System interface
func (p *Physics) Init() error {
	return nil
}

func (p *Physics) Update(elapsed float64, f *RefFrame, hotEnts *[]Id) error {
	// TODO: split into functions and loop over all ents _for each_ function

	for _, e := range *hotEnts {
		// 1. update linear acceleration from forces
		inverseMass := float64(1) / *(p.MC[f][e])
		lastAcc := p.AC[f][e]
		lastAcc.AddScaledVector(p.FC[f][e], inverseMass)

		// 2. update angular acceleration from torques
		angularAcc := p.IC[f][e].Transform(p.TC[f][e])

		// 3. update linear and angular velocity
		p.VC[f][e].AddScaledVector(lastAcc, elapsed)
		p.RC[f][e].AddScaledVector(angularAcc, elapsed)

		// 4. impose drag
		p.VC[f][e].MulScalar(p.VC[f][e], math.Pow(linearDamping, elapsed))
		p.RC[f][e].MulScalar(p.RC[f][e], math.Pow(angularDamping, elapsed))

		// 5. update linear position (V3.AddScaledVector)
		p.PC[f][e].AddScaledVector(p.VC[f][e], elapsed)

		// 6. update angular position (Q.AddScaledVector)
		p.OC[f][e].AddScaledVector(p.RC[f][e], elapsed)

		// 7. normalize orientation
		p.OC[f][e].Normalise()

		// TODO: 8. update derived data

		log.Debug("physics.Update", "e", e, "pos", p.PC[f][e], "vel", p.VC[f][e])
	}

	return nil
}

func (p *Physics) RegisterEntity(id Id) {
	p.Ents[id] = true
}
func (p *Physics) DeregisterEntity(id Id) {
	p.Ents[id] = false
}

func magnitude(x, y, z float64) float64 {
	return math.Sqrt(x*x + y*y + z*z)
}
