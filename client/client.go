package main

import (
	"bufio"
	"net"
	"os"
	"strings"
	"time"

	logger "../utils/logger"
	sio "../utils/sio"

	"github.com/jessevdk/go-flags"
)

type Message = sio.Message

type clientOptions struct {
	ServerHost string `short:"h" long:"host" default:"localhost" description:"Server to connect"`
	ServerPort string `short:"P" long:"port" default:"8080" description:"Server port"`
}

type cliCommand struct {
	action string
	params string
}

var (
	currentRequestId int
)

func init() {
	currentRequestId = 0
}

func main() {
	var opts clientOptions
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

	serverChan := make(chan Message)
	cliChan := make(chan cliCommand)

	go waitCliCommand(cliChan)
	go waitServerCommand(conn, serverChan)

	for {
		select {
		case cmd := <-cliChan:
			handleCliCommand(conn, cmd)
		case recvMsg := <-serverChan:
			handleServerMessage(conn, recvMsg)
		case <-time.After(100 * time.Millisecond):
			continue
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func waitCliCommand(cliChan chan cliCommand) {
	for {
		if cmd, err := getCommand(); err == nil {
			cliChan <- cmd
		} else {
			panic(err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func waitServerCommand(conn net.Conn, serverChan chan Message) {
	for {
		if recvMsg, err := sio.Recv(conn); err == nil {
			serverChan <- recvMsg
		} else {
			panic(err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func handleCliCommand(conn net.Conn, cmd cliCommand) {
	println("!!!!:", cmd.action)
	switch cmd.action {
	case "help", "h", "Help", "H":
		printHelp()
	case "login", "send": // можно не утруждать клиента валидацией экшена
		sendMsg := Message{Route: cmd.action, PostData: cmd.params, RequestId: currentRequestId}
		if err := sio.Send(conn, sendMsg); err != nil {
			panic(err)
		}
	case "exit":
		sendMsg := Message{Route: "logout"}
		if err := sio.Send(conn, sendMsg); err != nil {
			panic(err)
		}
		conn.Close()
		os.Exit(0)
	default:
		printHelp()
	}
}

func handleServerMessage(conn net.Conn, msg Message) {
	switch msg.Route {
	case "send":
		if msg.Code == 200 {
			println(msg.PostData)
		}
	case "login":
		if msg.Code == 200 {
			println("Auth succcess")
		} else {
			println("Auth failed")
		}
	}
}

func getCommand() (cliCommand, error) {
	print("> ")
	in := bufio.NewReader(os.Stdin)
	cmdRaw, err := in.ReadString('\n')
	if err != nil {
		return cliCommand{}, err
	}

	var action, params string
	if spaceIndex := strings.Index(cmdRaw, " "); spaceIndex == -1 {
		action = cmdRaw[:len(cmdRaw)-1]
	} else {
		action = cmdRaw[:spaceIndex]
		params = cmdRaw[spaceIndex+1 : len(cmdRaw)-1]
	}

	return cliCommand{action: action, params: params}, nil
}

func printHelp() {
	println("Help:")
	println("help - show this help")
	println("exit - close connection and exit")
	println("login user pass - authorize on server by login password")
	println("send text - send message to all users on server")
}
