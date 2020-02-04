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
)

type Star struct {}

// Star system barycenters are fixed points in the galactic reference frame.
// In single-star systems, the barycenter equals the center of the star.
type StarSystem struct {
    barycenter *V3 // Position fixed in galactic frame
    stars []Star
    planets []Planet
}
