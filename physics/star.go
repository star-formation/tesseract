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

package physics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strconv"

	xrand "golang.org/x/exp/rand"
	
	"github.com/star-formation/tesseract/lib"
)
var (
	massHistogramFile = "data/Galaxy_stellar_mass_histogram.txt"
)

type Body struct {
	Name string

	Mass   float64
	Radius float64

	Orbit    *OE
	Rotation *lib.V3

	Atmosphere *Atmosphere

	MagField float64
}

// Star is a unique star.
type Star struct {
	Entity uint64
	Body

	SpectralType rune
	Luminosity   float64
	SurfaceTemp  float64
}

// NewStar returns a procedurally generated star.  Many of the stars attributes
// are derived from the mass.
func NewStar(mass float64) *Star {
	star := &Star{}
	//star.Body.Name = starName()
	star.Body.Mass = mass
	star.Body.Radius = StarRadius(mass) * solarRadius

	star.SpectralType = SpectralType(mass)
	star.Luminosity = StarLum(mass) * solarLum
	star.SurfaceTemp = StarSurfaceTemp(star.Luminosity, star.Body.Radius)
	return star
}

func (s *Star) Debug() {
	fmt.Printf("Class: %v Mass: %.3f Radius: %.0f km Temp: %.0f K Lum: %.3g W\n", strconv.QuoteRune(s.SpectralType), s.Body.Mass, s.Body.Radius/1000, s.SurfaceTemp, s.Luminosity)
}

// DefaultOrbit returns an orbit suitable as destination for FTL drives.
func (s *Star) DefaultOrbit() *OE {
	e, i, Ω, ω, θ := 0.0, 0.0, 0.0, 0.0, 0.0
	μ := GravitationalConstant * s.Mass * solarMass
	//hzInner, hzOuter := s.HabitableZone()
	r := aum * 1.0 //(hzInner + (hzOuter-hzInner)/2.0) // middle of HZ

	// Eqn 2.44 (substitution for h)
	h := math.Sqrt((r + r*e*math.Cos(θ)) * μ)
	return &OE{h: h, μ: μ, e: e, i: i, Ω: Ω, ω: ω, θ: θ}
}

// https://www.planetarybiology.com/calculating_habitable_zone.html
// TODO: update to latest research
func (s *Star) HabitableZone() (float64, float64) {
	lum := StarLum(s.Body.Mass)
	ri := math.Sqrt(lum / 1.1)
	ro := math.Sqrt(lum / 0.53)
	return ri, ro
}

// https://en.wikipedia.org/wiki/Stellar_classification
func SpectralType(mass float64) rune {
	switch {
	case mass < 0.50:
		return 'M'
	case mass < 0.80:
		return 'K'
	case mass < 1.04:
		return 'G'
	case mass < 1.40:
		return 'F'
	case mass < 2.10:
		return 'A'
	case mass < 16.0:
		return 'B'
	default:
		return 'O'
	}
}

// TODO: update to the latest research on star mass radius relation
func StarRadius(mass float64) float64 {
	if mass < 1.66 {
		return 1.06 * math.Pow(mass, 0.945)
	} else {
		return 1.33 * math.Pow(mass, 0.555)
	}
}

// https://en.wikipedia.org/wiki/Mass%E2%80%93luminosity_relation
func StarLum(mass float64) float64 {
	var a, b float64
	switch {
	case mass < 0.43:
		a, b = 2.3, 0.23
	case mass < 2.0:
		a, b = 4, 1
	case mass < 20.0:
		a, b = 3.5, 1.5
	default:
		a, b = 1, 3200
	}
	return b * math.Pow(mass, a)
}

// https://en.wikipedia.org/wiki/Stefan%E2%80%93Boltzmann_law#Temperature_of_stars
func StarSurfaceTemp(lum, r float64) float64 {
	return math.Pow((lum / (4 * math.Pi * stefanBoltzmann * (r * r))), 0.25)
}

//
// Mass Histogram / Initial Stellar Mass Function (IMF)
//
// [1] https://en.wikipedia.org/wiki/Initial_mass_function
// [2] https://github.com/Azeret/galIMF
//
type MassRange struct {
	Start     float64
	Range     float64
	StarCount int
}

type MassHistogram struct {
	Ranges     []MassRange
	AvgMass    float64
	TotalStars int
}

// Data generated with [2]
func getMassHistogram() *MassHistogram {
	var err error
	fileContent, err := ioutil.ReadFile(massHistogramFile)
	if err != nil {
		panic(err)
	}

	br := bytes.NewReader(fileContent)
	csvReader := csv.NewReader(br)
	csvReader.Comma = ' '
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = 5

	ranges := make([]MassRange, 0)
	mh := &MassHistogram{ranges, 0.0, 0}
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		mRange, err := strconv.ParseFloat(r[1], 64)
		start, err := strconv.ParseFloat(r[3], 64)
		count, err := strconv.ParseUint(r[4], 10, 64)
		if err != nil {
			panic(err)
		}

		mr := MassRange{start, mRange, int(count)}
		mh.Ranges = append(mh.Ranges, mr)
		mh.TotalStars += int(count)
	}

	var m float64
	for _, r := range mh.Ranges {
		tm := float64(r.StarCount) * (r.Start + r.Range/2.0)
		m += tm
	}
	mh.AvgMass = m / float64(mh.TotalStars)

	return mh
}

func (mh *MassHistogram) randMass(r *xrand.Rand) float64 {
	x := r.Intn(mh.TotalStars)
	i, count := 0, 0
	for {
		count += mh.Ranges[i].StarCount
		if x <= count {

			return mh.Ranges[i].Start + mh.Ranges[i].Range*r.Float64()
		}
		i++
	}
}
