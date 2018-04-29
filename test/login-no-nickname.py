#!/usr/bin/env python3
import client

C = client.Client()
C.connect()
C._send_json({
    "message_type": "LOGIN",
    "role": "player"})

msg = C.recv_json()
print(msg)
assert msg["message_type"] == "KICK"
assert "Invalid first message" in msg["kick_reason"]
assert "Field 'nickname' is missing" in msg["kick_reason"]
