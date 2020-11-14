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
	"github.com/star-formation/tesseract/physics"
)

// Types and Interfaces for the Entity-Component-System (ECS)
// architectural pattern.
// See: https://en.wikipedia.org/wiki/Entity_component_system
//
// If Gustav had invented the Entity-Component-System (ECS) pattern,
// he would have named it Id-Data-Function (IDF), as in ECS
// entities are ids, components are data and systems are functions.

// Entities are ids.
//type Id uint64

// Components are data containers for entities.  They contain no game logic.
type Component interface {
	// TODO
}

// Systems contain and execute game logic.
// Systems operate on and manipulate component data (game state).
//
// Each system handles a subset of components, which may overlap with other
// component subsets handled by other systems.
//
// Examples of systems include physics, combat and player trading exchanges.
//
// While components hold the game state data, systems track some data such as
// which entities "are in the system" (have data in all related components)
type System interface {
	// Init is called by the game engine once before the game loop begins.
	// It should perform any initialization needed by the system.
	Init() error

	// Update is called by the game engine once for each hot ref frame
	// every game frame.  Update is not called if there are no hot ref frames
	// (all game entities idle) in a given game frame.
	// The system must perform any state updates for all hot entities in the
	// hot ref frame.  Depending on the system, it may also update state for
	// idle entities in the hot ref frame (e.g. if an area of effect weapon
	// hits idle entities).
	Update(worldTime, elapsed float64, rf *physics.RefFrame) error

	// IsHotPostUpdate is called by the game engine once for each hot entity
	// every game frame.  IsHotPostUpdate is not called if there are no hot
	// entities in the given game frame (all entities idle).
	// The system must return if the given entity remains hot after the last
	// update.
	//
	// Example: the classical mechanics system returns true for any entity
	//          that has active force generators and false otherwise.
	IsHotPostUpdate(uint64) bool
}

func removeEnt(ents []uint64, i int) []uint64 {
	ents[i] = ents[len(ents)-1]
	return ents[:len(ents)-1]
}
