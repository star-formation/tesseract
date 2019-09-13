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
	p := NewPhysics()
	// grid is top-level ref frame
	grid := p.NewFrame()

	// two entities
	e0 := Id(42)
	e1 := Id(43)

	p.NewEntity(e0, grid)
	p.NewEntity(e1, grid)

	p.Ents[e0] = true
	p.Ents[e1] = true

	// X,Y,Z coordinates in meters (m)
	p.PC[grid][e0] = &V3{1000.0, 1000.0, 1000.0}
	p.PC[grid][e1] = &V3{1000.0, 1000.0, 1000.0}

	// X,Y,Z,M velocity vector, M is meters per second (m/s)
	p.VC[grid][e0] = &V3{1, 0, 0}
	p.VC[grid][e1] = &V3{0, 1, 0}

	// Mass scalar value is kilogram (kg)
	var m0, m1 float64
	m0, m1 = 20.0, 40.0
	p.MC[grid][e0] = &m0
	p.MC[grid][e1] = &m1

	// "hot" ref frames and entities
	frames := []*RefFrame{grid}
	in := make(map[*RefFrame]*[]Id)
	in[grid] = &[]Id{e0, e1}
	hotEnts := &HotEnts{Frames: frames, In: in}

	engine := &Engine{
		systems: []System{p},
		mainBus: &MessageBus{},
		hot:     hotEnts,
	}

	return engine
}
