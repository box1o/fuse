package websocket

import (
	"context"
	"errors"
	"sync"
	"time"

	coderwebsocket "github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const defaultWriteTimeout = 10 * time.Second

var ErrNilConnection = errors.New("websocket connection is nil")

type Connection struct {
	connection   *coderwebsocket.Conn
	writeTimeout time.Duration
	writeMutex   sync.Mutex
	closeOnce    sync.Once
}

func NewConnection(connection *coderwebsocket.Conn, writeTimeout time.Duration) (*Connection, error) {
	if connection == nil {
		return nil, ErrNilConnection
	}

	if writeTimeout <= 0 {
		writeTimeout = defaultWriteTimeout
	}

	return &Connection{
		connection:   connection,
		writeTimeout: writeTimeout,
	}, nil
}

func (connection *Connection) WriteJSON(ctx context.Context, value any) error {
	connection.writeMutex.Lock()
	defer connection.writeMutex.Unlock()

	writeContext, cancel := context.WithTimeout(
		ctx,
		connection.writeTimeout,
	)
	defer cancel()

	return wsjson.Write(
		writeContext,
		connection.connection,
		value,
	)
}

func (connection *Connection) ReadJSON(ctx context.Context, target any) error {
	return wsjson.Read(
		ctx,
		connection.connection,
		target,
	)
}

func (connection *Connection) Close(status coderwebsocket.StatusCode, reason string) error {
	var closeError error

	connection.closeOnce.Do(func() {
		closeError = connection.connection.Close(status, reason)
	})

	return closeError
}

func (connection *Connection) CloseNormal() error {
	return connection.Close(
		coderwebsocket.StatusNormalClosure,
		"connection closed",
	)
}
