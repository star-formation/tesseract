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
	"encoding/csv"
	"io/ioutil"
	"strings"

	xrand "golang.org/x/exp/rand"
)

var (
	Rand *xrand.Rand
)

// Orbital elements generation is simplified by ensuring
// that the following always holds:
// * The Hill Spheres of two bodies never overlap; by comparing the
//   semi-major axes of two orbits
func rollOrbits(n uint) []*OE {
	// This Hill Sphere intersection test is simplified in that it does not
	// take into account orbital inclination, longitude of the ascending node,
	// argument of periapsis nor synchronocy of orbital periods.
	// Instead, we simply assume the orbits are "flat" circles
	// in the orbital plane and that the two bodies' hill spheres
	// never overlap at their apastron.

	return []*OE{}
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
*/
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

//naclbox "github.com/kevinburke/nacl/box"

/* TODO: this is for testing

   This will be replaced with a on-chain random beacon for source of
   deterministic _and_ unpredictable entropy.

   See https://dfinity.org/static/dfinity-consensus-0325c35128c72b42df7dd30c22c41208.pdf
   and https://github.com/ethereum/eth2.0-specs/blob/master/specs/core/0_beacon-chain.md
*/

func NewRand(seed uint64) (*xrand.Rand, error) {
	// We use https://www.godoc.org/golang.org/x/exp/rand#PCGSource
	// as the math/rand RNG algo is planned to be deprecated.
	// See https://github.com/golang/go/issues/21835
	//
	// x/exp/rand.NewSource defaults to PCGSource
	src := xrand.NewSource(seed)
	return xrand.New(src), nil
}
