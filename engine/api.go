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
	"time"

	"github.com/ethereum/go-ethereum/log"
)

const (
	dataChanBufferSize = 10
	subExpiry          = 4 * time.Second
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

func NewEntitySub(e Id) (<-chan []byte, chan<- bool) {
	dataChan := make(chan []byte, dataChanBufferSize)
	keepAliveChan := make(chan bool, 1)
	es := &EntitySub{e, dataChan, keepAliveChan}

	GE.subChan <- es
	S.EntitySubs[es] = struct{}{}

	unsub := func() {
		close(dataChan)
		close(keepAliveChan)
		S.EntitySubsCloseChan <- es
	}

	keepAliveChan <- true

	go func() {
		for {
			select {
			case ok := <-keepAliveChan:
				if ok {
					time.Sleep(subExpiry)
				} else {
					log.Info("Unsubscribe: keepAlive=false")
					unsub()
					return
				}
			default:
				log.Info("Unsubscribe: expired")
				unsub()
				return
			}
		}
	}()

	return dataChan, keepAliveChan
}
