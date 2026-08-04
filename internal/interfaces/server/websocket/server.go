package websocket

import (
	"fmt"
	"net/http"
	"time"

	coderwebsocket "github.com/coder/websocket"
)

type ServerConfig struct {
	OriginPatterns []string
	WriteTimeout   time.Duration
}

type Server struct {
	originPatterns []string
	writeTimeout   time.Duration
}

func NewServer(config ServerConfig) *Server {
	writeTimeout := config.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = defaultWriteTimeout
	}

	return &Server{
		originPatterns: config.OriginPatterns,
		writeTimeout:   writeTimeout,
	}
}

func (server *Server) Accept(writer http.ResponseWriter, request *http.Request) (*Connection, error) {
	rawConnection, err := coderwebsocket.Accept(
		writer,
		request,
		&coderwebsocket.AcceptOptions{
			OriginPatterns: server.originPatterns,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"accept websocket connection: %w",
			err,
		)
	}

	connection, err := NewConnection(rawConnection, server.writeTimeout)
	if err != nil {
		_ = rawConnection.Close(
			coderwebsocket.StatusInternalError,
			"failed to initialize connection",
		)

		return nil, fmt.Errorf(
			"initialize websocket connection: %w",
			err,
		)
	}

	return connection, nil
}
