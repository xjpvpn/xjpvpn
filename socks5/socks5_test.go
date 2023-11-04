package socks5_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
	"tunnel/socks5"
)

const (
	ProxyAddr = "localhost:11080"
)

func TestBasic(t *testing.T) {
	// socks5
	server := socks5.New()
	server.OnStream(func(ctx *socks5.Context) error {
		buffer := make([]byte, 1024)
		n, err := ctx.Conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		header := string(buffer[:n])
		if !strings.Contains(header, "GET / HTTP/1.1") {
			panic("unexpected behavior")
		}
		if _, err := ctx.Conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" + "Content-Length: 0\r\n\r\n",
		)); err != nil {
			panic(err)
		}
		return io.EOF
	})
	go server.Listen(ProxyAddr)
	time.Sleep(time.Second)

	// client
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(
					fmt.Sprintf("socks5://%v", ProxyAddr),
				)
			},
		},
	}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("unexpected behavior")
	}
}
