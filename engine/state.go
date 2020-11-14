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
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"

	"github.com/star-formation/tesseract/lib"
	"github.com/star-formation/tesseract/physics"
	//"github.com/star-formation/tesseract/gameplay"
)

// Dev/Test global in-memory game state.  Used to simplify
// rapid iteration of data structures and game design/logic.

// The state will likely go through many iterations before it's clear
// how to best encode and optimize it for merkle trees/proofs
// in the outer blockchain layer.

// If an object can rotate it has an inertia tensor
type Rotational struct {
	R    *lib.V3 // Rotation X,Y,Z 3D vector
	IITB *lib.M3 // Inverse Inertia Tensor Body space/coordinates
	IITW *lib.M3 // Inverse Inertia Tensor World space/coordinates
	T    *lib.M4 // 3x4 matrix transforming body space/coordinates to world
}

var S *State

type State struct {
	MsgBus    *lib.MessageBus
	ActionBus *lib.MessageBus

	EntCount  uint64
	EntFrames map[uint64]*physics.RefFrame

	HotEnts  map[*physics.RefFrame]map[uint64]struct{}
	IdleEnts map[*physics.RefFrame]map[uint64]struct{}

	IdleSince map[uint64]float64

	EntitySubs          map[*EntitySub]struct{}
	EntitySubsCloseChan chan *EntitySub

	//StarsById   map[uint64]*physics.Star
	//StarsByName map[string]*physics.Star
	//Sectors map[string]*gameplay.Sector

	//
	// Physics Components
	//
	// Mass component holds float64 values
	Mass map[uint64]*float64

	// Position and Velocity components holds 3x1 vectors
	Pos map[uint64]*lib.V3
	Vel map[uint64]*lib.V3

	// Orbit Component holds Keplerian Orbital Elements
	Orb map[uint64]*physics.OE

	// Orientation Component holds quaternions
	Ori map[uint64]*lib.Q

	// Holds rotational data for entities that can rotate
	Rot map[uint64]*Rotational

	// Holds force/torque generators for movable entities
	ForceGens map[uint64][]ForceGen

	//
	// Gameplay Components
	//
	// Holds ship class data
	//ShipClass map[uint64]ShipClass
}

func ResetState() {
	s := new(State)

	channels := make([]chan<- []byte, 0)
	s.MsgBus = &lib.MessageBus{channels}
	channels2 := make([]chan<- []byte, 0)
	s.ActionBus = &lib.MessageBus{channels2}

	s.EntitySubs = make(map[*EntitySub]struct{}, 0)
	s.EntitySubsCloseChan = make(chan *EntitySub, 100)

	s.EntFrames = make(map[uint64]*physics.RefFrame, 0)

	s.HotEnts = make(map[*physics.RefFrame]map[uint64]struct{}, 0)
	s.IdleEnts = make(map[*physics.RefFrame]map[uint64]struct{}, 0)
	s.IdleSince = make(map[uint64]float64, 0)

	//s.StarsById = make(map[uint64]*physics.Star, 0)
	//s.StarsByName = make(map[string]*physics.Star, 0)
	//s.Sectors = make(map[string]*gameplay.Sector, 0)
	
	s.Mass = make(map[uint64]*float64, 0)
	s.Pos = make(map[uint64]*lib.V3, 0)
	s.Vel = make(map[uint64]*lib.V3, 0)
	s.Orb = make(map[uint64]*physics.OE, 0)
	s.Ori = make(map[uint64]*lib.Q, 0)
	s.ForceGens = make(map[uint64][]ForceGen, 0)
	s.Rot = make(map[uint64]*Rotational, 0)
	//s.ShipClass = make(map[uint64]ShipClass, 0)
	S = s
}

func (s *State) NewEntity() uint64 {
	s.EntCount += 1
	return uint64(s.EntCount)
}

func (s *State) SetHot(e uint64, rf *physics.RefFrame) {
	s.ensureEntAlloc(rf)
	s.HotEnts[rf][e] = struct{}{}
	delete(s.IdleEnts[rf], e)
}

func (s *State) SetIdle(e uint64, rf *physics.RefFrame, since float64) {
	s.ensureEntAlloc(rf)
	delete(s.HotEnts[rf], e)
	s.IdleEnts[rf][e] = struct{}{}
	s.IdleSince[e] = since
}

func (s *State) ensureEntAlloc(rf *physics.RefFrame) {
	if rf == nil {
		panic("nil rf")
	}
	if s.HotEnts[rf] == nil {
		s.HotEnts[rf] = make(map[uint64]struct{}, 1)
	}
	if s.IdleEnts[rf] == nil {
		s.IdleEnts[rf] = make(map[uint64]struct{}, 1)
	}
}

/*
func (s *State) AddStar(star *physics.Star, pos *lib.V3) {
	s.Pos[star.Entity] = pos
	s.StarsById[star.Entity] = star
	s.StarsByName[star.Body.Name] = star

	s.SetIdle(star.Entity, S.EntFrames[star.Entity], 0)
}
*/

func (s *State) AddForceGen(e uint64, fg ForceGen) {
	s.ForceGens[e] = append(s.ForceGens[e], fg)
	rf := s.EntFrames[e]
	log.Debug("FUNK", "rf", rf)
	s.SetHot(e, rf)
}

func (s *State) AddEntitySub(e uint64) {
	//if S.EntitySubs[e]
}

//
// JSON Encoding
//
type EntJSON struct {
	Id  uint64       `json: "id"`
	Mas *float64 `json: "mass"`
	Pos *lib.V3      `json: "pos"`
	Vel *lib.V3      `json: "vel"`
	Ori *lib.Q       `json: "ori"`
	Rot *lib.V3      `json: "rot"`
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
