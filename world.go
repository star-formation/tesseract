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
	"fmt"
	"time"
)

func initRand(testSeed uint64) {
	seed := testSeed
	if testSeed == 0 {
		seed = uint64(time.Now().Nanosecond())
	}

	r, _ := NewRand(seed)
	Rand = r
}

func DevWorld2(testSeed uint64) {
	// init/reset state
	initRand(testSeed)
	ResetState()

	// our solar system and its galactic sector
	// TODO: correct time
	t := time.Now()
	ss := SolarSystem(&t)
	solPos := &V3{0.1, 0.1, auToGrid(4.2)}
	S.Pos[ss.Star.Entity] = solPos
	solSector := GetSector(solPos)
	solSector.addStarSystemFixed(ss, solPos)
	solSector.Mapped = 1.0
	DebugSectors(true)

	// setup player ship
	shipEnt := S.NewEntity()
	shipClass := &WarmJet{}

	var shipMassBase float64
	shipMassBase = shipClass.MassBase()

	shipIC := InertiaTensorCuboid(shipMassBase, 10, 10, 10)
	shipIC.Inverse()

	// add ship to state
	S.ShipClass[shipEnt] = shipClass
	S.Mass[shipEnt] = &shipMassBase
	S.Ori[shipEnt] = new(Q)
	S.Rot[shipEnt] = &Rotational{new(V3), new(M3), new(M3), new(M4)}
	S.Rot[shipEnt].IITB = shipIC

	fgs := make([]ForceGen, 0)
	S.ForceGens[shipEnt] = fgs

	// add player ship to Earth ref frame
	S.EntFrames[shipEnt] = S.EntFrames[ss.Planets[2].Entity]
	// set player ship orbit to low Earth orbit (LEO)
	S.Orb[shipEnt] = ss.Planets[2].DefaultOrbit()
	// TODO: ref frames and ship -> solar system links (for client API)
	// TODO: for now, manually link ref frames here;
	//       refactor after testing in client

	//
	go func() {
		time.Sleep(time.Second * 1)
		apiResp := APIGetGalaxy()
		j, err := json.MarshalIndent(apiResp, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("apiResp: %s\n", string(j))
	}()

	GE = NewEngine()
	GE.Loop()
}

/*
func DevWorldStars() {
	solPos := &V3{0.1, 0.1, auToGrid(4.2)}
	solLum := starLum(1.0)
	solTemp := starSurfaceTemp(solLum, 1.0)
	sol := &Star{
		Entity: S.NewEntity(),
		Body: Body{
			Name:       "Sun",
			Mass:       1.0,
			Radius:     1.0,
			Orbit:      nil,
			Rotation:   nil,
			Atmosphere: nil, // TODO
			MagField:   0.0, // TODO
		},
		SpectralType: spectralType(1.0),
		Luminosity:   solLum,
		SurfaceTemp:  solTemp,
	}
	log.Debug("FUNKY", "sol.Entity", sol.Entity)
	S.Pos[sol.Entity] = solPos

	solarSystem := &StarSystem{
		Entity:  S.NewEntity(),
		Name:    "Sol",
		Star:    sol,
		Planets: nil, // TODO
		Mapped:  1.0,
	}

	solSector := GetSector(solPos)
	solSector.addStarSystemFixed(solarSystem, solPos)
	solSector.Mapped = 1.0

	// traverse a few nearby sectors to get a few nearby stars
	north1 := GetSector(&V3{0.0, 0.0, auToGrid(5.1)})
	north2 := GetSector(&V3{0.0, 0.0, auToGrid(6.1)})
	north3 := GetSector(&V3{0.0, 0.0, auToGrid(7.1)})

	sectors := []*Sector{north1, north2, north3}
	for _, s := range sectors {
		for i := 0; i < 4; i++ {
			s.Traverse()
		}
		for _, ss := range s.StarSystems {
			ss.Traverse()
		}
	}
}
*/

/*
func DevShipOrbit() {
	// ==== SHIP
	e := S.NewEntity()
	log.Debug("DevShipOrbit", "shipEnt", e)
	shipClass := &WarmJet{}
	S.ShipClass[e] = shipClass

	var m0 float64
	m0 = shipClass.MassBase()
	S.Mass[e] = &m0

	S.Ori[e] = new(Q)

	S.Rot[e] = &Rotational{new(V3), new(M3), new(M3), new(M4)}
	shipIC := InertiaTensorCuboid(m0, 10, 10, 10)
	shipIC.Inverse()
	S.Rot[e].IITB = shipIC

	fgs := make([]ForceGen, 0)
	S.ForceGens[e] = fgs

	sol := S.StarSystemsByName["Sol"].Star
	solRF := S.EntFrames[sol.Entity]
	S.EntFrames[e] = solRF

	devheim := &Planet{
		Entity: S.NewEntity(),
		Body: Body{
			Mass:           1.0,
			Radius:         1.0 * earthRadius,
			AxialTilt:      0.0, // TODO: earth's
			RotationPeriod: 24 * 3600,
			SurfaceGravity: 1.0 * g0,
			Atmosphere: &Atmosphere{
				Height:           100 * 1000,
				ScaleHeight:      8.5 * 1000,
				PressureSeaLevel: 1.0 * earthSeaLevelPressure,
			},
		},
	}

	devheimRF := &RefFrame{
		Parent:      solRF,
		Pos:         nil,
		Orbit:       sol.DefaultOrbit(),
		Orientation: nil, // TODO
	}
	S.EntFrames[devheim.Entity] = devheimRF

	S.EntFrames[e] = devheimRF
	S.Orb[e] = devheim.DefaultOrbit()
	S.AddForceGen(e, &ThrustForceGen{m0 * g0 * 0.0001, 2})
}
*/

/*
func ShipAndStation() {
	// TODO: split into two dev worlds: one with this and one with ship
	//       starting in root frame / in hyperdrive

	// ==== SHIP
	ship := &WarmJet{}
	shipEnt := Id(42)

	S.SCC[shipEnt] = ship

	S.ORIC[shipEnt] = new(Q)

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
	S.ORIC[stationEnt] = new(Q)
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
*/

/*
func DevWorldBase() {
	solSectorPos := &V3{0.0, 0.0, (4 * aupc) / gridUnit}
	solLum := starLum(1.0)
	solTemp := starSurfaceTemp(solLum, 1.0)
	sol := &Star{
		Entity: S.NewEntity(),
		Body: Body{
			Name:       "Sol",
			Mass:       1.0,
			Radius:     1.0,
			Orbit:      nil,
			Rotation:   nil,
			Atmosphere: nil, // TODO
			MagField:   0.0, // TODO
		},
		SpectralType: spectralType(1.0),
		Luminosity:   solLum,
		SurfaceTemp:  solTemp,
	}

	earth := &Planet{
		Entity:         S.NewEntity(),
		Mass:           1.0,
		Radius:         1.0,
		SurfaceGravity: g0,
		Atmosphere: &Atmosphere{
			Height:           160000.0,
			ScaleHeight:      8500.0,
			PressureSeaLevel: 1.0,
		},
	}

	rootRF := &RefFrame{
		Parent:      nil,
		Pos:         nil,
		Orbit:       nil,
		Orientation: nil,
		DragCoef1:   0.0,
		DragCoef2:   0.0,
	}

	solRF := &RefFrame{
		Parent:      rootRF,
		Pos:         solSectorPos,
		Orbit:       nil,
		Orientation: nil, // TODO
		DragCoef1:   0.0,
		DragCoef2:   0.0,
	}

	earthOrbit := &OE{
		E: 0.0167,
		S: au,
		I: 7.155,
		L: -11.26064,
		A: 114.20783,
		T: 0.0, // TODO: time/epoch and starting position on the orbit.
	}
	earthRF := &RefFrame{
		Parent:      solRF,
		Pos:         nil,
		Orbit:       earthOrbit,
		Orientation: nil, // TODO
		DragCoef1:   0.0,
		DragCoef2:   0.0,
	}

	localOrbit := &OE{
		E: 0.0, // circular orbit
		S: earthRadius + 500,
		I: 0.0, // equatorial orbit
		L: 0.0, // TODO: planet reference frame direction/orientation
		A: 0.0, // TODO:
		T: 0.0, // TODO: time/epoch and starting position on the orbit.
	}
	localRF := &RefFrame{
		Parent:      earthRF,
		Pos:         nil,
		Orbit:       earthOrbit,
		Orientation: nil,  // TODO: test diff values
		Radius:      20.0, // km
		DragCoef1:   1.0,
		DragCoef2:   1.0,
	}

	S.EntFrames[sol.Entity] = solRF
}
*/

/*
	// TODO: read NewEntitySub and heimdall code, simplify
	// TODO: create API call for all relevant entity state + sub of future deltas
	// TODO: move to api.go (world.go should only setup world states)
	go func() {
		time.Sleep(2 * time.Second)
		dataChan, _ := NewEntitySub(2)
		for {
			select {
			case d, ok := <-dataChan:
				if ok {
					log.Debug("EntitySub", "len(dataChan)", len(dataChan), "d", string(d))
				} else {
					log.Debug("EntitySub dataChan closed")
					return
				}
			}
		}
	}()
*/
