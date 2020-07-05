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
	"strconv"
)

var (
	StarMassHist *Histogram
	StarMassAvg  float64
)

// Star is a unique star.
type Star struct {
	Entity Id
	Body

	SpectralType rune
	Luminosity   float64
	SurfaceTemp  float64
}

// NewStar returns a procedurally generated star.
func NewStar(mass float64) *Star {
	star := &Star{}
	star.Entity = S.NewEntity()
	star.Body.Name = starName()
	star.Body.Mass = mass
	star.Body.Radius = starRadius(mass) * solarRadius

	star.SpectralType = spectralType(mass)
	star.Luminosity = starLum(mass) * solarLum
	star.SurfaceTemp = starSurfaceTemp(star.Luminosity, star.Body.Radius)
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

// HabitableZone returns the inner and outer boundaries of the star's habitable zone in AU.
// See: https://www.planetarybiology.com/calculating_habitable_zone.html
// TODO: update to latest research
func (s *Star) HabitableZone() (float64, float64) {
	lum := starLum(s.Body.Mass)
	ri := math.Sqrt(lum / 1.1)
	ro := math.Sqrt(lum / 0.53)
	return ri, ro
}

// https://en.wikipedia.org/wiki/Stellar_classification
func spectralType(mass float64) rune {
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
func starRadius(mass float64) float64 {
	if mass < 1.66 {
		return 1.06 * math.Pow(mass, 0.945)
	} else {
		return 1.33 * math.Pow(mass, 0.555)
	}
}

// https://en.wikipedia.org/wiki/Mass%E2%80%93luminosity_relation
func starLum(mass float64) float64 {
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
func starSurfaceTemp(lum, r float64) float64 {
	return math.Pow((lum / (4 * math.Pi * stefanBoltzmann * (r * r))), 0.25)
}

// Mass Histogram / Initial Stellar Mass Function (IMF)
//
// [1] https://en.wikipedia.org/wiki/Initial_mass_function
// [2] https://github.com/Azeret/galIMF
//
// Data generated with [2]
func getStarStats() (*Histogram, float64) {
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

	endpoints := make([]float64, 0)
	counts := make([]uint64, 0)
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		start, err := strconv.ParseFloat(r[3], 64)
		count, err := strconv.ParseUint(r[4], 10, 64)
		if err != nil {
			panic(err)
		}

		endpoints = append(endpoints, start)
		counts = append(counts, count)
	}

	t := counts[0]
	m := float64(counts[0]) * (endpoints[0] / 2.0)
	for i := 1; i < len(endpoints); i++ {
		width := endpoints[i] - endpoints[i-1]
		m += float64(counts[i]) * (endpoints[i] + width/2.0)
		t += counts[i]
	}
	avgMass := m / float64(t)

	return NewHistogram(endpoints, counts), avgMass
}
