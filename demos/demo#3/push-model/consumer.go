/*
Copyright © 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	apiURL   = ":8080"
	endpoint = "/result"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("push model consumer: logger initialized")
}

func receiveHandler(writer http.ResponseWriter, request *http.Request) {
	// The Upgrader.Upgrade method upgrades the HTTP server
	// connection to the WebSocket protocol as described in the
	// WebSocket RFC. A summary of the process is this: The client
	// sends an HTTP request requesting that the server upgrade the
	// connection used for the HTTP request to the WebSocket
	// protocol. The server inspects the request and if all is good
	// the server sends an HTTP response agreeing to upgrade the
	// connection. From that point on, the client and server use
	// the WebSocket protocol over the network connection.
	connection, err := upgrader.Upgrade(writer, request, nil) // error ignored for sake of simplicity
	if err != nil {
		log.Error().Err(err).Msg("upgrader.Upgrade")
	}

	for {
		// try to read message from producer
		msgType, message, err := connection.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("connection.ReadMessage()")
			return
		}

		// print the message to the console
		fmt.Printf("received %d bytes from %s (message type %d)\n", len(message), connection.RemoteAddr(), msgType)

		// write response back to producer
		//if err = conn.WriteMessage(msgType, msg); err != nil {
		//	return
		//}
	}
}

func main() {
	http.HandleFunc(endpoint, receiveHandler)
	http.ListenAndServe(apiURL, nil)
}
