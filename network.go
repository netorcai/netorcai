package netorcai

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

type Client struct {
	Conn             net.Conn
	nickname         string
	state            int
	reader           *bufio.Reader
	writer           *bufio.Writer
	incomingMessages chan ClientMessage
	canTerminate     chan string
}

type ClientMessage struct {
	content map[string]interface{}
	err     error
}

func RunServer(port int, globalState *GlobalState, onexit,
	gameLogicExit chan int) {
	defer globalState.WaitGroup.Done()
	// Listen all incoming TCP connections on the specified port
	listenAddress := ":" + strconv.Itoa(port)
	globalState.Mutex.Lock()
	var err error
	globalState.Listener, err = net.Listen("tcp", listenAddress)
	globalState.Mutex.Unlock()
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
	defer globalState.Listener.Close()

	for {
		// Wait for an incoming connection.
		client := &Client{}
		client.Conn, err = globalState.Listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Warn("Could not accept incoming connection. Aborting server.")
			onexit <- 1
			return
		} else {
			// Handle connections in a new goroutine.
			client.reader = bufio.NewReader(client.Conn)
			client.writer = bufio.NewWriter(client.Conn)
			client.state = CLIENT_UNLOGGED
			client.incomingMessages = make(chan ClientMessage)
			client.canTerminate = make(chan string, 1)

			globalState.WaitGroup.Add(1)
			go handleClient(client, globalState, gameLogicExit)
		}
	}
}

func readClientMessage(client *Client, maximumAllowedSize uint32, errorFormatOnTooBigMessage string) bool {
	var msg ClientMessage
	// Receive message content size
	contentSizeBuf := make([]byte, 4)
	_, err := io.ReadFull(client.reader, contentSizeBuf)
	if err != nil {
		msg.err = fmt.Errorf("Remote endpoint closed? Read error: %v", err)
		client.incomingMessages <- msg
		return false
	}

	// Read message content size
	contentSize := binary.LittleEndian.Uint32(contentSizeBuf)
	if contentSize > maximumAllowedSize {
		msg.err = fmt.Errorf(errorFormatOnTooBigMessage, contentSize)
		client.incomingMessages <- msg
		return false
	}

	// Receive message content
	contentBuf := make([]byte, contentSize)
	_, err = io.ReadFull(client.reader, contentBuf)
	if err != nil {
		msg.err = fmt.Errorf("Remote endpoint closed? Read error: %v", err)
		client.incomingMessages <- msg
		return false
	}

	log.WithFields(log.Fields{
		"remote address": client.Conn.RemoteAddr(),
		"nickname":       client.nickname,
		"content size":   contentSize,
		"content":        string(contentBuf),
	}).Debug("New message received")
	// Read message content
	err = json.Unmarshal(contentBuf, &msg.content)
	if err != nil {
		log.WithFields(log.Fields{
			"err":             err,
			"message content": string(contentBuf),
		}).Debug("Non-JSON message received")
		msg.err = fmt.Errorf("Non-JSON message received")
		client.incomingMessages <- msg
		return false
	}

	client.incomingMessages <- msg
	return true
}

func readClientMessages(client *Client) {
	if readClientMessage(client, 1023, "Received message size of first message is too big: %v does not fit in 10 bits") {
		for readClientMessage(client, 16777215, "Received message size is too big: %v does not fit in 24 bits") {
		}
	}
}

func sendMessage(client *Client, content []byte) error {
	// Check content size
	contentSize := len(content)
	if contentSize >= 16777215 {
		return fmt.Errorf("content too big: size does not fit in 24 bits")
	}

	// Write content size on socket
	var contentSizeUint32 uint32 = uint32(contentSize) + 1 // +1 for \n
	contentSizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(contentSizeBuf, contentSizeUint32)
	_, err := client.writer.Write(contentSizeBuf)
	if err != nil {
		return fmt.Errorf("Remote endpoint closed? Write error: %v", err)
	}

	// Write content on socket
	_, err = client.writer.Write(content)
	if err != nil {
		return fmt.Errorf("Remote endpoint closed? Write error: %v", err)
	}

	// Write terminating "\n" character on socket
	err = client.writer.WriteByte(0x0A)
	if err != nil {
		return fmt.Errorf("Remote endpoint closed? Write error: %v", err)
	}

	// Flush socket
	client.writer.Flush()
	return nil
}
