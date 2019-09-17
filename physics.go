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

// Physics Engine
type Physics struct {
}

/*
func (p *Physics) TorqueAtBodyPoint(e Id, rf *RefFrame, force, bodyPoint *V3) *V3 {
	worldPoint := bodyToWorldPoint(bodyPoint)
	return TorqueAtPoint(e, rf, force, worldPoint)
}


// TODO: in body.cpp , the force is not actually split into force on center of
//       of mass and torque, but adds torque _on_top_of_ the force-at-point!
// See https://www.gamedev.net/forums/topic/664930-force-and-torque/
// TODO: clarify this assumption on e.g. position/source of force
func (p *Physics) TorqueAtPoint(e Id, rf *RefFrame, force, worldPoint *V3) *V3 {
	point := new(V3)
	*point = *worldPoint
	worldPoint.Sub(worldPoint, S.PC[rf][e])
	return worldPoint.VectorProduct(worldPoint, force)
}
*/

// System interface
func (p *Physics) Init() error {
	return nil
}

func (p *Physics) Update(elapsed float64, rf *RefFrame, hotEnts *[]Id) error {
	// TODO: split into functions and loop over all ents _for each_ function

	for _, e := range *hotEnts {
		// update force generators
		linearForce, torque := new(V3), new(V3)
		for _, fg := range S.FC[rf][e] {
			lf, t := fg.UpdateForce(e, rf, elapsed)
			if lf != nil {
				linearForce.Add(linearForce, lf)
			}
			if t != nil {
				torque.Add(torque, t)
			}
		}

		// update linear acceleration from forces
		inverseMass := float64(1) / *(S.MC[rf][e])
		lastAcc := new(V3)
		lastAcc.AddScaledVector(linearForce, inverseMass)

		// update linear velocity
		S.VC[rf][e].AddScaledVector(lastAcc, elapsed)

		// update angular acceleration from torques
		angularAcc := S.IC[rf][e].Transform(torque)

		// update angular velocity
		S.RC[rf][e].AddScaledVector(angularAcc, elapsed)

		// apply damping (universal)
		S.VC[rf][e].MulScalar(S.VC[rf][e], math.Pow(linearDamping, elapsed))
		S.RC[rf][e].MulScalar(S.RC[rf][e], math.Pow(angularDamping, elapsed))

		// update linear position (V3.AddScaledVector)
		S.PC[rf][e].AddScaledVector(S.VC[rf][e], elapsed)

		// update angular position (Q.AddScaledVector)
		S.OC[rf][e].AddScaledVector(S.RC[rf][e], elapsed)

		// normalize orientation
		S.OC[rf][e].Normalise()

		// TODO: update derived data

		log.Debug("physics.Update", "p", S.PC[rf][e], "v", S.VC[rf][e], "o", S.OC[rf][e], "r", S.RC[rf][e])
	}

	return nil
}
