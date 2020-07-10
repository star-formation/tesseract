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

type Body struct {
	Name string

	Mass   float64
	Radius float64

	Orbit *OE

	AxialTilt float64

	Rotation       *V3
	RotationPeriod float64

	Atmosphere *Atmosphere

	MagField float64

	SurfaceGravity float64
}

func (b *Body) Clone() interface{} {
	name := make([]byte, len(b.Name))
	copy(name, b.Name)

	var atm *Atmosphere
	if b.Atmosphere != nil {
		atm = b.Atmosphere.Clone().(*Atmosphere)
	}

	var orb *OE
	if b.Orbit != nil {
		orb = b.Orbit.Clone().(*OE)
	}

	var rot *V3
	if b.Rotation != nil {
		rot = b.Rotation.Clone().(*V3)
	}

	return &Body{
		string(name),
		b.Mass,
		b.Radius,
		orb,
		b.AxialTilt,
		rot,
		b.RotationPeriod,
		atm,
		b.MagField,
		b.SurfaceGravity,
	}
}
