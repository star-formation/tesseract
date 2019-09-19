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

// The Engine manages entities, components and systems and handles
// much of the core operations of the game engine.
//
// It is intentionally somewhat of a monolith / god object during development
// to avoid premature abstractions.
//
// Long term we want to abstract out functionality into separate modules for
// better SoC, but the engine will likely retain a considerable amount of
// core functions.
type Engine struct {
	systems []System
	mainBus *MessageBus
	hot     *HotEnts
}

var loopTarget = 20 * time.Millisecond

func (e *Engine) Loop() error {
	var err error
	var elapsed time.Duration
	var start, last, now time.Time

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

		r, err := NewRand()
		if err != nil {
			break
		}

		err = e.updatehot(r, elapsed.Seconds())
		if err != nil {
			break
		}

		//err = e.processTimers(r)
		//if err != nil {
		//    break
		//}

		err = e.handleUserActions(r)

		j, err := json.Marshal(S)
		if err != nil {
			return err
		}
		S.MB.Post(j)
	}

	// TODO: error handling
	log.Info("engine.Loop", "err", err)
	return err
}

// TODO: error handling
func (e *Engine) updatehot(r *xrand.Rand, elapsed float64) error {
	framePerm := r.Perm(len(e.hot.Frames))
	for _, i := range framePerm {
		e.updatehotFrame(r, elapsed, e.hot.Frames[i])
	}

	return nil
}

func (e *Engine) updatehotFrame(r *xrand.Rand, elapsed float64, f *RefFrame) {
	r.Shuffle(len(*e.hot.In[f]), func(i, j int) {
		(*e.hot.In[f])[i], (*e.hot.In[f])[j] = (*e.hot.In[f])[j], (*e.hot.In[f])[i]
	})

	for _, sys := range e.systems {
		sys.Update(elapsed, f, e.hot.In[f])
	}
}

/*
func (e *Engine) processTimers(r *xrand.Rand) error {
    var err error
    toRemove := []Id{}
    now := time.Now()
    for _, t := e.timerComponent.Sorted {
        if t.Time.After(now) {
            return nil
        }

        toRemove = append(toRemove, t.Id)
        err = t.Action.Execute()
        if err != nil {
            break
        }
    }

    // TODO: make O(1) for t.Sorted as we're removing in-order
    for _, i := range toRemove {
        e.timerComponent.RemoveEntity(i)
    }

    return err
}
*/

func (e *Engine) handleUserActions(rand *xrand.Rand) error {
	return nil
}
