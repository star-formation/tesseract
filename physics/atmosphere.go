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
	"math"
)

type Atmosphere struct {
	Height, ScaleHeight float64
	PressureSeaLevel    float64
}

// https://en.wikipedia.org/wiki/Scale_height
func (a *Atmosphere) PressureAtAltitude(alt float64) float64 {
	// TODO: handle increase pressure below sea level
	if alt < 0 {
		return a.PressureSeaLevel
	}

	if alt > a.Height {
		return 0
	}

	return a.PressureSeaLevel * math.Exp(-(alt / a.ScaleHeight))
}
