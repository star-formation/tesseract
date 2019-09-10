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

    "github.com/ethereum/go-ethereum/log"
)

/*  To reduce code verbosity in physics code, we use this naming:
    X, Y, Z = Cartesian 3D coordinates

    P = 3D position    (X,Y,Z)
    V = 3D velocity    (X,Y,Z)
    O = 3D orientation (X,Y,Z)
    R = 3D rotation    (X,Y,Z)

    M = Magnitude (of velocity or rotation)
*/
type V3 struct {
    X, Y, Z float64
}

type PComp struct { Ps map[*RefFrame]map[Id]*V3 }
type VComp struct { Vs map[*RefFrame]map[Id]*V3 }
type OComp struct { Os map[*RefFrame]map[Id]*V3 }
type RComp struct { Rs map[*RefFrame]map[Id]*V3 }

type MassComp struct { Masses map[*RefFrame]map[Id]*float64 }

type RadiusComp struct { Radii map[*RefFrame]map[Id]*float64 }

// TODO ...
type HotEnts struct {
    Frames []*RefFrame
    In map[*RefFrame]*[]Id
}

// Physics Engine
type Physics struct {
    ents map[Id]bool
    pComp *PComp
    vComp *VComp
}

// System interface
func (p *Physics) Init() error {
    return nil
}
func (p *Physics) Update(elapsed float64, f *RefFrame, hotEnts *[]Id) error {
    for _, e := range *hotEnts {
        // For debug
        if !p.ents[e] {
            panic("hot ent not in frame")
        }
        
        // Update velocity
        p.pComp.Ps[f][e].X += (p.vComp.Vs[f][e].X * elapsed)
        p.pComp.Ps[f][e].Y += (p.vComp.Vs[f][e].Y * elapsed)
        p.pComp.Ps[f][e].Z += (p.vComp.Vs[f][e].Z * elapsed)

        // Detect collisions
        
        // Resolve collisions

        log.Info("physics.Update", "pComp", p.pComp.Ps[f][e])
    }
    return nil
}

func (p *Physics) RegisterEntity(id Id) {
    p.ents[id] = true
}
func (p *Physics) DeregisterEntity(id Id) {
    p.ents[id] = false
}

func magnitude(x, y, z float64) float64 {
    return math.Sqrt(x*x + y*y + z*z)
}