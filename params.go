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
	"math"
	"time"
)

const (
	//
	// Math Constants
	//
	twoPi = 2 * math.Pi

	DBL_EPSILON = 2.2204460492503131E-16

	// acceptable tolerance of errors of float64 calculations compared to
	// analytical or precalculated solutions
	tolerance = 0.000000000000001

	eccentricAnomalyTolerance           = 1e-6
	hyperbolicEccentricAnomalyTolerance = 1e-6

	//
	// Physics and Astrophysics Constants (real world)
	//
	speedOfLight          = 299792458.0 // m/s
	GravitationalConstant = 6.674e-11
	stefanBoltzmann       = 5.670373e-8

	aum  = 149597870700.0    // meters per AU
	auly = 63241.07708426628 // AU per light year
	aupc = 2.06265e5         // AU per parsec
	lypc = 3.26156           // light years per parsec

	milkyWayDiscHeight = 2000.0 // ly
	radiusLocalBubble  = 150.0  // ly

	// https://en.wikipedia.org/wiki/Solar_mass
	// https://en.wikipedia.org/wiki/Solar_luminosity
	// https://en.wikipedia.org/wiki/Solar_radius
	solarMass   = 1.98855e30 // kg
	solarLum    = 3.828e26   // W
	solarRadius = 6.957e8    // m

	// https://en.wikipedia.org/wiki/Earth_mass
	earthMass             = 5.9722e24 // kg
	earthRadius           = 6.3781e6  // km
	earthMu               = 3.986004418e14
	earthSeaLevelPressure = 101325 // pascals
	g0                    = 9.80665

	marsMu = 4.282837e13

	//
	// Physics Engine
	//
	linearDamping  = float64(1.0)
	angularDamping = float64(1.0)

	//
	// Game Engine
	//
	loopTarget        = 1000 * time.Millisecond
	maxActionsPerLoop = 10

	//
	// Game Design
	//
	gridUnit              = 100.0 // AU
	sectorSize            = aupc
	minStellarProximity   = (1.5 * auly) / gridUnit
	sectorTraversalFactor = 0.25

	maxPlanets = 14
	maxMoons   = 6

	//
	// System Name Letter Relative Weights
	//
	// https://en.wikipedia.org/wiki/List_of_writing_systems#List_of_writing_scripts_by_adoption
	wLatin           = 6120
	wChinese         = 1340
	wDevanagari      = 820
	wArabic          = 660
	wBengaliAssamese = 300
	wCyrillic        = 250
	wKana            = 120
	wJavanese        = 80
	wHangul          = 79
	wTelugu          = 74
	wTamil           = 70
	wGujarati        = 48
)

var (
	//
	// Math
	//
	KHat = &V3{0, 0, 1}

	massHistogramFile = "data/Galaxy_stellar_mass_histogram.txt"

	exoplanetEUCatalogFile = "data/exoplanet.eu_catalog_2019_12_02.csv"

	//
	// Static Game State
	//
	rootRF = &RefFrame{
		Parent:      nil,
		Pos:         nil,
		Orbit:       nil,
		Orientation: nil,
	}
)
