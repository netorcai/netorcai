package client

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c *Client) Connect(hostname string, port int) error {
	var err error
	connectAddress := hostname + ":" + strconv.Itoa(port)

	c.conn, err = net.Dial("tcp", connectAddress)
	if err != nil {
		return err
	}

	c.reader = bufio.NewReader(c.conn)
	c.writer = bufio.NewWriter(c.conn)
	return nil
}

func (c *Client) Disconnect() error {
	c.reader = nil
	c.writer = nil
	return c.conn.Close()
}

func (c *Client) SendBytes(content []byte) error {
	contentSize := len(content)
	if contentSize >= 65535 {
		return fmt.Errorf("content too big: size does not fit in 16 bits")
	}

	// Write content size on socket
	var contentSizeUint16 uint16 = uint16(contentSize) + 1 // +1 for \n
	contentSizeBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(contentSizeBuf, contentSizeUint16)
	_, err := c.writer.Write(contentSizeBuf)
	if err != nil {
		return fmt.Errorf("Remote endpoint closed? Write error: %v", err)
	}

	// Write content on socket
	_, err = c.writer.Write(content)
	if err != nil {
		return fmt.Errorf("Remote endpoint closed? Write error: %v", err)
	}

	// Write terminating "\n" character on socket
	err = c.writer.WriteByte(0x0A)
	if err != nil {
		return fmt.Errorf("Remote endpoint closed? Write error: %v", err)
	}

	// Flush socket
	c.writer.Flush()
	return nil
}

func (c *Client) SendString(str string) error {
	return c.SendBytes([]byte(str))
}

func (c *Client) SendJSON(msg map[string]interface{}) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Cannot marshall JSON message: %v", err)
	} else {
		return c.SendBytes(content)
	}
}

func (c *Client) SendLogin(role, nickname, metaprotocolVersion string) error {
	msg := map[string]interface{}{
		"message_type":         "LOGIN",
		"role":                 role,
		"nickname":             nickname,
		"metaprotocol_version": metaprotocolVersion,
	}

	return c.SendJSON(msg)
}

func (c *Client) ReadMessage() (map[string]interface{}, error) {
	var msg map[string]interface{}
	contentSizeBuf := make([]byte, 2)
	_, err := io.ReadFull(c.reader, contentSizeBuf)
	if err != nil {
		return msg, fmt.Errorf("Remote endpoint closed? Read error: %v", err)
	}

	// Read message content size
	contentSize := binary.LittleEndian.Uint16(contentSizeBuf)

	// Receive message content
	contentBuf := make([]byte, contentSize)
	_, err = io.ReadFull(c.reader, contentBuf)
	if err != nil {
		return msg, fmt.Errorf("Remote endpoint closed? Read error: %v", err)
	}

	// Read message content
	err = json.Unmarshal(contentBuf, &msg)
	if err != nil {
		return msg, fmt.Errorf("Non-JSON message received")
	} else {
		return msg, nil
	}
}
