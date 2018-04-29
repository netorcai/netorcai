#!/usr/bin/env python3
import client

C = client.Client()
C.connect()
C._send_json({
    "message_type": "LOGIN",
    "nickname": "okay"})

msg = C.recv_json()
print(msg)
assert msg["message_type"] == "KICK"
assert "Invalid first message" in msg["kick_reason"]
assert "Field 'role' is missing" in msg["kick_reason"]
