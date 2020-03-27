package main

import (
	"net"
	"time"

	logger "../utils/logger"
	sio "../utils/sio"

	"github.com/jessevdk/go-flags"
)

type Message = sio.Message

type ClientOptions struct {
	ServerHost string `short:"h" long:"host" default:"localhost" description:"Server to connect"`
	ServerPort string `short:"P" long:"port" default:"8080" description:"Server port"`
	UserName   string `short:"u" long:"user" description:"Chat user"`
	Password   string `short:"p" long:"password" description:"Chat user password"`
}

func main() {
	var opts ClientOptions
	if _, err := flags.NewParser(&opts, flags.Default).Parse(); err != nil {
		panic(err)
	}

	logger.Log(logger.INFO, "Connect to "+opts.ServerHost+":"+opts.ServerPort)
	conn, err := net.Dial("tcp", opts.ServerHost+":"+opts.ServerPort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	logger.Log(logger.INFO, "Connected")

	// Цикл обработчика команд с консоли ?
	for {
		if err := sio.Send(conn, Message{Raw: "texttest"}); err != nil {
			panic(err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
