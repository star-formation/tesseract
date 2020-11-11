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

	"github.com/ethereum/go-ethereum/log"
)

// Dev/Test global in-memory game state.  Used to simplify
// rapid iteration of data structures and game design/logic.

// The state will likely go through many iterations before it's clear
// how to best encode and optimize it for merkle trees/proofs
// in the outer blockchain layer.

// If an object can rotate it has an inertia tensor
type Rotational struct {
	R    *V3 // Rotation X,Y,Z 3D vector
	IITB *M3 // Inverse Inertia Tensor Body space/coordinates
	IITW *M3 // Inverse Inertia Tensor World space/coordinates
	T    *M4 // 3x4 matrix transforming body space/coordinates to world
}

var S *State

type State struct {
	MsgBus    *MessageBus
	ActionBus *MessageBus

	EntCount  uint64
	EntFrames map[Id]*RefFrame

	HotEnts  map[*RefFrame]map[Id]struct{}
	IdleEnts map[*RefFrame]map[Id]struct{}

	IdleSince map[Id]float64

	EntitySubs          map[*EntitySub]struct{}
	EntitySubsCloseChan chan *EntitySub

	// Hyperspace component holds data used by the Hyperdrive System
	Hyperspace map[Id]*Hyperspace

	StarsById   map[Id]*Star
	StarsByName map[string]*Star

	Sectors map[string]*Sector

	//
	// Physics Components
	//
	// Mass component holds float64 values
	Mass map[Id]*float64

	// Position and Velocity components holds 3x1 vectors
	Pos map[Id]*V3
	Vel map[Id]*V3

	// Orbit Component holds Keplerian Orbital Elements
	Orb map[Id]*OE

	// Orientation Component holds quaternions
	Ori map[Id]*Q

	// Holds force/torque generators for movable entities
	ForceGens map[Id][]ForceGen

	// Holds rotational data for entities that can rotate
	Rot map[Id]*Rotational

	// Holds ship class data
	ShipClass map[Id]ShipClass
}

func ResetState() {
	s := new(State)

	channels := make([]chan<- []byte, 0)
	s.MsgBus = &MessageBus{channels}
	channels2 := make([]chan<- []byte, 0)
	s.ActionBus = &MessageBus{channels2}

	s.EntitySubs = make(map[*EntitySub]struct{}, 0)
	s.EntitySubsCloseChan = make(chan *EntitySub, 100)

	s.EntFrames = make(map[Id]*RefFrame, 0)

	s.HotEnts = make(map[*RefFrame]map[Id]struct{}, 0)
	s.IdleEnts = make(map[*RefFrame]map[Id]struct{}, 0)
	s.IdleSince = make(map[Id]float64, 0)

	s.Hyperspace = make(map[Id]*Hyperspace, 0)
	s.StarsById = make(map[Id]*Star, 0)
	s.StarsByName = make(map[string]*Star, 0)
	s.Sectors = make(map[string]*Sector, 0)
	s.Mass = make(map[Id]*float64, 0)
	s.Pos = make(map[Id]*V3, 0)
	s.Vel = make(map[Id]*V3, 0)
	s.Orb = make(map[Id]*OE, 0)
	s.Ori = make(map[Id]*Q, 0)
	s.ForceGens = make(map[Id][]ForceGen, 0)
	s.Rot = make(map[Id]*Rotational, 0)
	s.ShipClass = make(map[Id]ShipClass, 0)
	S = s
}

func (s *State) NewEntity() Id {
	s.EntCount += 1
	return Id(s.EntCount)
}

func (s *State) SetHot(e Id, rf *RefFrame) {
	s.ensureEntAlloc(rf)
	s.HotEnts[rf][e] = struct{}{}
	delete(s.IdleEnts[rf], e)
}

func (s *State) SetIdle(e Id, rf *RefFrame, since float64) {
	s.ensureEntAlloc(rf)
	delete(s.HotEnts[rf], e)
	s.IdleEnts[rf][e] = struct{}{}
	s.IdleSince[e] = since
}

func (s *State) ensureEntAlloc(rf *RefFrame) {
	if rf == nil {
		panic("nil rf")
	}
	if s.HotEnts[rf] == nil {
		s.HotEnts[rf] = make(map[Id]struct{}, 1)
	}
	if s.IdleEnts[rf] == nil {
		s.IdleEnts[rf] = make(map[Id]struct{}, 1)
	}
}

func (s *State) AddStar(star *Star, pos *V3) {
	s.Pos[star.Entity] = pos
	s.StarsById[star.Entity] = star
	s.StarsByName[star.Body.Name] = star

	s.SetIdle(star.Entity, S.EntFrames[star.Entity], 0)
}

func (s *State) AddForceGen(e Id, fg ForceGen) {
	s.ForceGens[e] = append(s.ForceGens[e], fg)
	rf := s.EntFrames[e]
	log.Debug("FUNK", "rf", rf)
	s.SetHot(e, rf)
}

func (s *State) AddEntitySub(e Id) {
	//if S.EntitySubs[e]
}

//
// JSON Encoding
//
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
	for eId, mass := range s.Mass {
		entJSON := EntJSON{
			eId,
			mass,
			s.Pos[eId],
			s.Vel[eId],
			s.Ori[eId],
			s.Rot[eId].R,
		}
		rfJSON.Ents = append(rfJSON.Ents, entJSON)
	}
	rfJSONs = append(rfJSONs, rfJSON)
	return json.Marshal(rfJSONs)
}
