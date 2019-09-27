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
	"github.com/ethereum/go-ethereum/log"
)

/* Force Generators
 */

type ForceGen interface {
	// Returns linear force and torque
	UpdateForce(e Id, duration float64) (*V3, *V3)

	IsExpired() bool
}

// TODO: apply same drag function on rotation
type DragForceGen struct {
	DragCoef1, DragCoef2 float64
}

func (d *DragForceGen) UpdateForce(e Id, duration float64) (*V3, *V3) {
	if S.MC[e].V.IsZero() {
		return nil, nil
	}
	vel := new(V3)
	*vel = *(S.MC[e].V)
	velMag := vel.Magnitude()
	drag := velMag*d.DragCoef1 + velMag*velMag*d.DragCoef2
	force := vel
	force.Normalise()
	force.MulScalar(force, -drag)
	log.Debug("DragForceGen.UpdateForce", "f", force)
	return force, nil
}

func (d *DragForceGen) IsExpired() bool {
	return false
}

// For e.g. center-of-mass-aligned engines
type ThrustForceGen struct {
	thrust   float64
	timeLeft float64
}

func (t *ThrustForceGen) UpdateForce(e Id, elapsed float64) (*V3, *V3) {
	log.Debug("ThrustForceGen.UpdateForce", "t", t.thrust)
	var f float64
	if t.timeLeft > elapsed {
		f = t.thrust * elapsed
		t.timeLeft -= elapsed
	} else {
		f = t.thrust * t.timeLeft
		t.timeLeft = 0
	}

	// align thrust force to orientation
	// rotate the forward axis with e's orientation, then scalar mult force
	fv := S.OC[e].ForwardVector()
	return fv.MulScalar(fv, f), nil
}

func (t *ThrustForceGen) IsExpired() bool {
	return t.timeLeft == 0
}

// For Ship turning
type TurnForceGen struct {
	torque   *V3
	timeLeft float64
}

func (t *TurnForceGen) UpdateForce(e Id, elapsed float64) (*V3, *V3) {
	//log.Debug("TurnForceGen.UpdateForce", "tq", *t.torque)
	if t.timeLeft > elapsed {
		t.timeLeft -= elapsed
		return nil, new(V3).MulScalar(t.torque, elapsed)
	}

	resTorque := new(V3).MulScalar(t.torque, t.timeLeft)
	t.timeLeft = 0
	return nil, resTorque
}

func (t *TurnForceGen) IsExpired() bool {
	return t.timeLeft == 0
}
