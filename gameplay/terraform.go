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

package gameplay

/*
import (
	"encoding/csv"
	"io/ioutil"
	"strings"

	//"github.com/star-formation/tesseract/lib"
	"github.com/star-formation/tesseract/physics"
)

var (
	MassHist *physics.MassHistogram
)

func init() {
	MassHist = getMassHistogram()
}

//
// PLANETS
//
type PlanetBase uint8

const (
	Terra = iota
	Gas
)

//
// MOONS
//
type Moon struct {
}

// Orbital elements generation is simplified by ensuring
// that the following always holds:
// * The Hill Spheres of two bodies never overlap; by comparing the
//   semi-major axes of two orbits
func rollOrbits(n uint) []*physics.OE {
	// This Hill Sphere intersection test is simplified in that it does not
	// take into account orbital inclination, longitude of the ascending node,
	// argument of periapsis nor synchronocy of orbital periods.
	// Instead, we simply assume the orbits are "flat" circles
	// in the orbital plane and that the two bodies' hill spheres
	// never overlap at their apastron.

	return []*physics.OE{}
}

//
// Planets
//
func rollPlanetCount() uint {
	return 1
}

func rollPlanets() []*physics.Planet {
	c := rollPlanetCount()
	//orbits := rollOrbits(c)
	planets := make([]*physics.Planet, c)
	for i := 0; i < int(c); i++ {
		planets = append(planets, rollPlanet())
	}
	return planets
}

func rollPlanet() *physics.Planet {
	return &physics.Planet{}
}

func rollTerraOrGas(mass, float64, orbit *physics.OE, star physics.Star) PlanetBase {
	return Terra
}

func rollPlanetMass(orbit *physics.OE) float64 {
	return 0
}

func genPlanetName(p *physics.Planet) string {
	if true {
		return rollNameEpic(p)
	} else {
		return ""
	}
}

func rollNameEpic(p *physics.Planet) string {
	return ""
}

//
// Moons
//
func rollMoonCount(p *physics.Planet) uint {
	return 1
}

func rollMoons(p *physics.Planet) []*Moon {
	c := rollMoonCount(p)
	//orbits := rollOrbits(c)
	moons := make([]*Moon, c)
	for i := 0; i < int(c); i++ {
		moons = append(moons, rollMoon())
	}
	return moons
}

func rollMoon() *Moon {
	return &Moon{}
}

//
// Exoplanet Catalog
//
/*
func parseExoplanetCatalog() {
	csvReader := csvReader(exoplanetEUCatalogFile)

	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			log.Info("EOF: ", "f", exoplanetEUCatalogFile)
			break
		}
		if err != nil {
			panic(err)
		}
	}

}

//
// File Utils
//
func csvReader(fileName string) *csv.Reader {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(strings.NewReader(string(b)))
	// TODO: skip all leading empty lines or lines starting with #
	_, _ = csvReader.Read() // Skip header line
	return csvReader
}
*/
