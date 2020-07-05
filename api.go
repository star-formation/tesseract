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
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"
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

// DevAPISub must be called before DevAPISend.
// It returns a data channel and a keepalive channel.
//
// The data channel is used for data that is sent to the client.
// The game engine writes all API responses and subscription events
// to the data channel.
//
// The game engine terminates the API sub and data channel if
// 1. It has not sent anything to the data channel and not received
//    a signal on the keepalive channel within apiSubExpiry time.
// 2. OR any API params from the client are invalid
// 3. OR an error is encountered during processing a call.
func DevAPISub() (<-chan []byte, chan<- bool, error) {
	toClientChan := make(chan []byte, toClientChanBufferSize)
	keepAliveChan := make(chan bool, 1)
	keepAliveChan <- true

	unsub := func() {
		log.Info("APISub unsub")
		close(toClientChan)
		close(keepAliveChan)
	}

	go func() {
		for {
			select {
			case ok := <-keepAliveChan:
				log.Info("APISub", "keepAliveChan", ok)
				if ok {
					time.Sleep(apiSubExpiry)
				} else {
					unsub()
					return
				}
			default:
				log.Info("APISub keepAliveChan expired")
				unsub()
				return
			}
		}
	}()

	return toClientChan, keepAliveChan, nil
}

// DevAPISend calls the game engine API.
// fields is a map of string attributes to values of any type.
// toClientChan is a data channel to which will be sent
// all API responses and all API subscription event data.
// TODO: make all type assertions safe
func DevAPISend(fields map[string]interface{}, toClientChan chan<- []byte) error {
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
		return handleGetState(params, toClientChan)
	case valueSubState:
		return handleSubState(params, toClientChan)
	default:
		return invalid(keyCallType, callType)
	}
}

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

func handleGetState(params map[string]interface{}, toClientChan chan<- []byte) error {
	return nil
}

func handleSubState(params map[string]interface{}, toClientChan chan<- []byte) error {
	// TODO
	return nil
}

func invalid(k string, v interface{}) error {
	return fmt.Errorf("invalid type or value for key: %s value: %s", k, reflect.TypeOf(v).String())
}

type EnvFull struct {
	RefFrame *RefFrame
	// TODO: starsystem
}

func sendEnvFull(e Id, toClientChan chan<- []byte) {
	env := EnvFull{
		RefFrame: S.EntFrames[e],
	}

	b, err := json.Marshal(env)
	if err != nil {
		panic(err)
	}

	toClientChan <- b
}
