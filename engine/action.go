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

/*
import (
	"fmt"
	"strconv"
	//"github.com/ethereum/go-ethereum/log"

	"github.com/star-formation/tesseract/lib"
)

// Actions are authenticated requests to modify the game state.
// Most actions originate from users, where we consider them
// authenticated post user account signature verification.
// Actions can also originate from the game engine itself;
// such actions are always considered authenticated.
type Action interface {
	Execute() error
}

// The input to this function is raw bytes from other layers, e.g.
// the binary payload of a WebSocket message.  As such these bytes have
// not yet been validated, and may be malicious.
// TODO: refactor and document security assumptions and input validation
// in diff layers.
// TODO: for dev/test we use a simple JSON schema
func HandleAction(j map[string]interface{}) error {
	e, err := strconv.ParseUint(j["entity"].(string), 10, 64)
	if err != nil {
		return err
	}
	params := j["params"].(map[string]interface{})
	duration := params["duration"].(float64)

	var ar Action
	if j["action"] == "rotate" {
		torque := params["force"].(map[string]interface{})
		x := torque["x"].(float64)
		y := torque["y"].(float64)
		z := torque["z"].(float64)
		ar = &ActionRotate{e, &lib.V3{x, y, z}, duration}
	} else {
		f := params["force"].(float64)
		ar = &ActionEngineThrust{e, f, duration}
	}
	GE.actionChan <- ar

	return nil
}

type ActionRotate struct {
	entity   uint64
	t        *lib.V3 // torque
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
	entity   uint64
	thrust   float64
	duration float64
}

func (a *ActionEngineThrust) Execute() error {
	S.AddForceGen(a.entity, &ThrustForceGen{a.thrust, a.duration})
	return nil
}
*/
