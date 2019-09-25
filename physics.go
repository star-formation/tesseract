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
	worldPoint.Sub(worldPoint, S.PC[e])
	return worldPoint.VectorProduct(worldPoint, force)
}
*/

// System interface
func (p *Physics) Init() error {
	return nil
}

func (p *Physics) Update(elapsed float64) error {
	// TODO: split into functions and loop over all ents _for each_ function
	for e, _ := range S.HotEnts {
		// update force generators
		linearForce, torque := new(V3), new(V3)
		expiredFGs := make(map[int]bool, 0)
		for i, fg := range *S.MC[e].FGs {
			lf, t := fg.UpdateForce(e, elapsed)
			//log.Debug("Physics.Update", "lf", lf, "t", t)
			if lf != nil {
				linearForce.Add(linearForce, lf)
			}
			if t != nil {
				torque.Add(torque, t)
			}
			if fg.IsExpired() {
				expiredFGs[i] = true
			}
		}

		// update linear acceleration from forces
		inverseMass := float64(1) / *(S.MassC[e])
		lastAcc := new(V3)
		lastAcc.AddScaledVector(linearForce, inverseMass)

		// update linear velocity
		S.MC[e].V.AddScaledVector(lastAcc, elapsed)

		// update angular acceleration from torques
		angularAcc := S.RC[e].IITW.Transform(torque)
		log.Debug("FUNK", "torque", torque, "angularAcc", angularAcc)

		// update angular velocity
		S.RC[e].R.AddScaledVector(angularAcc, elapsed)

		// apply damping (universal)
		S.MC[e].V.MulScalar(S.MC[e].V, math.Pow(linearDamping, elapsed))
		S.RC[e].R.MulScalar(S.RC[e].R, math.Pow(angularDamping, elapsed))

		// update linear position (V3.AddScaledVector)
		S.PC[e].AddScaledVector(S.MC[e].V, elapsed)

		// update angular position (Q.AddScaledVector)
		S.OC[e].AddScaledVector(S.RC[e].R, elapsed)

		// normalize orientation
		S.OC[e].Normalise()

		updateTransformMatrix(S.RC[e].T, S.PC[e], S.OC[e])

		// update inverse inertia tensor in world coordinates
		updateInertiaTensor(S.RC[e].IITW, S.RC[e].IITB, S.RC[e].T, S.OC[e])

		// Clear force/torque accumulator
		newFGs := make([]ForceGen, 0, len(*S.MC[e].FGs)-len(expiredFGs))
		for i, fg := range *S.MC[e].FGs {
			if !expiredFGs[i] {
				newFGs = append(newFGs, fg)
			}
		}
		*S.MC[e].FGs = newFGs

		// "p", S.PC[e], "v", S.VC[e]
		log.Debug("physics.Update", "o", S.OC[e], "r", S.RC[e].R)
	}

	return nil
}

// https://en.wikipedia.org/wiki/List_of_moments_of_inertia
// mass, width, height, depth
func InertiaTensorCuboid(m, w, h, d float64) *M3 {
	h2 := h * h
	w2 := w * w
	d2 := d * d
	x := (1.0 / 12.0) * m
	return &M3{x * (h2 + d2), 0, 0, 0, x * (w2 + d2), 0, 0, 0, x * (w2 + h2)}
}

// Update transform matrix (m) from position (p) and orientation (o)
func updateTransformMatrix(m *M4, p *V3, o *Q) {
	m[0] = 1 - 2*o.J*o.J - 2*o.K*o.K
	m[1] = 2*o.I*o.J - 2*o.R*o.K
	m[2] = 2*o.I*o.K + 2*o.R*o.J
	m[3] = p.X

	m[4] = 2*o.I*o.J + 2*o.R*o.K
	m[5] = 1 - 2*o.I*o.I - 2*o.K*o.K
	m[6] = 2*o.J*o.K - 2*o.R*o.I
	m[7] = p.Y

	m[8] = 2*o.I*o.K - 2*o.R*o.J
	m[9] = 2*o.J*o.K + 2*o.R*o.I
	m[10] = 1 - 2*o.I*o.I - 2*o.J*o.J
	m[11] = p.Z
}

// update the inverse inertia tensor in world space coordinates (tw) using
// the inverse inertia tensor in body space coordinates (tb),
// the transform matrix (tm) and the orientation (o)
func updateInertiaTensor(tw, tb *M3, tm *M4, o *Q) {
	t4 := tm[0]*tb[0] + tm[1]*tb[3] + tm[2]*tb[6]
	t9 := tm[0]*tb[1] + tm[1]*tb[4] + tm[2]*tb[7]
	t14 := tm[0]*tb[2] + tm[1]*tb[5] + tm[2]*tb[8]
	t28 := tm[4]*tb[0] + tm[5]*tb[3] + tm[6]*tb[6]
	t33 := tm[4]*tb[1] + tm[5]*tb[4] + tm[6]*tb[7]
	t38 := tm[4]*tb[2] + tm[5]*tb[5] + tm[6]*tb[8]
	t52 := tm[8]*tb[0] + tm[9]*tb[3] + tm[10]*tb[6]
	t57 := tm[8]*tb[1] + tm[9]*tb[4] + tm[10]*tb[7]
	t62 := tm[8]*tb[2] + tm[9]*tb[5] + tm[10]*tb[8]

	tw[0] = t4*tm[0] + t9*tm[1] + t14*tm[2]
	tw[1] = t4*tm[4] + t9*tm[5] + t14*tm[6]
	tw[2] = t4*tm[8] + t9*tm[9] + t14*tm[10]
	tw[3] = t28*tm[0] + t33*tm[1] + t38*tm[2]
	tw[4] = t28*tm[4] + t33*tm[5] + t38*tm[6]
	tw[5] = t28*tm[8] + t33*tm[9] + t38*tm[10]
	tw[6] = t52*tm[0] + t57*tm[1] + t62*tm[2]
	tw[7] = t52*tm[4] + t57*tm[5] + t62*tm[6]
	tw[8] = t52*tm[8] + t57*tm[9] + t62*tm[10]
}
