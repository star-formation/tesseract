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
    "errors"
    "encoding/binary"
    //"github.com/ethereum/go-ethereum/log"
)

// Actions are authenticated attempts to modify the game state.
// Most actions originate from users, where we consider them
// authenticated post user account signature verification.
// Actions can also originate from the game engine itself, to e.g. trigger
// regular story events, in which case they are always authenticated.

type Action interface {
    Execute() error
}

// Game action codes.
const (
    // TODO: for testing, this returns the entire global game state.
    GetGlobalState = 601
)

// The input to this function is raw bytes handed from other layers, e.g.
// the binary payload of a WebSocket message.  As such these bytes have
// not yet been validated, and may be malicious.
// TODO: refactor and document security assumptions and input validation
// in diff layers.
func HandleMsg(msg []byte) error {
    // TODO: limit actions per account over time
    // TODO: make lookup for codes
    if len(msg) < 2 {
        return errors.New("invalid command")
    } 
    code := binary.BigEndian.Uint16(msg[0:1])
    if code != GetGlobalState {
        return errors.New("invalid action code")
    }

    return nil
}

