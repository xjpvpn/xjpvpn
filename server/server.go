package server

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

const (
	BufferSize = 4096
)

type Server struct {
	keyPlain []byte
}

func New(keyPlain string) *Server {
	return &Server{keyPlain: []byte(keyPlain)}
}

func (s *Server) Listen(serverAddr string) {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.BasicAuth(
		func(user, password string, ctx echo.Context) (bool, error) {
			key := []byte(fmt.Sprintf("%s:%s", user, password))
			if subtle.ConstantTimeCompare(key, s.keyPlain) == 1 {
				return true, nil
			}
			return false, nil
		}),
	)

	// connect
	router.GET("/connect", func(c echo.Context) error {
		base64 := base64.RawURLEncoding

		// addr
		addrBase64 := c.QueryParam("addr")
		if addrBase64 == "" {
			return echo.ErrBadRequest
		}
		addr, err := base64.DecodeString(addrBase64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		// connect
		conn, err := net.DialTimeout("tcp", string(addr), time.Minute)
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err)
		}
		defer conn.Close()

		// websocket
		websocket.Handler(func(ws *websocket.Conn) {
			wg := &sync.WaitGroup{}

			// push
			wg.Add(1)
			pusher := func() {
				defer wg.Done()
				defer ws.Close()
				buf := make([]byte, BufferSize)
				for done := false; !done; {
					n, err := conn.Read(buf)
					if err != nil {
						done = true
					}
					if _, err := ws.Write(buf[:n]); err != nil {
						done = true
					}
				}
			}
			go pusher()

			// pull
			wg.Add(1)
			puller := func() {
				defer wg.Done()
				defer ws.Close()
				buf := make([]byte, BufferSize)
				for done := false; !done; {
					n, err := ws.Read(buf)
					if err != nil {
						done = true
					}
					if _, err := conn.Write(buf[:n]); err != nil {
						done = true
					}
				}
			}
			go puller()

			//
			wg.Wait()
		}).ServeHTTP(c.Response(), c.Request())

		//
		return err
	})

	//
	if err := router.Start(serverAddr); err != nil {
		slog.Error("[SERVER]", "error", err)
	}
}
