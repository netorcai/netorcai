#!/usr/bin/env python3
import client

C = client.Client()
C.connect()
C._send_json({
    "message_type": "LOGIN",
    "nickname": "okay",
    "role": "player"})

msg = C.recv_json()
print(msg)
assert msg["message_type"] == "LOGIN_ACK"
