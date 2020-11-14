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
package engine

import (
	"github.com/star-formation/tesseract/lib"
	"github.com/star-formation/tesseract/physics"
	
	"github.com/ethereum/go-ethereum/log"
)

//
// References:
//
// [1] Millington, Ian. Game physics engine development (Second Edition). CRC Press, 2010.
// [2] https://github.com/idmillington/cyclone-physics/blob/master/include/cyclone/core.h
//

// The Physics system simulates classical mechanics.
type Physics struct{}

//
// System interface
//
func (p *Physics) Init() error {
	return nil
}

func (p *Physics) Update(worldTime, elapsed float64, rf *physics.RefFrame) error {
	//log.Debug("Physics ====")
	for e, _ := range S.HotEnts[rf] {
		// TODO: after initial orbit debug, add len == 0 check
		if S.ForceGens[e] != nil && len(S.ForceGens[e]) > 0 {
			updateClassicalMechanics(worldTime, elapsed, rf, e)
		}
	}

	// TODO: update ref frames
	return nil
}

func (p *Physics) IsHotPostUpdate(e uint64) bool {
	return S.ForceGens[e] != nil && len(S.ForceGens[e]) > 0
}

//
// Internal functions
//
func updateClassicalMechanics(worldTime, elapsed float64, rf *physics.RefFrame, e uint64) {
	var pos, vel *lib.V3
	if S.Orb[e] != nil {
		log.Debug("updateClassicalMechanics", "oe", S.Orb[e].Fmt())
		pos, vel = S.Orb[e].OrbitalToStateVector()
		log.Debug("updateClassicalMechanics", "pos", pos.Fmt(), "vel", vel.Fmt())
	} else {
		pos = S.Pos[e]
		vel = S.Vel[e]
	}

	log.Debug("updateClassicalMechanics", "e", e, "len(fgs)", len(S.ForceGens[e]))

	// update force generators
	linearForce, torque := new(lib.V3), new(lib.V3)
	expiredFGs := make(map[int]bool, 0)
	for i, fg := range S.ForceGens[e] {
		lf, t := fg.UpdateForce(e, elapsed)
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

	// TODO: skip updates if resulting linearForce and/or torque is zero.
	log.Debug("updateClassicalMechanics", "lf", linearForce, "tq", torque)

	if !linearForce.IsZero() {
		// update linear acceleration from forces
		inverseMass := float64(1) / *(S.Mass[e])
		acc := new(lib.V3)
		acc.AddScaledVector(linearForce, inverseMass)
		// update linear velocity
		vel.AddScaledVector(acc, elapsed)
	}

	if !torque.IsZero() {
		// update angular acceleration from torques
		angularAcc := S.Rot[e].IITW.Transform(torque)
		// update angular velocity
		S.Rot[e].R.AddScaledVector(angularAcc, elapsed)
	}

	// update linear position
	pos.AddScaledVector(vel, elapsed)

	// update angular position (orientation)
	S.Ori[e].AddScaledVector(S.Rot[e].R, elapsed)
	// normalize orientation
	S.Ori[e].Normalise()

	// apply damping (universal)
	//vel.MulScalar(vel, math.Pow(linearDamping, elapsed))
	//S.Rot[e].R.MulScalar(S.Rot[e].R, math.Pow(angularDamping, elapsed))

	updateTransformMatrix(S.Rot[e].T, pos, S.Ori[e])

	// update inverse inertia tensor in world coordinates
	updateInertiaTensor(S.Rot[e].IITW, S.Rot[e].IITB, S.Rot[e].T)

	if S.Orb[e] != nil {
		S.Orb[e] = physics.StateVectorToOrbital(pos, vel, S.Orb[e].Mu())
		log.Debug("updateClassicalMechanics", "oe2", S.Orb[e].Fmt())
	} else {
		S.Pos[e] = pos
		S.Vel[e] = vel
	}

	// Clear force/torque accumulator
	newFGCount := len(S.ForceGens[e]) - len(expiredFGs)
	if newFGCount == 0 {
		S.ForceGens[e] = []ForceGen{}
	} else {
		newFGs := make([]ForceGen, 0, len(S.ForceGens[e])-len(expiredFGs))
		for i, fg := range S.ForceGens[e] {
			if !expiredFGs[i] {
				newFGs = append(newFGs, fg)
			}
		}
		S.ForceGens[e] = newFGs
	}

	//log.Debug("physics.Update", "p", S.PC[e], "v", S.MC[e].V, "o", S.ORIC[e], "r", S.RC[e].R)
}

// TODO: check inertia tensor functions and cuboid tensor for Y axis

// https://en.wikipedia.org/wiki/List_of_moments_of_inertia
//                       mass, width, height, depth
func InertiaTensorCuboid(m, w, h, d float64) *lib.M3 {
	h2 := h * h
	w2 := w * w
	d2 := d * d
	x := (1.0 / 12.0) * m
	return &lib.M3{x * (h2 + d2), 0, 0, 0, x * (w2 + d2), 0, 0, 0, x * (w2 + h2)}
}

// Update transform matrix (m) from position (p) and orientation (o)
func updateTransformMatrix(m *lib.M4, p *lib.V3, o *lib.Q) {
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
// the inverse inertia tensor in body space coordinates (tb) and
// the transform matrix (tm)
func updateInertiaTensor(tw, tb *lib.M3, tm *lib.M4) {
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

/*
func (p *Physics) TorqueAtBodyPoint(e uint64, rf *physics.RefFrame, force, bodyPoint *lib.V3) *lib.V3 {
	worldPoint := bodyToWorldPoint(bodyPoint)
	return TorqueAtPoint(e, rf, force, worldPoint)
}

// TODO: in body.cpp , the force is not actually split into force on center of
//       of mass and torque, but adds torque _on_top_of_ the force-at-point!
// See https://www.gamedev.net/forums/topic/664930-force-and-torque/
// TODO: clarify this assumption on e.g. position/source of force
func (p *Physics) TorqueAtPoint(e uint64, rf *physics.RefFrame, force, worldPoint *lib.V3) *lib.V3 {
	point := new(lib.V3)
	*point = *worldPoint
	worldPoint.Sub(worldPoint, S.PC[e])
	return worldPoint.VectorProduct(worldPoint, force)
}
*/
