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
package gameplay

//import "github.com/ethereum/go-ethereum/log"

/*
//
// The hyperdrive system implements travel through hyperspace.
//
type Hyperdrive struct {
}

// Data for an instance of one ship in hyperdrive
type Hyperspace struct {
	// Start position in galactic grid units (see galaxy.go)
	Start *V3
	
	// Target the hyperdrive is locked onto
	// TODO: support non-star targets
	// TODO: support hyperdrive in arbitrary directions without target lock
	Target *Star

	// Time when reaching target
	TargetTime float64

	// Multiple of speed of light in vacuum
	Speed float64

	// If exited by user action
	Exited bool
}

func NewHyperspace(start *V3, target *Star, wTime, tTime float64) *Hyperspace {
	if wTime >= tTime {
		panic("target time must be in the future")
	}
	travelSeconds := tTime - wTime

	log.Debug("debug", "d", target.Entity)
	dist := new(V3).Sub(start, S.Pos[target.Entity]).Magnitude()
	lightSeconds := ((dist * gridUnit) * aum) / speedOfLight
	speed := lightSeconds / travelSeconds
	return &Hyperspace{
		start,
		target,
		tTime,
		speed,
		false,
	}
}

//
// System interface
//
func (hd *Hyperdrive) Init() error {
	return nil
}

func (hd *Hyperdrive) Update(wTime, elapsed float64, rf *RefFrame) error {
	log.Debug("Hyperdrive.Update")

	for e, _ := range S.HotEnts[rf] {
		if S.Hyperspace[e] != nil {
			updateHyperdrive(wTime, elapsed, rf, e)
		}
	}

	return nil
}

func updateHyperdrive(wTime, elapsed float64, rf *RefFrame, e Id) {
	hs := S.Hyperspace[e]
	if hs.TargetTime < wTime {
		// exit hyperdrive at destination, even if user sent exit action
		// since last update

		// TODO: how to retrieve / instantiate target ref frame
		//to := &RefFrame{}
		//updateEntityRefFrame(e, rf)

		S.Orb[e] = hs.Target.DefaultOrbit()
		delete(S.Hyperspace, e)
	}

	if hs.Exited {
		panic("todo")
		// TODO: check if local frame exists
		// TODO: disallow hyperdrive through ref frames like star systems
		// TODO: galactic coll det/resp - interstellar clouds, etc
		//delete(S.HSC, e)
		//continue
	}
}

func (hd *Hyperdrive) IsHotPostUpdate(e Id) bool {
	return S.Hyperspace[e] != nil
}
*/
