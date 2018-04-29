#!/usr/bin/env python3
"""Handle a netorcai metaprotocol client."""
import json
import socket
import struct

class Client:
    """Handles a netorcai metaprotocol client."""
    def __init__(self, hostname=None, port=None):
        hostname = "localhost" if hostname is None else hostname
        port = 4242 if port is None else hostname

        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect((hostname, port))

    def __del__(self):
        self.sock.close()

    def _send_string(self, string):
        buf = string.encode("utf-8")
        self.sock.send(struct.pack("<H", len(buf)+1))
        self.sock.send(buf)
        self.sock.send("\n".encode("utf-8"))

    def _send_json(self, data):
        self._send_string(json.dumps(data))

    def _recv_string(self):
        buf = self.sock.recv(2)
        content_size = struct.unpack("<H", buf)[0]
        buf = self.sock.recv(content_size)
        return buf.decode('utf-8')

    def recv_json(self):
        """Receive and returns a message as a JSON dictionary."""
        msg = self._recv_string()
        return json.loads(msg)

    def send_login(self, nickname, role):
        """Send a LOGIN message."""
        self._send_json({
            "message_type":"LOGIN",
            "nickname": nickname,
            "role": role})
