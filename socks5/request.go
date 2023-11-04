package socks5

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	RepSuccess            byte = 0x00
	RepSocks5Failure      byte = 0x01
	RepConnNotAllowed     byte = 0x02
	RepNetworkUnreachable byte = 0x03
	RepHostUnreachable    byte = 0x04
	RepConnRefused        byte = 0x05
	RepTTLExpired         byte = 0x06
	RepCmdNotSupported    byte = 0x07
	RepAddrNotSupported   byte = 0x08
	RepUnknown            byte = 0xff
)

func (s *Socks5) socks5auth(ctx *Context) error {
	// | VER | NMETHODS | METHODS  |
	// |  1  |    1     | 1 to 255 |
	// VER: 0x05 SOCKS5 version protocol.
	// NMETHODS: contains the number of method identifier that
	//           appear in the METHODS field.
	// METHODS: X'00' NO AUTHENTICATION REQUIRED
	header := readBytes(ctx.Conn, 2)
	if header[0] != 0x05 {
		return fmt.Errorf("protocol error, %v", header[0])
	}
	readBytes(ctx.Conn, int(header[1]))

	// | VER | METHOD |
	// |  1  |   1    |
	// METHOD: The SOCKS5 selects from one of the methods given in
	//         METHODS, and sends a METHOD selection message.
	writeBytes(ctx.Conn, []byte{0x05, 0x00})

	//
	return nil
}

func (s *Socks5) socks5reply(conn net.Conn, repStatus byte) {
	// | VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// |  1  |  1  | X'00' |  1   | Variable |    2     |
	writeBytes(conn, []byte{0x05, repStatus, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	)
}

func (s *Socks5) socks5request(ctx *Context) error {
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// VER: 0x05 SOCKS5
	// CMD: 0x01 CONNECT 0x02 BIND 0x03 UDP ASSOCIATE
	// ATYP: 0x01 IPV4 0x03 DOMAIN 0x04 IPV6
	header := readBytes(ctx.Conn, 4)

	// version
	if header[0] != 0x05 {
		return fmt.Errorf("protocol error, %v", header[0])
	}

	// address
	var addr string
	addrType := header[3]
	if addrType == 0x01 {
		addr = net.IP(readBytes(ctx.Conn, 4)).String()
	} else if addrType == 0x03 {
		size := int(readBytes(ctx.Conn, 1)[0])
		addr = string(readBytes(ctx.Conn, size))
	} else if addrType == 0x04 {
		addr = net.IP(readBytes(ctx.Conn, 16)).String()
	} else {
		return fmt.Errorf("addr not supported, %v", addrType)
	}

	// port
	bPort := readBytes(ctx.Conn, 2)
	port := strconv.Itoa(
		int(binary.BigEndian.Uint16(bPort)),
	)
	ctx.Addr = fmt.Sprintf("%v:%v", addr, port)

	// command
	if header[1] == 0x01 {
		ctx.Cmd = "connect"
	} else if header[1] == 0x03 {
		ctx.Cmd = "associate"
	} else {
		return fmt.Errorf("command not supported, %v", header[1])
	}

	// callback
	if s.requestCallback != nil {
		return s.requestCallback(ctx)
	}

	//
	return nil
}

func (s *Socks5) socks5stream(ctx *Context) error {
	if s.streamCallback == nil {
		return nil
	}
	if err := s.streamCallback(ctx); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}

func (s *Socks5) handleConnect(ctx *Context) {
	defer ctx.Conn.Close()

	// panic
	defer func() {
		if err := recover(); err != nil {
			if s.errorCallback != nil {
				s.errorCallback(ctx, err.(error))
			}
		}
	}()

	// init
	if s.initCallback != nil {
		if err := s.initCallback(ctx); err != nil {
			if s.errorCallback != nil {
				s.errorCallback(ctx, err)
			}
			return
		}
	}
	if s.closeCallback != nil {
		defer func() {
			if err := s.closeCallback(ctx); err != nil {
				if s.errorCallback != nil {
					s.errorCallback(ctx, err)
				}
			}
		}()
	}

	// auth
	if err := s.socks5auth(ctx); err != nil {
		if s.errorCallback != nil {
			s.errorCallback(ctx, err)
		}
		return
	}

	// request
	if err := s.socks5request(ctx); err != nil {
		if s.errorCallback != nil {
			s.errorCallback(ctx, err)
		}
		s.socks5reply(ctx.Conn, RepSocks5Failure)
		return
	}
	s.socks5reply(ctx.Conn, RepSuccess)

	// stream
	if err := s.socks5stream(ctx); err != nil {
		if s.errorCallback != nil {
			s.errorCallback(ctx, err)
		}
		return
	}
}
