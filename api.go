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
	"time"
)

const (
	toClientChanBufferSize = 10
	apiSubExpiry           = 4 * time.Second

	// Mandatory API fields:
	keyCallType = "callType"
	keyParams   = "params"

	// The API has three types of calls:
	// 1. actions: see action.go
	// 2. getState: read-only state getters
	// 3. subState: subscribe to future state deltas
	valueAction   = "action"
	valueGetState = "getState"
	valueSubState = "subState"

	// actions
	keyActionName = "actionName"
	valueRotate   = "rotate"
	valueThrust   = "thrust"

	// getState
	keyStateType = "stateType"
	valueEnvFull = "envFull"

	// Common parameters
	keyEntity   = "entity"
	keyDuration = "duration"
	keyForce    = "force"
	keyX        = "x"
	keyY        = "y"
	keyZ        = "z"
)

type Executer interface {
	Execute() interface{}
}

type APIExec struct {
	Ex       Executer
	RespChan chan interface{}
}

//
// API Actions
//
func APIActionRotate(e Id, t *V3, d float64) error {
	a := &ActionRotate{entity: e, t: t, duration: d}
	x := call(GE.actionChan, a)
	switch resp := x.(type) {
	case error:
		return resp
	case nil:
		return nil
	}
	return nil
}

func APIActionThrust(e Id, t, d float64) {
	a := &ActionEngineThrust{entity: e, thrust: t, duration: d}
	cast(GE.actionChan, a)
}

//
// API State Getters
//
func APIGetGalaxy() []Sector {
	r := AllSectors{}
	x := call(GE.getStateChan, r)
	return x.([]Sector)
}

// call makes a synchronous call to the game engine by sending a request
// and waiting until a reply arrives.
func call(engineChan chan APIExec, ex Executer) interface{} {
	respChan := make(chan interface{})
	engineChan <- APIExec{Ex: ex, RespChan: respChan}
	resp := <-respChan
	return resp
}

// cast sends an asynchronous request to the game engine.
func cast(engineChan chan APIExec, ex Executer) {
	engineChan <- APIExec{Ex: ex}
}

//
// TODO: API Call Types
//
type AllSectors struct {
}

func (as AllSectors) Execute() interface{} {
	// Copy all reference type to avoid any links to engine state
	sectors := make(map[string]*Sector)
	/*
		for k, v := range S.Sectors {
			// TODO: add .Clone() to all relevant state types
			c := &V3{v.Corner.X, v.Corner.Y, v.Corner.Z}
			sss := make([]*StarSystem, len(v.StarSystems))
			for i, ss := range v.StarSystems {
				sss
			}
			s := Sector{Corner: c, Mapped: v.Mapped, StarSystems: ss}
			sectors[k] = s
		}
	*/
	return sectors
}

type SpatialState struct {
	Entity Id

	OE *OE

	Pos *V3
	Vel *V3

	Ori *Q
	Rot *V3
}

func (ss *SpatialState) Execute() (interface{}, error) {
	e := ss.Entity
	ss.OE = S.Orb[e]
	ss.Pos = S.Pos[e]
	ss.Vel = S.Vel[e]
	ss.Ori = S.Ori[e]
	ss.Rot = S.Rot[e].R // TODO: include body/world transform?

	return ss, nil
}

/*
func API(fields map[string]interface{}, respChan chan<- []byte) error {
	// Decode mandatory fields
	callType, ok := fields[keyCallType].(string)
	if !ok {
		return invalid(keyCallType, fields[keyCallType])
	}

	params, ok := fields[keyParams].(map[string]interface{})
	if !ok {
		return invalid(keyParams, fields[keyParams])
	}

	switch callType {
	case valueAction:
		return handleAction(params)
	case valueGetState:
		return handleGetState(params, respChan)
	case valueSubState:
		return handleSubState(params, respChan)
	default:
		return invalid(keyCallType, callType)
	}
}

func handleAction(params map[string]interface{}) error {
	eStr, ok := params[keyEntity].(string)
	if !ok {
		return invalid(keyEntity, params[keyEntity])
	}
	e, err := strconv.ParseUint(eStr, 10, 64)
	if err != nil {
		return err
	}

	actionName, ok := params[keyActionName].(string)
	if !ok {
		return invalid(actionName, params[keyActionName])
	}

	// TODO: make duration optional
	duration, ok := params[keyDuration].(float64)
	if !ok {
		return invalid(keyDuration, params[keyDuration])
	}

	// TODO: safe field error handling
	var ar Action
	switch actionName {
	case valueRotate:
		torque := params[keyForce].(map[string]interface{})
		x := torque[keyX].(float64)
		y := torque[keyY].(float64)
		z := torque[keyZ].(float64)
		ar = &ActionRotate{Id(e), &V3{x, y, z}, duration}
	case valueThrust:
		f := params[keyForce].(float64)
		ar = &ActionEngineThrust{Id(e), f, duration}
	default:
		return invalid(keyActionName, params[keyActionName])
	}
	GE.actionChan <- ar
	return nil
}

func handleGetState(params map[string]interface{}, respChan chan<- []byte) error {
	return nil
}

func handleSubState(params map[string]interface{}, respChan chan<- []byte) error {
	eStr, ok := params[keyEntity].(string)
	if !ok {
		return invalid(keyEntity, params[keyEntity])
	}
	e, err := strconv.ParseUint(eStr, 10, 64)
	if err != nil {
		return err
	}
	// TODO: send to client subState request
	return nil
}

func invalid(k string, v interface{}) error {
	return fmt.Errorf("invalid type or value for key: %s value: %s", k, reflect.TypeOf(v).String())
}

type EnvFull struct {
	RefFrame *RefFrame
	// TODO: starsystem
}

func sendEnvFull(e Id, respChan chan<- []byte) {
	env := EnvFull{
		RefFrame: S.EntFrames[e],
	}

	b, err := json.Marshal(env)
	if err != nil {
		panic(err)
	}

	respChan <- b
}
*/

/*
// APISubInit must be called before API.
// It returns a data channel and a keepalive channel.
//
// The data channel is used for data that is sent to the client.
// The game engine writes all API responses and subscription events
// to the data channel.
//
// The game engine terminates the data channel if
// 1. It has not sent anything to the data channel and not received
//    a signal on the keepalive channel within apiSubExpiry time.
// 2. OR any API params from the client are invalid
// 3. OR an error is encountered during processing a call.
func APISubInit() (<-chan []byte, chan<- bool, error) {
	respChan := make(chan []byte, respChanBufferSize)
	keepAliveChan := make(chan bool, 1)
	keepAliveChan <- true

	closeChans := func() {
		log.Info("API closing chans")
		close(respChan)
		close(keepAliveChan)
	}

	go func() {
		for {
			select {
			case ok := <-keepAliveChan:
				log.Info("API", "keepAliveChan", ok)
				if ok {
					time.Sleep(apiSubExpiry)
				} else {
					closeChans()
					return
				}
			default:
				log.Info("API keepAliveChan expired")
				closeChans()
				return
			}
		}
	}()

	return respChan, keepAliveChan, nil
}
*/

/*
	e, err := strconv.ParseUint(fields[keyEntity].(string), 10, 64)
	if err != nil {
		return err
	}
	params := fields[keyParams].(map[string]interface{})
	duration := params[keyDuration].(float64)

	var ar Action
	if fields[keyAction] == keyRotate {
		torque := params[keyForce].(map[string]interface{})
		x := torque[keyX].(float64)
		y := torque[keyY].(float64)
		z := torque[keyZ].(float64)
		ar = &ActionRotate{Id(e), &V3{x, y, z}, duration}
	} else {
		f := params[keyForce].(float64)
		ar = &ActionEngineThrust{Id(e), f, duration}
	}
	GE.actionChan <- ar

	return nil
}
*/
