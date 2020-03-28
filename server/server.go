package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"time"

	logger "../utils/logger"
	sio "../utils/sio"

	"github.com/jessevdk/go-flags"
	uuid "github.com/nu7hatch/gouuid"
)

type ServerOptions struct {
	ConfigFile string `short:"c" long:"config" default:"config.json" description:"Path to the configuration file" `
}

type NetworkConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type DBConfig struct {
	DBHost string `json:"host"`
	DBPort string `json:"port"`
	DBName string `json:"name"`
	DBUser string `json:"user"`
	DBPass string `json:"password"`
}

type ServerConfig struct {
	Network              NetworkConfig `json:"server"`
	DB                   DBConfig      `json:"db"`
	ConnectionTimeoutSec int64         `json:"client_connection_timeout_sec"`
}

type ClientSession struct {
	SessionId   string
	Conn        net.Conn
	ConnectedIn time.Time
	UserName    string
	isAuth      bool
}

type serverContext struct {
	clientSessions map[string]ClientSession
	serverRouter   router
	// Указатель на БД
}

func createContext() serverContext {
	var sc serverContext
	sc.clientSessions = make(map[string]ClientSession)
	sc.serverRouter = createRouter()
	return sc
}

func main() {
	var opts ServerOptions
	if _, err := flags.NewParser(&opts, flags.Default).Parse(); err != nil {
		panic(err)
	}

	var cfg = parseConfigFile(opts.ConfigFile)
	println("Starting server on ", cfg.Network.Host, ":", cfg.Network.Port)

	listener, err := net.Listen("tcp", cfg.Network.Host+":"+cfg.Network.Port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	var serverContext serverContext = createContext()
	serverContext.serverRouter.addRoute("login", loginHandler)

	go clearSessions(serverContext.clientSessions, cfg)
	for {
		connection, err := listener.Accept()
		if err != nil {
			println(err.Error())
			continue
		}

		cuuid, _ := uuid.NewV4()
		client := ClientSession{SessionId: cuuid.String(), Conn: connection, ConnectedIn: time.Now(), isAuth: false}
		serverContext.clientSessions[client.SessionId] = client

		go clientHandler(serverContext, client)
	}
}

func parseConfigFile(ConfigFileName string) (cfg ServerConfig) {
	configFile, err := os.Open(ConfigFileName)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	configBytes, _ := ioutil.ReadAll(configFile)
	if err = json.Unmarshal(configBytes, &cfg); err != nil {
		panic(err)
	}

	return cfg
}

func clearSessions(clientSessions map[string]ClientSession, cfg ServerConfig) {
	for {
		for _, session := range clientSessions {
			if session.isAuth == false && time.Now().Unix()-session.ConnectedIn.Unix() > cfg.ConnectionTimeoutSec {
				println("Kill session: " + session.SessionId)
				delete(clientSessions, session.SessionId)
				session.Conn.Close()
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func clientHandler(context serverContext, client ClientSession) {
	println("New connection from " + client.Conn.RemoteAddr().String())
	println("Seesion id: " + client.SessionId)

	// main loop
	for {
		msg, err := sio.Recv(client.Conn)
		if err != nil {
			logger.Log(logger.ERR, "Error on reading data")
			break
		}

		logger.Log(logger.INFO, "Recv: "+msg.PostData)

		if handler, err := context.serverRouter.getHandler(msg.Route); err == nil {
			respCode := handler(context, msg.PostData)
			sio.Send(client.Conn, sio.Message{Route: msg.Route, Code: respCode})
		} else {
			logger.Log(logger.ERR, "Unknown route")
		}
	}
}

// TODO: подключить базу и сделать аутентификацию
func loginHandler(context serverContext, postData string) (statusCode int) {
	return 200
}
