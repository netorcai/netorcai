package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strconv"
)

func server(port int, globalState *GlobalState, onexit chan int) {
	// Listen all incoming TCP connections on the specified port
	listenAddress := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"err":            err,
			"network":        "tcp",
			"listen address": listenAddress,
		}).Error("Cannot listen incoming connections")
		onexit <- 1
		return
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Listening incoming connections")
	defer listener.Close()

	for {
		// Wait for an incoming connection.
		var client Client
		client.conn, err = listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Warn("Could not accept incoming connection")
		} else {
			// Handle connections in a new goroutine.
			client.reader = bufio.NewReader(client.conn)
			client.writer = bufio.NewWriter(client.conn)
			client.state = CLIENT_UNLOGGED
			client.incomingMessages = make(chan ClientMessage)

			go handleClient(&client, globalState)
		}
	}

	onexit <- 0
}

type Client struct {
	conn             net.Conn
	nickname         string
	state            int
	reader           *bufio.Reader
	writer           *bufio.Writer
	incomingMessages chan ClientMessage
}

type ClientMessage struct {
	content map[string]interface{}
	err     error
}

func readClientMessages(client *Client) {
	var msg ClientMessage

	for {
		// Receive message content size
		contentSizeBuf := make([]byte, 2)
		_, err := io.ReadFull(client.reader, contentSizeBuf)
		if err != nil {
			msg.err = fmt.Errorf("Cannot read content size. "+
				"Remote endpoint closed? Err=%v", err)
			client.incomingMessages <- msg
			return
		}

		// Read message content size
		contentSize := binary.LittleEndian.Uint16(contentSizeBuf)

		// Receive message content
		contentBuf := make([]byte, contentSize)
		_, err = io.ReadFull(client.reader, contentBuf)
		if err != nil {
			msg.err = fmt.Errorf("Cannot read content. "+
				"Remote endpoint closed? Err=%v", err)
			client.incomingMessages <- msg
			return
		}

		// Read message content
		err = json.Unmarshal(contentBuf, &msg.content)
		if err != nil {
			log.WithFields(log.Fields{
				"err":             err,
				"message content": string(contentBuf),
			}).Debug("Non-JSON message received")
			msg.err = fmt.Errorf("Non-JSON message received")
			client.incomingMessages <- msg
			return
		}

		client.incomingMessages <- msg
	}
}

func sendMessage(client *Client, content []byte) error {
	// Check content size
	contentSize := len(content)
	if contentSize >= 65535 {
		return fmt.Errorf("content too big: size does not fit in 16 bits")
	}

	// Write content size on socket
	var contentSizeUint16 uint16 = uint16(contentSize) + 1 // +1 for \n
	contentSizeBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(contentSizeBuf, contentSizeUint16)
	_, err := client.writer.Write(contentSizeBuf)
	if err != nil {
		return fmt.Errorf("Cannot write content size. "+
			"Remote endpoint closed? Err=%v", err)
	}

	// Write content on socket
	_, err = client.writer.Write(content)
	if err != nil {
		return fmt.Errorf("Cannot write content. Remote endpoint closed? "+
			"Err=%v", err)
	}

	// Write terminating "\n" character on socket
	err = client.writer.WriteByte(0x0A)
	if err != nil {
		return fmt.Errorf("Cannot write terminating character. "+
			"Remote endpoint closed? Err=%v", err)
	}

	// Flush socket
	client.writer.Flush()
	return nil
}
