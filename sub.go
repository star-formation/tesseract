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
	"encoding/json"
)

type EntitySub struct {
	entity        Id
	dataChan      chan []byte
	keepAliveChan chan bool
}

type EntitySubData struct {
	OE *OE

	Pos *V3
	Vel *V3

	Ori *Q
	Rot *V3
}

func (es *EntitySub) Update() {
	e := es.entity
	data := EntitySubData{
		OE:  S.Orb[e],
		Pos: S.Pos[e],
		Vel: S.Vel[e],
		Ori: S.Ori[e],
		Rot: S.Rot[e].R, // TODO: include body/world transform?
	}

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	es.dataChan <- b
}
