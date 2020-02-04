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

import "math"

func InitWorld() {
	ResetState()

	planetRF := &RefFrame{
		Parent:      nil,
		Pos:         nil,
		Orbit:       nil, //&OE{},
		Orientation: nil,
		DragCoef1:   0,
		DragCoef2:   0,
	}

	earth := &Planet{
		Entity:         0,
		Mass:           5.97237 * math.Pow(10, 24), // kg
		Radius:         6378,                       // km
		SurfaceGravity: 1.0,
		Atmosphere:     nil, // TODO

	}

	orbit := &OE{
		E: 0.0, // circular orbit
		S: earth.Radius + 500,
		I: 0.0, // equatorial orbit
		L: 0.0, // TODO: planet reference frame direction/orientation
		A: 0.0, // TODO: check
		T: 0.0, // TODO: starting position on the orbit. TODO: time
	}
	localRF := &RefFrame{
		Parent:      planetRF,
		Pos:         nil,
		Orbit:       orbit,
		Orientation: &Q{}, // TODO: test diff values
		Radius:      20.0, // km
		DragCoef1:   1.0,
		DragCoef2:   1.0,
	}

	// ==== SHIP
	ship := &WarmJet{}
	shipEnt := Id(42)

	S.SCC[shipEnt] = ship

	S.OC[shipEnt] = new(Q)

	fgs := make([]ForceGen, 0)
	S.MC[shipEnt] = Mobile{new(V3), &fgs}

	S.RC[shipEnt] = Rotational{new(V3), new(M3), new(M3), new(M4)}

	var m0 float64
	m0 = ship.MassBase()
	S.MassC[shipEnt] = &m0

	shipIC := InertiaTensorCuboid(m0, 10, 10, 10)
	shipIC.Inverse()

	*S.RC[shipEnt].IITB = *shipIC
	//log.Debug("devscene", "S.ICB[e0]", S.ICB[e0])
	//S.AddForceGen(e0, &ThrustForceGen{&V3{m0 * g0 * 0.01, 0, 0}})
	//S.AddForceGen(e0, grid, &DragForceGen{10.0, 40.0})

	// X,Y,Z position in kilometers (km) from origin
	S.PC[shipEnt] = &V3{1, 1, 1}
	r := ship.BoundingSphereRadius()
	S.SRC[shipEnt] = &r

	// ==== STATION
	// The Station is stationary; no mass, velocity, force gens, etc.
	stationEnt := Id(43)
	S.OC[stationEnt] = new(Q)
	S.PC[stationEnt] = &V3{0, 0, 0}
	stationRadius := 200.0
	S.SRC[stationEnt] = &stationRadius

	entMap := make(map[Id]bool, 2)
	S.EntsInFrames[localRF] = entMap

	S.EntsInFrames[localRF][stationEnt] = true
	S.EntsInFrames[localRF][shipEnt] = true

	actionChan := make(chan Action, 10)
	GE = &GameEngine{
		systems:    []System{&Physics{}},
		actionChan: actionChan,
	}

	go GE.Loop()
}
