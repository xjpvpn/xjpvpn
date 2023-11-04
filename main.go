package main

import (
	"flag"
	"log"
	"log/slog"
	"strings"
	"tunnel/client"
	"tunnel/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var runTest bool
	var keyPlain string
	var localAddr string
	var serverAddr string
	flag.BoolVar(&runTest, "test", false, "run test")
	flag.StringVar(
		&keyPlain, "password", "", "password for authorize",
	)
	flag.StringVar(
		&localAddr, "local_addr", "", "local address used in local",
	)
	flag.StringVar(
		&serverAddr, "server_addr", "", "local address used in client",
	)
	flag.Parse()

	if runTest {
		keyPlain = "admin:123456"
		localAddr = "localhost:1081"
		serverAddr = "localhost:8081"
	}

	// validate
	if strings.Contains(localAddr, "0.0.0.0") {
		slog.Error("localAddr is not safe")
		return
	}
	if !strings.Contains(keyPlain, ":") {
		slog.Error("key should like user:password")
		return
	}

	//
	local := client.New(keyPlain)
	server := server.New(keyPlain)

	// start serve
	if runTest {
		if len(localAddr) > 0 && len(serverAddr) > 0 {
			go server.Listen(serverAddr)
			local.Listen(localAddr, serverAddr)
		} else {
			slog.Error("[MAIN]", "localAddr", localAddr, "serverAddr", serverAddr)
		}
	} else {
		if len(localAddr) == 0 && len(serverAddr) > 0 {
			server.Listen(serverAddr)
		} else if len(localAddr) > 0 && len(serverAddr) > 0 {
			local.Listen(localAddr, serverAddr)
		} else {
			slog.Error("[MAIN]", "localAddr", localAddr, "serverAddr", serverAddr)
		}
	}
}
