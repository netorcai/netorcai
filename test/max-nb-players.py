#!/usr/bin/env python3
from collections import namedtuple
from time import sleep
import client
import sys

LoginInfo = namedtuple('LoginInfo', ['client', 'login_reply'])

nb_players = 10
nb_players_max = 4
info = []

# Do sequential LOGIN attempts
for i in range(nb_players):
    C = client.Client()
    C.connect()
    C.send_login("okay", "player")
    reply = C.recv_json()
    info.append(LoginInfo(C, reply))

# Check LOGIN results
error = False
for i, inf in enumerate(info):
    if i < nb_players_max:
        if inf.login_reply["message_type"] != "LOGIN_ACK":
            print("Player {}: Expected LOGIN_ACK, got {}".format(i, inf.login_reply))
            error = True
    else:
        if inf.login_reply["message_type"] != "KICK":
            print("Player {}: Expected KICK, got {}".format(i, inf.login_reply))
            error = True

if error:
    sys.exit(1)

# Close sockets so there is room for one player
for i in range(nb_players_max - 1, nb_players):
    info[i].client.close()

# Wait some time
sleep(0.2)

# Connect a new player
C = client.Client()
C.connect()
C.send_login("okay", "player")
reply = C.recv_json()
if reply["message_type"] != "LOGIN_ACK":
    print("Last player: Expected LOGIN_ACK, got {}".format(reply))
    error = True

# Leave
sys.exit(int(error))
