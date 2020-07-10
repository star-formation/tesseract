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

// Star system barycenters are fixed points in the galactic reference frame.
// In single-star systems, the barycenter equals the center of the star.
type StarSystem struct {
	Entity Id
	Name   string

	// For now, we only support single-star systems.
	// barycenter *V3 // Position fixed in galactic frame
	Star    *Star
	Planets []*Planet

	// Value between 0.0 and 1.0 denoting how much the star system has been
	// explored by players == procedurally generated
	Mapped float64
}

// TODO: expand to multiple-star systems
func NewStarSystem(starMass float64) *StarSystem {
	ss := &StarSystem{}
	ss.Star = NewStar(starMass)
	return ss
}

func (ss *StarSystem) Clone() interface{} {
	name := make([]byte, len(ss.Name))
	copy(name, ss.Name)

	planets := []*Planet{}
	for _, p := range ss.Planets {
		c := p.Clone().(*Planet)
		planets = append(planets, c)
	}

	return &StarSystem{
		Entity:  ss.Entity,
		Name:    string(name),
		Star:    ss.Star.Clone().(*Star),
		Planets: planets,
		Mapped:  ss.Mapped,
	}
}

// Traverse enacts partial procedural generation of a star system.
// Traverse is generally called by the TODO system when a player is traversing
// the star system (having spent X time and/or moved Y distance in it).
func (ss *StarSystem) Traverse() {
	if ss.Mapped == 1.0 {
		return
	}

	if ss.Mapped == 0.0 {
		ss.procgenPlanets()
		// TODO: init stuff

	}

	// TODO: (scale AU distances to star)
	// 1. inner planets
	// 1.1. mercury-like, some even hotter / more active, volcano/lava
	// 1.2. hot jupiters (make for interesting gas harvesting)
	//                    - intense star light makes viable solar panels
	//                    - heat mgt becomes critical
	//
	// 2. habitable zone
	//
	// 3. outer planets
	//
	// 4. Kuiper Belt equivalent
	// 4.1. Ice worlds, crystalline ice
	//      - https://en.wikipedia.org/wiki/Haumea (challenging rotation!)
	//

	prob := 0.5 // TODO
	if Rand.Float64() < prob {
		// TODO: add planets
	}

	ss.Mapped += starSystemTraversalFactor
	if ss.Mapped > 1.0 {
		ss.Mapped = 1.0
	}
}

// TODO: temp testing: full procgen on init:
// TODO: instead of using HZ boundaries and arbitrary subdivision of it,
//       simply calculate the effective temp based on stellar flux
//       (determined by distance to star and the star's lum)
//       - then, adjust temp based on atmosphere, which is procgen
//         in part from magnetic field
func (ss *StarSystem) procgenPlanets() {
	planets := []*Planet{}

	// TODO: step 1: exoplanet occurence rates
	// planetCount :=
	ss.Planets = planets
}
