package sockets

import (
	"context"
	"msg/models"
	"net/http"

	"github.com/gorilla/websocket"
)

type ServerMessage struct {
	Msg   *models.StorageMessage
	Error error
}

type ServerSocket struct {
	clients map[*websocket.Conn]struct{}

	server   http.Server
	upgrader websocket.Upgrader

	storage  models.Storage
	callback func(*ServerMessage)
}

func NewServerSocket(addr string, storage models.Storage, upgrader *websocket.Upgrader, callback func(*ServerMessage)) *ServerSocket {
	if upgrader == nil {
		upgrader = &websocket.Upgrader{}
	}

	serverSocket := &ServerSocket{
		clients: make(map[*websocket.Conn]struct{}),

		server:   http.Server{Addr: addr},
		upgrader: *upgrader,
		storage:  storage,

		callback: callback,
	}

	return serverSocket
}

func (socket *ServerSocket) Start() error {
	socket.server.Handler = socket
	err := socket.server.ListenAndServe()

	return err
}

func (socket *ServerSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection, err := socket.upgrader.Upgrade(w, r, nil)
	if err != nil {
		socket.callback(&ServerMessage{nil, err})
		return
	}
	defer connection.Close()

	socket.clients[connection] = struct{}{}
	defer delete(socket.clients, connection)

	socket.onConnection(connection)

	for {
		msg := &models.SocketMessage{}
		err := connection.ReadJSON(msg)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			socket.callback(&ServerMessage{nil, err})
			break
		}

		enhancedMsg := models.CreateMessageDB(msg)
		err = socket.storage.Append(enhancedMsg)
		if err != nil {
			socket.callback(&ServerMessage{nil, err})
			break
		}

		socket.callback(&ServerMessage{enhancedMsg, nil})

		go socket.writeClients(enhancedMsg)
	}
}

func (socket *ServerSocket) Stop() error {
	err := socket.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (socket *ServerSocket) onConnection(conn *websocket.Conn) {
	latestMsgs, err := socket.storage.GetLatest(10)
	if err != nil {
		socket.callback(&ServerMessage{nil, err})
		return
	}
	for _, msg := range latestMsgs {
		err = conn.WriteJSON(msg)

		if err != nil {
			socket.callback(&ServerMessage{nil, err})
			return
		}
	}
}

func (socket *ServerSocket) writeClients(msg *models.StorageMessage) {
	for conn := range socket.clients {
		err := conn.WriteJSON(msg)
		if err != nil {
			socket.callback(&ServerMessage{nil, err})
			return
		}
	}
}
