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
	"time"

	"github.com/ethereum/go-ethereum/log"
)

func InitWorld() {
}

func DevWorld(testSeed uint64) {
	seed := testSeed
	if testSeed == 0 {
		seed = uint64(time.Now().Nanosecond())
	}

	r, _ := NewRand(seed)
	Rand = r

	ResetState()

	//for i := 0; i < 40; i++ {
	//	fmt.Println(starName())
	//}
	DevWorldStars()
	DebugSectors(true)

	// TODO: begin stationary relative top-level galactic grid
	// TODO: implement hyperdrive in any direction; sector traversal triggering
	//       star procgen
	//
	// TODO: and then - hyperdrive to a new star; triggering system procgen! :D  for E
	//
	//DevHyperdrive()

	// TODO: this comes after system procgen
	//DevShipOrbit()

	systems := []System{
		&Physics{},
		//&Hyperdrive{},
	}
	actionChan := make(chan Action, 10)
	subChan := make(chan *EntitySub, 10)
	GE = &GameEngine{
		systems:    systems,
		actionChan: actionChan,
		subChan:    subChan,
	}

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

	GE.Loop()
}

func DevWorldStars() {
	solPos := &V3{0.1, 0.1, (4.2 * sectorSize) / gridUnit}
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

	solSector := GetSector(solPos)
	solSector.addStarFixed(sol, solPos)
	solSector.Mapped = 1.0

	// traverse a few nearby sectors to get a few nearby stars
	north1 := GetSector(&V3{0.0, 0.0, (5.1 * aupc) / gridUnit})
	north2 := GetSector(&V3{0.0, 0.0, (6.1 * aupc) / gridUnit})
	north3 := GetSector(&V3{0.0, 0.0, (7.1 * sectorSize) / gridUnit})

	sectors := []*Sector{north1, north2, north3}
	for _, s := range sectors {
		for i := 0; i < 4; i++ {
			s.Traverse()
		}
	}
}

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

	sol := S.StarsByName["Sol"]
	solRF := S.EntFrames[sol.Entity]
	S.EntFrames[e] = solRF

	devheim := &Planet{
		Entity:         S.NewEntity(),
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
func DevWorld() {
	fmt.Printf("Average: %v\n", MassHist.AvgMass)

	solSectorPos := &V3{0.0, 0.0, (4 * aupc) / gridUnit}
	solSector := GetSector(solSectorPos)

	solSector.Debug()

	solSector.Traverse()
	solSector.Debug()

	solSector.Traverse()
	solSector.Debug()

	solSector.Traverse()
	solSector.Debug()

	solSector.Traverse()
	solSector.Debug()

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
