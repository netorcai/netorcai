#!/usr/bin/env python3
from collections import namedtuple
from time import sleep
import client
import sys

LoginInfo = namedtuple('LoginInfo', ['client', 'login_reply'])

nb_game_logics = 10
nb_game_logics_max = 1
info = []

# Do sequential LOGIN attempts
for i in range(nb_game_logics):
    C = client.Client()
    C.connect()
    C.send_login("okay", "game logic")
    reply = C.recv_json()
    info.append(LoginInfo(C, reply))

# Check LOGIN results
error = False
for i, inf in enumerate(info):
    if i < nb_game_logics_max:
        if inf.login_reply["message_type"] != "LOGIN_ACK":
            print("game_logic {}: Expected LOGIN_ACK, got {}".format(i, inf.login_reply))
            error = True
    else:
        if inf.login_reply["message_type"] != "KICK":
            print("game_logic {}: Expected KICK, got {}".format(i, inf.login_reply))
            error = True

if error:
    sys.exit(1)

# Close sockets so there is room for one game_logic
for i in range(nb_game_logics):
    info[i].client.close()

# Wait some time
sleep(0.2)

# Connect a new game_logic
C = client.Client()
C.connect()
C.send_login("okay", "game logic")
reply = C.recv_json()
if reply["message_type"] != "LOGIN_ACK":
    print("Last game_logic: Expected LOGIN_ACK, got {}".format(reply))
    error = True

# Leave
sys.exit(int(error))
