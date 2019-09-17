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

const (
	g0 = 9.80665

	linearDamping  = float64(1.0)
	angularDamping = float64(1.0)
)

type FComp map[*RefFrame]map[Id]*float64 // Generic float64 component
type V3Comp map[*RefFrame]map[Id]*V3     // Generic 3x1 vector component
type M3Comp map[*RefFrame]map[Id]*M3     // Generic 3x3 matrix component
type QComp map[*RefFrame]map[Id]*Q       // Generic quaternion component

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

	// Position, Velocity and Rotation components holds 3x1 vectors
	PC, VC, RC V3Comp

	// Orientation Component holds quaternions
	OC QComp

	// Inertia component holds 3x3 matrices
	IC M3Comp

	FC map[*RefFrame]map[Id][]ForceGen
}

func NewPhysics() *Physics {
	p := new(Physics)

	p.Ents = make(map[Id]bool, 10)

	p.MC = make(FComp, 1)

	p.PC = make(V3Comp, 1)
	p.VC = make(V3Comp, 1)
	p.RC = make(V3Comp, 1)

	p.OC = make(QComp, 1)

	p.IC = make(M3Comp, 1)

	p.FC = make(map[*RefFrame]map[Id][]ForceGen, 1)

	return p
}

func (p *Physics) NewFrame() *RefFrame {
	rf := new(RefFrame)

	p.MC[rf] = make(map[Id]*float64, 10)

	p.PC[rf] = make(map[Id]*V3, 10)
	p.VC[rf] = make(map[Id]*V3, 10)
	p.RC[rf] = make(map[Id]*V3, 10)

	p.OC[rf] = make(map[Id]*Q, 10)

	p.IC[rf] = make(map[Id]*M3, 10)

	p.FC[rf] = make(map[Id][]ForceGen, 10)

	return rf
}

// TODO: auto-increment entity ID and decouple its components
func (p *Physics) NewEntity(e Id, rf *RefFrame) {
	p.VC[rf][e] = &V3{} // velocity

	p.IC[rf][e] = &M3{} // inertia tensor

	p.OC[rf][e] = &Q{}  // orientation
	p.RC[rf][e] = &V3{} // rotation

	p.FC[rf][e] = []ForceGen{} // Force Generators
}

func (p *Physics) AddForceGen(e Id, rf *RefFrame, fg ForceGen) {
	p.FC[rf][e] = append(p.FC[rf][e], fg)
}

/*
func (p *Physics) TorqueAtBodyPoint(e Id, rf *RefFrame, force, bodyPoint *V3) *V3 {
	worldPoint := bodyToWorldPoint(bodyPoint)
	return TorqueAtPoint(e, rf, force, worldPoint)
}
*/

// TODO: in body.cpp , the force is not actually split into force on center of
//       of mass and torque, but adds torque _on_top_of_ the force-at-point!
// See https://www.gamedev.net/forums/topic/664930-force-and-torque/
// TODO: clarify this assumption on e.g. position/source of force
func (p *Physics) TorqueAtPoint(e Id, rf *RefFrame, force, worldPoint *V3) *V3 {
	point := new(V3)
	*point = *worldPoint
	worldPoint.Sub(worldPoint, p.PC[rf][e])
	return worldPoint.VectorProduct(worldPoint, force)
}

// System interface
func (p *Physics) Init() error {
	return nil
}

func (p *Physics) Update(elapsed float64, rf *RefFrame, hotEnts *[]Id) error {
	// TODO: split into functions and loop over all ents _for each_ function

	for _, e := range *hotEnts {
		// update force generators
		linearForce, torque := new(V3), new(V3)
		for _, fg := range p.FC[rf][e] {
			lf, t := fg.UpdateForce(e, rf, p, elapsed)
			if lf != nil {
				linearForce.Add(linearForce, lf)
			}
			if t != nil {
				torque.Add(torque, t)
			}
		}

		// update linear acceleration from forces
		inverseMass := float64(1) / *(p.MC[rf][e])
		lastAcc := new(V3)
		lastAcc.AddScaledVector(linearForce, inverseMass)

		// update linear velocity
		p.VC[rf][e].AddScaledVector(lastAcc, elapsed)

		// update angular acceleration from torques
		angularAcc := p.IC[rf][e].Transform(torque)

		// update angular velocity
		p.RC[rf][e].AddScaledVector(angularAcc, elapsed)

		// apply damping (universal)
		p.VC[rf][e].MulScalar(p.VC[rf][e], math.Pow(linearDamping, elapsed))
		p.RC[rf][e].MulScalar(p.RC[rf][e], math.Pow(angularDamping, elapsed))

		// update linear position (V3.AddScaledVector)
		p.PC[rf][e].AddScaledVector(p.VC[rf][e], elapsed)

		// update angular position (Q.AddScaledVector)
		p.OC[rf][e].AddScaledVector(p.RC[rf][e], elapsed)

		// normalize orientation
		p.OC[rf][e].Normalise()

		// TODO: update derived data

		log.Debug("physics.Update", "p", p.PC[rf][e], "v", p.VC[rf][e], "o", p.OC[rf][e], "r", p.RC[rf][e])
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
