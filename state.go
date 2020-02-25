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

// Dev/Test global in-memory game state.  Used to simplify
// rapid iteration of data structures and game design/logic.

// The state will likely go through many iterations before it's clear
// how to best encode and optimize it for merkle trees/proofs
// in the outer blockchain layer.
type Mobile struct {
	V   *V3         // velocity
	FGs *[]ForceGen // force/torque generators
}

// If an object can rotate it has an inertia tensor
type Rotational struct {
	R    *V3 // Rotation X,Y,Z 3D vector
	IITB *M3 // Inverse Inertia Tensor Body space/coordinates
	IITW *M3 // Inverse Inertia Tensor World space/coordinates
	T    *M4 // 3x4 matrix transforming body space/coordinates to world
}

var S *State

type State struct {
	MB *MessageBus
	AB *MessageBus

	EntCount uint64

	EntsInFrames map[*RefFrame]map[Id]bool

	// TODO: track frames one level below root
	// TODO: think about search, log N complexity
	GalacticFrames map[Id]*RefFrame

	// Hyperspace component holds data used by the Hyperdrive System
	HSC map[Id]*Hyperspace

	// TODO: consolidate
	StarC map[Id]*Star
	Stars map[string]*Star

	AllSectors map[string]*Sector

	//
	// Physics Components
	//
	// Mass component holds float64 values
	MassC map[Id]*float64

	// Position and Velocity components holds 3x1 vectors
	PC map[Id]*V3

	// Orbit Component holds Keplerian Orbital Elements
	ORBC map[Id]*OE

	// Orientation Component holds quaternions
	ORIC map[Id]*Q

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

	s.EntsInFrames = make(map[*RefFrame]map[Id]bool, 0)
	s.GalacticFrames = make(map[Id]*RefFrame, 0)
	s.HSC = make(map[Id]*Hyperspace, 0)
	s.StarC = make(map[Id]*Star, 0)
	s.Stars = make(map[string]*Star, 0)
	s.AllSectors = make(map[string]*Sector, 0)
	s.MassC = make(map[Id]*float64, 0)
	s.PC = make(map[Id]*V3, 0)
	s.ORBC = make(map[Id]*OE, 0)
	s.ORIC = make(map[Id]*Q, 0)
	s.MC = make(map[Id]Mobile, 0)
	s.RC = make(map[Id]Rotational, 0)
	s.SRC = make(map[Id]*float64, 0)
	s.SCC = make(map[Id]ShipClass, 0)
	S = s
}

func (s *State) NewEntity() Id {
	s.EntCount += 1
	return Id(s.EntCount)
}

func (s *State) AddStar(star *Star, pos *V3) {
	s.PC[star.Entity] = pos
	s.StarC[star.Entity] = star
	s.Stars[star.Body.Name] = star
}

/*
func (s *State) NewEntity(e Id) {
	s.ORIC[e] = new(Q)

	fgs := make([]ForceGen, 0)
	S.MC[e] = Mobile{new(V3), &fgs}

	s.RC[e] = Rotational{new(V3), new(M3), new(M3), new(M4)}
}
*/

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
	rfJSONs := make([]RefFrameJSON, 0)
	rfJSON := RefFrameJSON{}
	for eId, mass := range s.MassC {
		entJSON := EntJSON{
			eId,
			mass,
			s.PC[eId],
			s.MC[eId].V,
			s.ORIC[eId],
			s.RC[eId].R,
		}
		rfJSON.Ents = append(rfJSON.Ents, entJSON)
	}
	rfJSONs = append(rfJSONs, rfJSON)
	return json.Marshal(rfJSONs)
}
