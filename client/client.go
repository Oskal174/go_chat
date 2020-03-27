package main

import (
	"bufio"
	"net"
	"os"
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
	cliChan := make(chan Message)

	go handleCliCommand(cliChan)
	go handleServerMsg(conn, serverChan)

	for {
		select {
		case msg := <-cliChan:
			sendMsg(conn, msg)
		case msg := <-serverChan:
			logger.Log(logger.INFO, msg.PostData)
		case <-time.After(500 * time.Millisecond):
			continue
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func handleCliCommand(cliChan chan Message) {
handleCliCommandLoop:
	for {
		cmd := getCommand()
		switch cmd.action {
		case "help", "h", "Help", "H":
			println("PRINT HELP TODO")
		case "exit":
			cliChan <- Message{Route: "/exit", PostData: "", RequestId: currentRequestId}
			break handleCliCommandLoop
		default:
			println("PRINT HELP TODO")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func handleServerMsg(conn net.Conn, serverChan chan Message) {
	for {

		time.Sleep(100 * time.Millisecond)
	}
}

// TODO: сделать проверку ситаксиса входной команды
func getCommand() cliCommand {
	print("> ")
	in := bufio.NewReader(os.Stdin)
	cmdRaw, _ := in.ReadString('\n')

	var cmd = cliCommand{}
	for i, ch := range cmdRaw {
		if ch == ' ' {
			cmd = cliCommand{action: cmdRaw[:i], params: cmdRaw[i+1:]}
			break
		}
	}

	return cmd
}

func sendMsg(conn net.Conn, msg Message) {
	if err := sio.Send(conn, msg); err != nil {
		panic(err)
	}
}
