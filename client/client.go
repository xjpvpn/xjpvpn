package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"tunnel/socks5"

	"golang.org/x/net/websocket"
)

const (
	BufferSize = 4096
)

type Client struct {
	keyPlain []byte
}

func New(keyPlain string) *Client {
	return &Client{keyPlain: []byte(keyPlain)}
}

func (c *Client) Listen(socks5Addr, serverAddr string) {
	server := socks5.New()
	server.OnError(func(ctx *socks5.Context, err error) {
		slog.Error("[CLIENT]", "error", err, "addr", ctx.Addr)
	})
	server.OnRequest(func(ctx *socks5.Context) error {
		base64 := base64.RawURLEncoding

		// websocket
		path := fmt.Sprintf(
			"ws://%s/%s?addr=%s", serverAddr, ctx.Cmd,
			base64.EncodeToString([]byte(ctx.Addr)),
		)
		origin := fmt.Sprintf("http://%s/", serverAddr)
		config, err := websocket.NewConfig(path, origin)
		if err != nil {
			panic(err)
		}
		config.Header.Set(
			"Authorization", fmt.Sprintf(
				"basic %s", base64.EncodeToString(c.keyPlain),
			),
		)
		ws, err := websocket.DialConfig(config)
		if err != nil {
			return err
		}
		ctx.Data["ws"] = ws

		//
		return nil
	})
	server.OnStream(func(ctx *socks5.Context) error {
		wg := &sync.WaitGroup{}
		ws := ctx.Data["ws"].(*websocket.Conn)

		// push
		wg.Add(1)
		pusher := func() {
			defer wg.Done()
			defer ws.Close()
			msg := make([]byte, BufferSize)
			for done := false; !done; {
				n, err := ctx.Conn.Read(msg)
				if err != nil {
					done = true
				}
				if _, err := ws.Write(msg[:n]); err != nil {
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
			msg := make([]byte, BufferSize)
			for done := false; !done; {
				n, err := ws.Read(msg)
				if err != nil {
					done = true
				}
				if _, err := ctx.Conn.Write(msg[:n]); err != nil {
					done = true
				}
			}
		}
		go puller()

		//
		wg.Wait()
		return io.EOF
	})
	server.Listen(socks5Addr)
}
