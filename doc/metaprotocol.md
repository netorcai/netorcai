Network protocol description
============================
This protocol is based on TCP and is *mostly* textual, as all messages are
composed by two parts:
1. CONTENT_SIZE, a 16-bit little-endian unsigned integer corresponding to
   the size of the message content (therefore excluding the 2 octets used
   to store CONTENT_SIZE).
2. CONTENT, an UTF-8 string of CONTENT_SIZE octets,
   terminated by an UTF-8 *Line Feed* character (U+000A).

The content of each message must be a valid
[JSON](https://www.json.org/) object.  
Messages are [typed](#message-types) and clients must follow their
[specified behaviour](#expected-client-behavior).

Network entities (endpoints)
----------------------------
This protocol allows multiple entities to communicate:
- the unique **game logic** entity, in charge of managing the game itself.
- **clients** entities, that are in one of the following types:
  - *player*, in charge of taking actions to play the game
  - *visualization*, in charge of displaying the game progress
- the unique **netorcai** entity, the central orchestrator (broker) between
  the game logic and the clients.

![entities figure](./fig/entities.svg "entities figure")

Message types
-------------
Each message has a type.
This type is set as a string in the `message_type` field of the main message
JSON object.
The other fields of the main JSON object depend on the message type.

List of messages between **clients** and **netorcai**:
- [LOGIN](#login)
- [LOGIN_ACK](#login_ack)
- [KICK](#kick)
- [GAME_STARTS](#game_starts)
- [GAME_ENDS](#game_ends)
- [TURN](#turn)
- [TURN_ACK](#turn_ack)

List of messages between **netorcai** and **game logic**:
- ([LOGIN](#login))
- ([LOGIN_ACK](#login_ack))
- ([KICK](#kick))
- [DO_INIT](#do_init)
- [DO_INIT_ACK](#do_init_ack)
- [DO_TURN](#do_turn)
- [DO_TURN_ACK](#do_turn_ack)

### LOGIN
This message type is sent from (**clients** or **game logic**) to **netorcai**.

This is the first message sent by clients and game logic.
It allows them to indicate they want to participate in the game.
**netorcai** answers this message with a [LOGIN_ACK](#login_ack) message
if the logging in is accepted, or by a [KICK](#kick) message otherwise.

Fields:
- `nickname` (string): The name the clients wants to have.
  Must respect the `\A\S{1,10}\z` regular expression (go syntax).
- `role` (string). Must be `player`, `visualization` or `game logic`.

Example:
```json
{
  "message_type": "LOGIN",
  "nickname": "strutser",
  "role": "player"
}
```

### LOGIN_ACK
This message type is sent from **netorcai** to **clients** or **game logic**.

It tells a client or the game logic that its [LOGIN](#login) is accepted.

Fields: None.

Example:
```json
{
  "message_type": "LOGIN_ACK"
}
```

### KICK
This message type is sent from **netorcai** to **clients** or **game logic**.

It tells a client (or game logic) that it is about to be kicked out of a game.
After sending this message, **netorcai** will no longer read incoming messages
from the kicked client (or game logic).
It also means that **netorcai** is about to close the socket.

It can be sent for multiple reasons:
- As a negative acknowledge to a [LOGIN](#login) message
- If a message is invalid: Its content is not valid JSON, a message
  field is missing or has an invalid value
  (as specified in each [message type](#message-types)).
- If a client does not follow its
  [expected behavior](#expected_client_behavior).
- If **netorcai** is about to terminate.

Fields:
- `kick_reason` (string): The reason why the client (or game logic) has been
  kicked

Example:
```json
{
  "message_type": "KICK",
  "kick_reason": "Invalid message: Content is not valid JSON"
}
```

### GAME_STARTS
This message type is sent from **netorcai** to **clients**.

It tells the client that the game is about to start.

Fields:
- `player_id`: (integral non-negative number or -1):
  The unique identifier of the client if its role is `player`,
  -1 otherwise.
- `players_info`: (array of objects):
  If this message is sent to a `player`, this array is empty.  
  If this message is sent to a `visualization`, this array contains information
  about each player:
  - `player_id` (integral non-negative number): The unique player identifier.
  - `nickname` (string): The player nickname.
  - `remote_address` (string): The player network remote address.
  - `is_connected` (bool): Whether the player is currently connected to
    **netorcai**.
- `nb_players` (integral positive number): The number of players of the game.
- `nb_turns_max` (integral positive number):
  The maximum number of turns of the game.
- `milliseconds_before_first_turn` (non-negative number):
  The number of milliseconds before the first game [TURN](#turn).
- `milliseconds_between_turns` (non-negative number):
  The minimum number of milliseconds between two consecutive
  game [TURNs](#turn).
- `initial_game_state` (object): Game-dependent content.

Example:
```json
{
  "message_type": "GAME_STARTS",
  "player_id": -1,
  "players_info": [
    {
      "player_id": 0,
      "nickname": "jugador",
      "remote_address": "127.0.0.1:59840",
      "is_connected": true
    }
  ],
  "nb_players": 4,
  "nb_turns_max": 100,
  "milliseconds_before_first_turn": 1000,
  "milliseconds_between_turns": 1000,
  "initial_game_state": {}
}
```

### GAME_ENDS
This message type is sent from **netorcai** to **clients**.

It tells the client that the game is finished.
The client can safely close the socket after receiving this message.

Fields:
- `winner_player_id` (integral non-negative number or -1):
  The unique identifier of the player that won the game.
  Can be -1 if there is no winner.
- `game_state` (object): Game-dependent content.

Example:
```json
{
  "message_type": "GAME_ENDS",
  "winner_player_id": 0,
  "game_state": {}
}
```

### TURN
This message type is sent from **netorcai** to **clients**.

It tells the client a new turn has started.

Fields:
- `turn_number` (non-negative integral number):
  The number of the current turn.
- `game_state` (object): Game-dependent content that directly corresponds to
  the `game_state`field of a `DO_TURN_ACK` message.
- `players_info`: (array of objects):
  If this message is sent to a `player`, the array is empty.  
  If this message is sent to a `visualization`, the array contains information
  about each player:
  - `player_id` (integral non-negative number): The unique player identifier.
  - `nickname` (string): The player nickname.
  - `remote_address` (string): The player network remote address.
  - `is_connected` (bool): Whether the player is currently connected to
    **netorcai**.

Example:
```json
{
  "message_type": "TURN",
  "turn_number": 0,
  "game_state": {},
  "players_info": [
    {
      "player_id": 0,
      "nickname": "jugador",
      "remote_address": "127.0.0.1:59840",
      "is_connected": true
    }
  ]
}
```

### TURN_ACK
This message type is sent from **clients** to **netorcai**.

It tells netorcai that the client has managed a turn.
For players, it contains the actions the player wants to do.

Fields:
- `turn_number` (non-negative integral number):
  The number of the turn that the client has managed.
  Value must match the `turn_number` of the latest `TURN` received by
  the client.
- `actions` (array): Game-dependent content. Must be empty for visualizations.

Example:
```json
{
  "message_type": "TURN_ACK",
  "turn_number": 0,
  "actions": []
}
```

### DO_INIT
This message type is sent from **netorcai** to **game logic**.

This message initiates the sequence to start the game.
**netorcai** gives information to the game logic, such that the game logic can generate the game initial state.

Fields:
- `nb_players` (integral positive number): The number of players in the game.
- `nb_turns_max` (integral positive number):
  The maximum number of turns of the game.

Example:
```json
{
  "message_type": "DO_INIT",
  "nb_players": 4,
  "nb_turns_max": 100
}
```

### DO_INIT_ACK
This message is sent from **game logic** to **netorcai**.

It means that the game logic has finished its initialization.
It sends initial information about the game, which is forwarded to the clients.

Fields:
- `initial_game_state` (object): The initial game state, as it should be
  transmitted to clients.
  Only the `all_clients` key of this object is currently implemented, which
  means the associated game-dependent object will be transmitted to all the
  clients (players and visualizations).

Example:
```json
{
  "initial_game_state": {
    "all_clients": {}
  }
}
```

### DO_TURN
This message type is sent from **netorcai** to **game logic**.

It tells the game logic to do a new turn.

Fields:
- `player_actions` (array): The actions decided by the players.
  There is at most one array element per player.
  This array contains objects that must contain the following fields:
    - `player_id` (non-negative integral number):
      The unique identifier of the player who decided the actions.
    - `turn_number` (non-negative integral number):
      The turn whose the actions comes from
      (received from [TURN_ACK](#turn_ack)).
    - `actions` (array): The actions of the player.
      Game-dependent content (received from [TURN_ACK](#turn_ack)).

Example:
```json
{
  "message_type": "DO_TURN",
  "player_actions": [
    {
      "player_id": 0,
      "turn_number": 0,
      "actions": []
    }
  ]
}
```

### DO_TURN_ACK
This message type is sent from **game logic** to **netorcai**.

Game logic has computed a new turn and transmits its results.

Fields:
- `winner_player_id` (non-negative integral number or -1):
  The unique identifier of the player currently winning the game.
  Can be -1 if there is no current winner.
- `game_state` (object): The current game state, as it should be
  transmitted to clients.
  Only the `all_clients` key of this object is currently implemented, which
  means the associated game-dependent object will be transmitted to all the
  clients (players and visualizations).

Example:
```json
{
  "message_type": "DO_TURN_ACK",
  "winner_player_id": 0,
  "game_state": {
    "all_clients": {}
  }
}
```

## Expected client behavior
**netorcai** manages the clients by associating them with a state.
In a given state, a client can only receive and send certain types of messages.
A client that sends an unexpected type of message is [kicked](#kick)
by **netorcai**.

The following figure summarizes the expected behavior of a client:
- Each node is a client state.
- Edges are transitions between states:
	- ?MSG_TYPE means that the client receives a message of type MSG_TYPE.
	- !MSG_TYPE means that the client sends a message of type MSG_TYPE.

![client expected behavior figure](./fig/expected_behavior_client.svg "client expected behavior figure")

## Expected game logic behavior
Similarly to clients, **netorcai** manages the game logic by associating it
with a state.

Its expected behavior is described in the following figure:

![game logic expected behavior figure](./fig/expected_behavior_gamelogic.svg "game logic expected behavior figure")
