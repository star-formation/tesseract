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
	"fmt"
	//"github.com/ethereum/go-ethereum/log"
)

// Actions are authenticated requests to modify the game state.
// Most actions originate from users, where we consider them
// authenticated post user account signature verification.
// Actions can also originate from the game engine itself;
// such actions are always considered authenticated.
type Action interface {
	Execute() error
}

type ActionRotate struct {
	entity   Id
	t        *V3 // torque
	duration float64
}

func (a *ActionRotate) Execute() error {
	max := S.ShipClass[a.entity].CMGTorqueCap()
	if a.t.X > max.X || a.t.Y > max.Y || a.t.Z > max.Z {
		return fmt.Errorf("torque %v %v %v larger than CMG cap %v %v %v", a.t.X, a.t.Y, a.t.Z, max.X, max.Y, max.Z)
	}

	S.AddForceGen(a.entity, &TurnForceGen{a.t, a.duration})
	return nil
}

type ActionEngineThrust struct {
	entity   Id
	thrust   float64
	duration float64
}

func (a *ActionEngineThrust) Execute() error {
	S.AddForceGen(a.entity, &ThrustForceGen{a.thrust, a.duration})
	return nil
}
