#!/usr/bin/env python3
from collections import namedtuple
from time import sleep
import client
import sys

LoginInfo = namedtuple('LoginInfo', ['client', 'login_reply'])

nb_visus = 10
nb_visus_max = 4
info = []

# Do sequential LOGIN attempts
for i in range(nb_visus):
    C = client.Client()
    C.connect()
    C.send_login("okay", "visualization")
    reply = C.recv_json()
    info.append(LoginInfo(C, reply))

# Check LOGIN results
error = False
for i, inf in enumerate(info):
    if i < nb_visus_max:
        if inf.login_reply["message_type"] != "LOGIN_ACK":
            print("visu {}: Expected LOGIN_ACK, got {}".format(i, inf.login_reply))
            error = True
    else:
        if inf.login_reply["message_type"] != "KICK":
            print("visu {}: Expected KICK, got {}".format(i, inf.login_reply))
            error = True

if error:
    sys.exit(1)

# Close sockets so there is room for one visu
for i in range(nb_visus_max - 1, nb_visus):
    info[i].client.close()

# Wait some time
sleep(0.2)

# Connect a new visu
C = client.Client()
C.connect()
C.send_login("okay", "visualization")
reply = C.recv_json()
if reply["message_type"] != "LOGIN_ACK":
    print("Last visu: Expected LOGIN_ACK, got {}".format(reply))
    error = True

# Leave
sys.exit(int(error))
