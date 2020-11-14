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
	//"github.com/star-formation/tesseract/lib"
	"github.com/star-formation/tesseract/physics"
)

const (
	//
	// Physics Engine
	//
	linearDamping  = float64(1.0)
	angularDamping = float64(1.0)

	//
	// Game Design
	//
	gridUnit              = 100.0 // AU
	sectorSize            = physics.AUPC
	minStellarProximity   = (1.5 * physics.AULY) / gridUnit
	sectorTraversalFactor = 0.25

	maxPlanets = 14
	maxMoons   = 6
)

var (
	exoplanetEUCatalogFile = "data/exoplanet.eu_catalog_2019_12_02.csv"

	//
	// Static Game State
	//
	rootRF = &physics.RefFrame{
		Parent:      nil,
		Pos:         nil,
		Orbit:       nil,
		Orientation: nil,
	}
)
