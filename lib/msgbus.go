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
package lib

// MessageBus is a channel-based message / event bus where systems
// post messages delivered to all subscribers.
type MessageBus struct {
	Channels []chan<- []byte
}

func (mb *MessageBus) Subscribe() <-chan []byte {
	c := make(chan []byte, 10)
	mb.Channels = append(mb.Channels, c)
	return c
}

func (mb *MessageBus) Post(msg []byte) {
	for _, c := range mb.Channels {
		c <- msg
	}
}
