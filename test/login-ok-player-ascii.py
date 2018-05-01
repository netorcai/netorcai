#!/usr/bin/env python3
import client

C = client.Client()
C.connect()
C.send_login("okay", "player")

msg = C.recv_json()
print(msg)
assert msg["message_type"] == "LOGIN_ACK"
