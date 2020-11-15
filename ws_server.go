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
	//"errors"
	//"encoding/binary"
	"net/http"
	
	"github.com/gorilla/websocket"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// note: localhost:8081 as addr string works with raw TCP but not with HTTP
	host = ":8081"
)

var upgrader = websocket.Upgrader{}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("func handler: ", "r", r)
	// TODO: secure origin check
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade err:", "err", err)
		return
	}
	defer c.Close()

	sps := websocket.Subprotocols(r)
	if len(sps) != 1 || sps[0] != "client0.argonavis.io" {
		log.Error("Unsupported WebSocket Subprotocol: ", "subs", sps)
		WriteControlClose(c, websocket.CloseProtocolError, "unsupported subprotocol")
		return
	}

	// setup MessageBus sub to engine loop
	go func() {
		ch := S.MsgBus.Subscribe()
		for {
			stateJSON := <-ch
			err = c.WriteMessage(websocket.BinaryMessage, stateJSON)
			if err != nil {
				log.Error("write err:", "err", err)
				break
			}
		}
	}()

	for {
		log.Info("waiting on c.ReadMessage: ")
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Error("read:", "err", err)
			break
		}
		log.Info("recv: ", "msg", msg)
		
		// handle action
		err = HandleMsg(msg)
		if err != nil {
			log.Error("HandleMsg", "err", err)
			log.Info("Closing websocket conn")
			WriteControlClose(c, websocket.CloseInternalServerErr, err.Error())
			return
		}
	}
}

func HandleMsg(msg []byte) error {
	return nil
}

func WriteControlClose(c *websocket.Conn, closeCode int, str string) error {
	msg := websocket.FormatCloseMessage(closeCode, str)
	return c.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second*5))
}

func StartWebSocket() {
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Error("http.ListenAndServe", "err", err)
		return
	}
	
	log.Info("WebSocket Server Started", "host", host)
	return
}
