package sockets

import (
	"context"
	"fmt"
	"msg/models"
	"net/url"

	"github.com/gorilla/websocket"
)

type ClientMessage struct {
	Msg   *models.StorageMessage
	Error error
}

type ClientSocket struct {
	conn *websocket.Conn
	user string

	callback func(*ClientMessage)

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClientSocket(ctx context.Context, user string, u url.URL, cb func(*ClientMessage)) (*ClientSocket, error) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	clientSocket := &ClientSocket{
		conn:     c,
		user:     user,
		callback: cb,
		ctx:      ctx,
		cancel:   cancel,
	}
	go func() {
		for {
			msg := &models.StorageMessage{}
			err := c.ReadJSON(msg)

			select {
			case <-ctx.Done():
				return
			default:
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					cancel()
					return
				}
				cb(&ClientMessage{Msg: msg, Error: err})
			}
		}
	}()

	return clientSocket, nil
}

func (cs *ClientSocket) SendMessage(msg string) error {
	socketMesasge := models.CreateMessageSocket(cs.user, msg)

	select {
	case <-cs.ctx.Done():
		return fmt.Errorf("connection closed")
	default:
	}

	err := cs.conn.WriteJSON(socketMesasge)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ClientSocket) Stop() error {
	cs.SendMessage("Has left the channel")

	err := cs.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	cs.cancel()

	return nil
}
