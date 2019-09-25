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
	"errors"
	"time"

	xrand "golang.org/x/exp/rand"

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

// The Engine manages entities, components and systems and handles
// much of the core operations of the game engine.
//
// It is intentionally somewhat of a monolith / god object during development
// to avoid premature abstractions.
//
// Long term we want to abstract out functionality into separate modules for
// better SoC, but the engine will likely retain a considerable amount of
// core functions.
var GE *GameEngine

type GameEngine struct {
	systems    []System
	actionChan chan Action
}

func (e *GameEngine) Loop() error {
	var err error
	var elapsed time.Duration
	var start, last, now time.Time

	var r *xrand.Rand
	var j []byte
	// debug
	debug := 0

	start = time.Now()
	last = start
	for err == nil {
		debug++
		now = time.Now()
		elapsed = now.Sub(last)
		//log.Debug("engine.Loop", "c", debug, "run", time.Now().Sub(start))

		if elapsed < loopTarget {
			time.Sleep(loopTarget - elapsed)
			elapsed = loopTarget
			last = time.Now()
		} else {
			last = now
		}

		r, err = NewRand()
		if err != nil {
			break
		}

		err = e.updateHotEnts(r, elapsed.Seconds())
		if err != nil {
			break
		}

		//err = e.processTimers(r)
		//if err != nil {
		//    break
		//}

		err = e.handleUserActions(r)
		if err != nil {
			break
		}

		j, err = json.Marshal(S)
		if err != nil {
			break
		}
		S.MB.Post(j)
	}

	// TODO: error handling
	log.Info("engine.Loop", "err", err)
	return err
}

func (e *GameEngine) updateHotEnts(r *xrand.Rand, elapsed float64) error {
	for _, sys := range e.systems {
		err := sys.Update(elapsed)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *GameEngine) handleUserActions(rand *xrand.Rand) error {
	var actions []Action
	select {
	default:
		return nil
	case action := <-e.actionChan:
		log.Debug("handleUserActions", "e.actionChan", e.actionChan, "action", action)
		//log.Debug("handleUserActions", "a", a)
		actions = make([]Action, 0)
		actions = append(actions, action)
		for i := 0; i < maxActionsPerLoop; i++ {
			select {
			default:
				for _, a := range actions {
					log.Debug("handleUserActions", "a", a)
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
	return errors.New("wtf")
}
