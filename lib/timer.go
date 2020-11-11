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
    "sort"
    "time"
)

// Timers are actions scheduled for future execution.

// The timer component is attachable to any entity and used for things like
// delayed-effect weapons, manufacturing processes and skill training.
type Timer struct {
    Id Id
    MS time.Time
    Action *Action
}

// Alongside the map of entity ids to Timers, we also maintain
// this list (slice) of Timers Sorted in chronological order.
// The engine loop uses this to only traverse expired Timers.
type TimerList []*Timer
// sort.Interface
func (tl TimerList) Len() int           { return len(tl) }
func (tl TimerList) Swap(i, j int)      { tl[i], tl[j] = tl[j], tl[i] }
func (tl TimerList) Less(i, j int) bool { return tl[i].MS.Before(tl[j].MS) }

type TimerComponent struct {
    Timers map[Id]*Timer
    Sorted TimerList
}

func (tc *TimerComponent) Init() error {
    return nil
}

func (tc *TimerComponent) AddEntityData(id Id, t *Timer) {
    tc.Timers[id] = t
    tc.Sorted = append(tc.Sorted, t)
    sort.Stable(tc.Sorted)
}

func (tc *TimerComponent) GetEntityData(id Id) (*Timer) {
    return tc.Timers[id]
}

func (tc *TimerComponent) RemoveEntity(id Id) {
    delete(tc.Timers, id)

    f := func(i int) bool {return tc.Sorted[i].Id == id}
    i := sort.Search(len(tc.Sorted), f)
    tc.Sorted = append(tc.Sorted[:i], tc.Sorted[i+1:]...)
}
