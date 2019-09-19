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

type FComp map[*RefFrame]map[Id]*float64 // Generic float64 component
type V3Comp map[*RefFrame]map[Id]*V3     // Generic 3x1 vector component
type M3Comp map[*RefFrame]map[Id]*M3     // Generic 3x3 matrix component
type QComp map[*RefFrame]map[Id]*Q       // Generic quaternion component

type RadiusComp map[*RefFrame]map[Id]*float64

type HotEnts struct {
	Frames []*RefFrame
	In     map[*RefFrame]*[]Id
}

// TODO: JSON tags and encoding for entire state.  Easiest during dev/test.
//       Replaced by compact, fast binary encoding for prod.
type State struct {
	MB *MessageBus

	RefFrames []*RefFrame

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

func ResetState() {
	s := new(State)

	channels := make([]chan<- []byte, 0)
	s.MB = &MessageBus{channels}

	s.Ents = make(map[Id]bool, 10)

	s.MC = make(FComp, 1)

	s.PC = make(V3Comp, 1)
	s.VC = make(V3Comp, 1)
	s.RC = make(V3Comp, 1)

	s.OC = make(QComp, 1)

	s.IC = make(M3Comp, 1)

	s.FC = make(map[*RefFrame]map[Id][]ForceGen, 1)

	S = s
}

func (s *State) NewFrame() *RefFrame {
	rf := new(RefFrame)

	s.MC[rf] = make(map[Id]*float64, 10)

	s.PC[rf] = make(map[Id]*V3, 10)
	s.VC[rf] = make(map[Id]*V3, 10)
	s.RC[rf] = make(map[Id]*V3, 10)

	s.OC[rf] = make(map[Id]*Q, 10)

	s.IC[rf] = make(map[Id]*M3, 10)

	s.FC[rf] = make(map[Id][]ForceGen, 10)

	return rf
}

// TODO: auto-increment entity ID and decouple its component
func (s *State) NewEntity(e Id, rf *RefFrame) {
	s.VC[rf][e] = &V3{} // velocity

	s.IC[rf][e] = &M3{} // inertia tensor

	s.OC[rf][e] = &Q{}  // orientation
	s.RC[rf][e] = &V3{} // rotation

	s.FC[rf][e] = []ForceGen{} // Force Generators
}

func (s *State) AddForceGen(e Id, rf *RefFrame, fg ForceGen) {
	s.FC[rf][e] = append(s.FC[rf][e], fg)
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
	for _, rf := range s.RefFrames {
		rfJSON := RefFrameJSON{}
		for eId, m := range s.MC[rf] {
			entJSON := EntJSON{
				eId,
				m,
				s.PC[rf][eId],
				s.VC[rf][eId],
				s.OC[rf][eId],
				s.RC[rf][eId],
			}
			rfJSON.Ents = append(rfJSON.Ents, entJSON)
		}

		rfJSONs = append(rfJSONs, rfJSON)
	}
	return json.Marshal(rfJSONs)
}
