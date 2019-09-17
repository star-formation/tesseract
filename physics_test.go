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

/*
func TestTorque(t *testing.T) {
	ResetState()
	grid := S.NewFrame()
	e0 := Id(0)
	S.NewEntity(e0, grid)

	var m0 float64
	m0 = 42000.0
	S.MC[grid][e0] = &m0

	S.PC[grid][e0] = &V3{10.0, 10.0, 10.0}

	lf := &V3{0.0, 0.0, 1.0}
	worldPoint := &V3{14.0, 14.0, 10.0}
	tq := S.TorqueAtPoint(e0, grid, lf, worldPoint)
	log.Debug("TestTorque", "lfm", lf.Magnitude(), "t", tq, "tm", tq.Magnitude())
}
*/

func setup() *Engine {
	ResetState()
	// grid is top-level ref frame
	grid := S.NewFrame()
	grid.DragCoef1 = 0.2
	grid.DragCoef2 = 0.4

	e0 := Id(42)

	S.NewEntity(e0, grid)

	S.Ents[e0] = true

	// Mass scalar value is kilogram (kg)
	var m0 float64
	m0 = 42000.0
	S.MC[grid][e0] = &m0

	// X,Y,Z position in meters (m) from origin
	S.PC[grid][e0] = &V3{0.0, 0.0, 0.0}

	S.AddForceGen(e0, grid, &ThrustForceGen{&V3{m0 * g0 * 1.0, 0, 0}})
	//S.AddForceGen(e0, grid, &DragForceGen{10.0, 40.0})

	// "hot" ref frames and entities
	frames := []*RefFrame{grid}
	in := make(map[*RefFrame]*[]Id)
	in[grid] = &[]Id{e0}
	hotEnts := &HotEnts{Frames: frames, In: in}

	engine := &Engine{
		systems: []System{&Physics{}},
		mainBus: &MessageBus{},
		hot:     hotEnts,
	}

	return engine
}
