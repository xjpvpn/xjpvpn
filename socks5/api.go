package socks5

import (
	"net"
)

type Context struct {
	Cmd  string
	Addr string

	// socket
	Conn net.Conn

	// managed data
	Data map[string]interface{}
}

type Socks5 struct {

	// callbacks
	initCallback    func(*Context) error
	requestCallback func(*Context) error
	streamCallback  func(*Context) error
	closeCallback   func(*Context) error
	errorCallback   func(*Context, error)
}

func New() *Socks5 {
	return &Socks5{}
}

func (s *Socks5) OnInit(
	handler func(ctx *Context) error) {
	s.initCallback = handler
}

func (s *Socks5) OnRequest(
	handler func(ctx *Context) error) {
	s.requestCallback = handler
}

func (s *Socks5) OnStream(
	handler func(ctx *Context) error) {
	s.streamCallback = handler
}

func (s *Socks5) OnClose(
	handler func(ctx *Context) error) {
	s.closeCallback = handler
}

func (s *Socks5) OnError(
	handler func(ctx *Context, err error)) {
	s.errorCallback = handler
}

func (s *Socks5) Listen(localAddr string) {
	// listen
	ln, err := net.Listen("tcp", localAddr)
	if err != nil {
		if s.errorCallback != nil {
			s.errorCallback(nil, err)
		}
		return
	}

	// handle
	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.errorCallback != nil {
				s.errorCallback(nil, err)
			}
			continue
		}
		go s.handleConnect(&Context{
			Conn: conn, Data: make(
				map[string]interface{},
			),
		})
	}
}
