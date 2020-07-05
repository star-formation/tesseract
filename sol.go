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

import "time"

// SolarSystem returns our solar system.
// t is the system time that sets the true anomaly of the orbital bodies.
func SolarSystem(t *time.Time) *StarSystem {
	sunLum := starLum(1.0)
	sun := &Star{
		Entity: S.NewEntity(),
		Body: Body{
			Name:       "Sun",
			Mass:       1.0,
			Radius:     1.0,
			Orbit:      nil,
			Rotation:   nil,
			Atmosphere: nil, // TODO
			MagField:   0.0, // TODO
		},
		SpectralType: spectralType(1.0),
		Luminosity:   sunLum,
		SurfaceTemp:  starSurfaceTemp(sunLum, 1.0),
	}

	mercury := &Planet{
		Entity: S.NewEntity(),
		Body: Body{
			Name:   "Mercury",
			Mass:   0.055,  // Earth masses
			Radius: 0.3829, // Earth radii
			Orbit: &OE{
				h: 0.0,  // set by engine in real-time (derived from true anomaly)
				i: 6.34, // degrees to invariable plane
				Ω: 48.331,
				e: 0.20563,
				ω: 29.124,
				θ: 0.0, // set by engine in real-time (derived from time)
				μ: sunMu,
			},
			AxialTilt:      0.027,              // degrees
			Rotation:       &V3{0.0, 0.0, 0.0}, // TODO: from axial tilt?
			RotationPeriod: 58.646,             // days
			Atmosphere:     nil,
			MagField:       300.0, // equatorial field strength in nano teslas
		},
	}

	venus := &Planet{
		Entity: S.NewEntity(),
		Body: Body{
			Name:   "Venus",
			Mass:   0.815,
			Radius: 0.9499,
			Orbit: &OE{
				h: 0.0,
				i: 2.19,
				Ω: 76.68,
				e: 0.006772,
				ω: 54.884,
				θ: 0.0,
				μ: sunMu,
			},
			AxialTilt:      2.64,
			Rotation:       &V3{0.0, 0.0, 0.0},
			RotationPeriod: -243.025, // days, negative period == retrograde
			Atmosphere: &Atmosphere{
				Height:           400.0, // TODO: lowest stable orbit
				ScaleHeight:      15.9,
				PressureSeaLevel: 92.1,
			},
			MagField: 0.0, // Venus is known not to have a magnetic field
		},
	}

	earth := &Planet{
		Entity: S.NewEntity(),
		Body: Body{
			Name:   "Earth",
			Mass:   1.0,
			Radius: 1.0,
			Orbit: &OE{
				h: 0.0,
				i: 1.57869,
				Ω: -11.26064,
				e: 0.0167086,
				ω: 114.20783,
				θ: 0.0,
				μ: sunMu,
			},
			AxialTilt:      23.4392811,
			Rotation:       &V3{0.0, 0.0, 0.0},
			RotationPeriod: 0.99726968,
			Atmosphere: &Atmosphere{
				Height:           200.0, // TODO: lowest stable orbit
				ScaleHeight:      8.5,
				PressureSeaLevel: 1.013,
			},
			MagField: 25000.0, // TODO: approx between 25,000 and 65,000
		},
	}

	mars := &Planet{
		Entity: S.NewEntity(),
		Body: Body{
			Name:   "Mars",
			Mass:   0.107,
			Radius: 0.533,
			Orbit: &OE{
				h: 0.0,
				i: 1.67,
				Ω: 49.558,
				e: 0.0934,
				ω: 286.502,
				θ: 0.0,
				μ: sunMu,
			},
			AxialTilt:      25.19,
			Rotation:       &V3{0.0, 0.0, 0.0},
			RotationPeriod: 1.025957,
			Atmosphere: &Atmosphere{
				Height:           125.0, // TODO: lowest stable orbit
				ScaleHeight:      10.8,
				PressureSeaLevel: 0.006,
			},
			MagField: 0.0,
		},
	}

	jupiter := &Planet{
		Entity: S.NewEntity(),
		Body: Body{
			Name:   "Jupiter",
			Mass:   317.8,
			Radius: 11.2, // TODO: mean radius / flattening
			Orbit: &OE{
				h: 0.0,
				i: 0.32,
				Ω: 100.464,
				e: 0.0489,
				ω: 273.867,
				θ: 0.0,
				μ: sunMu,
			},
			AxialTilt:      3.13,
			Rotation:       &V3{0.0, 0.0, 0.0},
			RotationPeriod: 9.925 / 24.0,
			/*
				Atmosphere: &Atmosphere{ // TODO: gas planet atmosphere handling!
					Height:           0.0, // TODO: lowest stable orbit
					ScaleHeight:      27.0,
					PressureSeaLevel: 0.0,
				},
			*/
			Atmosphere: nil,
			MagField:   420000.0, // TODO: 4.2 gauss at equator to 10-14 at the poles
		},
	}

	planets := []*Planet{mercury, venus, earth, mars, jupiter}

	solarSystem := &StarSystem{
		Entity:  S.NewEntity(),
		Name:    "Solar System",
		Star:    sun,
		Planets: planets,
		Mapped:  1.0,
	}

	return solarSystem
}
