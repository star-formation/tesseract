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

import "github.com/ethereum/go-ethereum/log"

func DevScene1() {
	ResetState()

	e0 := Id(42)
	S.NewEntity(e0)
	S.Ents[e0] = true
	S.HotEnts[e0] = true

	// Mass scalar value is kilogram (kg)
	var m0 float64
	m0 = 4200.0
	S.MassC[e0] = &m0

	// X,Y,Z position in kilometers (km) from origin
	S.PC[e0] = &V3{-200, 600, -200}

	// TODO: add ship class entry for e0
	S.SCC[e0] = &WarmJet{}

	ic := InertiaTensorCuboid(m0, 10, 10, 10)
	ic.Inverse()
	log.Debug("devscene", "ic", ic)
	*S.RC[e0].IITB = *ic
	//log.Debug("devscene", "S.ICB[e0]", S.ICB[e0])
	//S.AddForceGen(e0, &ThrustForceGen{&V3{m0 * g0 * 0.01, 0, 0}})
	//S.AddForceGen(e0, grid, &DragForceGen{10.0, 40.0})

	actionChan := make(chan Action, 10)
	GE = &GameEngine{
		systems:    []System{&Physics{}},
		actionChan: actionChan,
	}

	go GE.Loop()
}
