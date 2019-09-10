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

/*
 See https://en.wikipedia.org/wiki/Orbital_elements

 To reduce code verbosity in physics code, we use this naming:
 E = Eccentricity
 S = Semimajor Axis
 I = Inclination
 L = Longitude of the ascending node
 A = Argument of periapsis
 T = True anomaly

*/
type OE struct {
    E,S,I,L,A,T float64
}

type OEComp struct { OEs map[RefFrame]map[Id]OE }
