package sio

import (
	"encoding/json"
	"net"
)

type Message struct {
	RequestId int    `json:"request_id"`
	Route     string `json:"route"`
	PostData  string `json:"post_data"`
}

const MSG_SIZE = 1024

func Send(conn net.Conn, msg Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if n, err := conn.Write(msgBytes); n < len(msgBytes) || err != nil {
		return err
	}

	return nil
}

func Recv(conn net.Conn) (msg Message, err error) {
	var recvBytes = make([]byte, MSG_SIZE)
	n, err := conn.Read(recvBytes)
	if err != nil {
		return msg, err
	}

	if err := json.Unmarshal(recvBytes[:n], &msg); err != nil {
		return msg, err
	}

	return msg, nil
}
