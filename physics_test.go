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
	"testing"
)

func TestPhysics(t *testing.T) {
	engine := setup()
	engine.Loop()
}

func setup() *Engine {
	// grid is top-level ref frame
	grid := &RefFrame{Parent: nil, Position: nil, Orbit: nil, Orientation: nil}

	// two entities
	e0 := Id(42)
	e1 := Id(43)
	e2 := Id(44)
	ents := make(map[Id]bool)
	ents[e0] = true
	ents[e1] = true
	ents[e2] = true

	// TODO: more granular setup functions
	ps := make(map[*RefFrame]map[Id]*V3)
	pComp := PComp{D: ps}
	pComp.D[grid] = make(map[Id]*V3)

	vs := make(map[*RefFrame]map[Id]*V3)
	vComp := VComp{D: vs}
	vComp.D[grid] = make(map[Id]*V3)

	// X,Y,Z coordinates in meters
	pComp.D[grid][e0] = &V3{1000.0, 1000.0, 1000.0}
	pComp.D[grid][e1] = &V3{1000.0, 1000.0, 1000.0}
	pComp.D[grid][e2] = &V3{1000.0, 1000.0, 1000.0}

	// X,Y,Z,M velocity vector, M is meters per second (m/s)
	vComp.D[grid][e0] = &V3{1, 0, 0}
	vComp.D[grid][e1] = &V3{1, 1, 0}
	vComp.D[grid][e2] = &V3{1, 1, 1}

	physics := &Physics{ents: ents, pComp: &pComp, vComp: &vComp}

	// "hot" ref frames and entities
	frames := []*RefFrame{grid}
	in := make(map[*RefFrame]*[]Id)
	in[grid] = &[]Id{e0, e1, e2}
	hotEnts := &HotEnts{Frames: frames, In: in}

	engine := &Engine{systems: []System{physics},
		mainBus: &MessageBus{},
		hot:     hotEnts}

	return engine
}
