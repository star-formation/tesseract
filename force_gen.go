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
	UpdateForce(e Id, rf *RefFrame, duration float64) (*V3, *V3)

	IsExpired() bool
}

type DragForceGen struct {
	DragCoef1, DragCoef2 float64
}

func (d *DragForceGen) UpdateForce(e Id, rf *RefFrame, duration float64) (*V3, *V3) {
	if S.VC[rf][e].IsZero() {
		return nil, nil
	}
	vel := new(V3)
	*vel = *(S.VC[rf][e])
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

type ThrustForceGen struct {
	thrust *V3
	// TODO: think about temporary and persistent thrust
	//expiry float64
}

func (t *ThrustForceGen) UpdateForce(e Id, rf *RefFrame, duration float64) (*V3, *V3) {
	return t.thrust, nil
}

func (t *ThrustForceGen) IsExpired() bool {
	return false
}
