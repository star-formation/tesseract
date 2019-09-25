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
	"encoding/json"
)

// Dev/Test global in-memory game state.  Used for dev/tests to simplify
// rapid iteration to data structures and game design/logic.

// This is especially important as the state will go through many iterations
// before it's clear to best encode optimize it for merkle trees/proofs in
// the outer blockchain layer.
var S *State

type Mobile struct {
	V   *V3         // velocity
	FGs *[]ForceGen // force/torque generators
}

// If an object can rotate it has an inertia tensor
type Rotational struct {
	R    *V3 // Rotation X,Y,Z 3D vector
	IITB *M3 // Inverse Inertia Tensor Body space/coordinates
	IITW *M3 // Inverse Inertia Tensor World space/coordinates
	T    *M4 // Transform 3x4 matrix transforming body space/coordinates to world
}

// TODO: JSON tags and encoding for entire state.  Easiest during dev/test.
//       Replaced by compact, fast binary encoding for prod.
type State struct {
	// Game Engine State
	MB *MessageBus
	AB *MessageBus

	Ents    map[Id]bool
	HotEnts map[Id]bool

	//RefFrames []*RefFrame
	//RFC map[Id]*RefFrame
	//HotRefFrames []*RefFrame

	//
	// Physics Components
	//
	// Mass component holds float64 values
	MassC map[Id]*float64

	// Position, Velocity components holds 3x1 vectors
	PC map[Id]*V3

	// Orientation Component holds quaternions
	OC map[Id]*Q

	// Holds velocity and force generators for movable entities
	MC map[Id]Mobile

	// Holds rotational data for entities that can rotate
	RC map[Id]Rotational

	// TODO: for now all shapes are perfect bounding spheres
	// Sphere Radius Component
	SRC map[Id]*float64

	// Cargo Objects Component (entities in a cargo bay)
	COC map[Id][]Id

	// Volume Component
	VOC map[Id]*float64

	// Total Aerodynamic Lift and Drag Coefficients
	// These change depending on attached modules and player skills
	// TODO: split into subsonic, supersonic, hypersonic
	AeroLiftCoef map[Id]*float64
	AeroDragCoef map[Id]*float64

	//
	// Ship Class Component
	///
	SCC map[Id]ShipClass
}

func ResetState() {
	s := new(State)

	channels := make([]chan<- []byte, 0)
	s.MB = &MessageBus{channels}

	channels2 := make([]chan<- []byte, 0)
	s.AB = &MessageBus{channels2}

	s.Ents = make(map[Id]bool, 10)
	s.HotEnts = make(map[Id]bool, 10)

	s.MassC = make(map[Id]*float64, 1)
	s.PC = make(map[Id]*V3, 1)

	s.OC = make(map[Id]*Q, 1)

	s.MC = make(map[Id]Mobile, 1)
	s.RC = make(map[Id]Rotational, 1)

	s.SCC = make(map[Id]ShipClass, 1)

	S = s
}

// TODO: auto-increment entity ID and decouple its component
func (s *State) NewEntity(e Id) {
	s.OC[e] = new(Q)

	fgs := make([]ForceGen, 0)
	S.MC[e] = Mobile{new(V3), &fgs}

	s.RC[e] = Rotational{new(V3), new(M3), new(M3), new(M4)}
}

func (s *State) AddForceGen(e Id, fg ForceGen) {
	*(s.MC[e].FGs) = append(*(s.MC[e].FGs), fg)
}

type EntJSON struct {
	Id  Id       `json: "id"`
	Mas *float64 `json: "mass"`
	Pos *V3      `json: "pos"`
	Vel *V3      `json: "vel"`
	Ori *Q       `json: "ori"`
	Rot *V3      `json: "rot"`
}

type RefFrameJSON struct {
	Ents []EntJSON `json: "ents"`
}

// Encode state as an array of reference frames, each having
// an array of entities where each entity has mass, position, etc.
// TODO: for now, we assume all entities have all components
func (s *State) MarshalJSON() ([]byte, error) {
	//log.Debug("MarshalJSON")
	rfJSONs := make([]RefFrameJSON, 0)
	rfJSON := RefFrameJSON{}
	for eId, mass := range s.MassC {
		entJSON := EntJSON{
			eId,
			mass,
			s.PC[eId],
			s.MC[eId].V,
			s.OC[eId],
			s.RC[eId].R,
		}
		rfJSON.Ents = append(rfJSON.Ents, entJSON)
	}
	rfJSONs = append(rfJSONs, rfJSON)
	return json.Marshal(rfJSONs)
}
