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
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
)

var (
	PlanetMassHist   *Histogram
	PlanetRadiusHist *Histogram
)

type PlanetBase uint8

const (
	Lava = iota
	Scorched
	World
	Gas
)

type Planet struct {
	Entity Id
	Body
}

func NewPlanet() *Planet {
	p := &Planet{}
	p.Body.Name = planetName()
	return p
}

func (p *Planet) Clone() interface{} {
	return &Planet{
		Entity: p.Entity,
		Body:   *p.Body.Clone().(*Body),
	}
}

func planetName() string {
	// TODO: number in system and possibly unique/rare name
	return ""
}

// DefaultOrbit returns a circular, prograde orbit 100km above a planet's
// surface or above its atmosphere height (if it has one).
func (p *Planet) DefaultOrbit() *OE {
	e, i, Ω, ω, θ := 0.0, 0.0, 0.0, 0.0, 0.0
	μ := GravitationalConstant * p.Body.Mass * earthMass
	r := 100000.0
	if p.Body.Atmosphere != nil {
		r += p.Body.Atmosphere.Height
	}
	h := math.Sqrt((r + r*e*math.Cos(θ)) * μ)
	return &OE{h: h, μ: μ, e: e, i: i, Ω: Ω, ω: ω, θ: θ}
}

// https://en.wikipedia.org/wiki/Gravity_of_Earth#Altitude
// p.surface_gravity has been pre-calculated by world building scripts
func (p *Planet) GravityAtAltitude(alt float64) float64 {
	// TODO: for debugging...
	if alt < -p.Body.Radius {
		panic("kraken")
	}

	// TODO: handle negative altitude (caves, canyons, etc)
	if alt < 0 {
		return p.Body.SurfaceGravity
	}

	x := p.Body.Radius / (p.Body.Radius + alt)
	return p.Body.SurfaceGravity * (x * x)
}

// Where we lock X,Y,Z coords does not matter, as we do not yet have
// surface features.  Simply select a orientation derived from the
// planet's orientation 3D vector.
func (p *Planet) GeodeticToCartesian(lat, lon float64) float64 {
	return 0
}

// getExoplanetHistograms returns mass and radius histograms from
// a CSV export of http://exoplanet.eu/catalog/
func getExoplanetHistograms() (*Histogram, *Histogram) {
	var err error
	fileContent, err := ioutil.ReadFile(exoplanetEUCatalogFile)
	if err != nil {
		panic(err)
	}

	br := bytes.NewReader(fileContent)
	csvReader := csv.NewReader(br)
	csvReader.Comma = ','
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = 98

	masses, radii := []float64{}, []float64{}
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if r[2] != "" {
			mass, err := strconv.ParseFloat(r[2], 64)
			if err != nil {
				panic(err)
			}
			masses = append(masses, mass)
		}

		if r[8] != "" {
			radius, err := strconv.ParseFloat(r[8], 64)
			if err != nil {
				panic(err)
			}
			radii = append(radii, radius)
		}
	}

	fmt.Printf("masses: %d, radii: %d \n", len(masses), len(radii))

	sort.Float64s(masses)
	sort.Float64s(radii)

	/*
		for j := 0; j < 24; j++ {
			fmt.Printf("m: %.8f \n", masses[len(masses)-j-1])
		}
		for j := 0; j < 24; j++ {
			fmt.Printf("m: %.8f \n", masses[j])
		}
	*/

	massCounts := make([]uint64, len(masses))
	radiusCounts := make([]uint64, len(radii))

	for i := 0; i < len(masses); i++ {
		massCounts[i] = 1
	}
	for i := 0; i < len(radii); i++ {
		radiusCounts[i] = 1
	}

	massHist := NewHistogram(masses, massCounts)
	radiusHist := NewHistogram(radii, radiusCounts)
	return massHist, radiusHist
}
