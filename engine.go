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

/*  NOTES

    NEW ARCH:
    1. Entity Component System (Ids, Data, Functions...:D)
    2. Systems write on message bus subscribed to by other systems.
    3. Order of updates:
    3.1. Physics: dynamic to env first, as dynamic to dynamic has last say
    3.2. Physics before "game" logic - game has final say over physics as its
         a layer on top of physics

    X. randomize order of entities update() call - to avoid
      potential exploit of being updated earlier than others

    X. the game engine's time logic: this is the core
        - logical frames?
        - one block/frame every 1s?

*/

const (
	loopTarget        = 1000 * time.Millisecond
	maxActionsPerLoop = 10
)

type GameEngine struct {
	systems    []System
	actionChan chan Action
}

var GE *GameEngine

func (ge *GameEngine) Loop() error {
	var err error
	var elapsed time.Duration
	var start, last, t time.Time

	var j []byte

	start = time.Now()
	last = start

	debug := 0
	for err == nil {
		debug++
		t = time.Now()
		elapsed = t.Sub(last)
		//log.Debug("engine.Loop", "c", debug, "run", time.Now().Sub(start))

		if elapsed < loopTarget {
			time.Sleep(loopTarget - elapsed)
			last = time.Now()
			elapsed = t.Sub(last)
		} else {
			last = t
		}

		err = ge.update(elapsed.Seconds())
		if err != nil {
			break
		}

		err = ge.handleUserActions()
		if err != nil {
			break
		}

		// TODO: per client
		j, err = json.Marshal(S)
		if err != nil {
			break
		}
		S.MB.Post(j)
	}

	log.Info("engine.Loop", "err", err)
	return err
}

// TODO: update only hot ents/frames
// TODO: derive the update order for ref frames and ents from random beacon
func (ge *GameEngine) update(elapsed float64) error {
	for rf, entMap := range S.EntsInFrames {
		for _, sys := range ge.systems {
			err := sys.Update(elapsed, rf, entMap)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *GameEngine) handleUserActions() error {
	var actions []Action
	select {
	default:
		return nil
	case action := <-e.actionChan:
		actions = make([]Action, 0)
		actions = append(actions, action)
		for i := 1; i < maxActionsPerLoop; i++ {
			select {
			default:
				for _, a := range actions {
					err := a.Execute()
					if err != nil {
						return err
					}
				}
				return nil
			case a := <-e.actionChan:
				actions = append(actions, a)
			}
		}
	}

	//return errors.New("wtf")
	return nil
}
