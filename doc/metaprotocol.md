Network protocol description
============================
This protocol is based on TCP and is *mostly* textual, as all messages are
composed by two parts:
1. CONTENT_SIZE, a 16-bit little-endian unsigned integer corresponding to
   the size of the message content (therefore excluding the 2 octets used
   to store CONTENT_SIZE).
2. CONTENT, an UTF-8 string of MESSAGE_SIZE octets,
   terminated by an UTF-8 *Line Feed* character (U+000A).

The content of each message must be a valid
[JSON](https://www.json.org/) object.

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
- [GAME_STARTS](#game_starts)
- [GAME_ENDS](#game_ends)
- [TURN](#turn)
- [TURN_ACK](#turn_ack)
- [KICK](#kick)

List of messages between **netorcai** and **game logic**:
- [DO_TURN](#do_turn)
- [DO_TURN_ACK](#do_turn_ack)

### LOGIN
This message type is sent from **clients** to **netorcai**.

This is the first message sent by a client.
It allows the client to indicate it wants to participate in the game.
**netorcai** answers this message with a [LOGIN_ACK](#login_ack) message
if the logging in is accepted, or by a [KICK](#kick) message otherwise.

Fields:
- `nickname` (string): The name the clients wants to have.
  Must respect the `\A\S{1,10}\z` regular expression (go syntax).
- `role` (string). Must be `player` or `visualization`.

Example:
```json
{
  "message_type": "LOGIN",
  "nickname": "strutser",
  "role": "player"
}
```

### LOGIN_ACK
This message type is sent from **netorcai** to **clients**.

It tells a client that his [LOGIN](#login) is accepted.

Fields: None.

Example:
```json
{
  "message_type": "LOGIN_ACK"
}
```

### GAME_STARTS
This message type is sent from **netorcai** to **clients**.

It tells the clients that the game is about to start.

Fields:
- `player_id`: (integral non-negative number or `null`):
  The unique identifier of the client if its role is `player`,
  `null` otherwise.
- `nb_players` (integral positive number): The number of players of the game.
- `nb_turns_max` (integral number): The maximum number of turns of the game.
- `milliseconds_before_first_turn` (number): The number of milliseconds before
  the first game [TURN](#turn).
- `data` (object): Game-dependent content.

Example:
```json
{
  "message_type": "GAME_STARTS",
  "player_id": 0,
  "nb_players": 4,
  "nb_turns_max": 100,
  "milliseconds_before_first_turn": 1000,
  "data": {}
}
```

### GAME_ENDS
This message type is sent from **netorcai** to **clients**.

It tells the clients that the game is finished.
Clients can safely close the socket after receiving this message.

Fields:
- `winner_player_id` (integral non-negative number or null):
  The unique identifier of the player that won the game.
  Can be null if there is no winner.
- `data` (object): Game-dependent content.

Example:
```json
{
  "message_type": "GAME_ENDS",
  "winner_player_id": 0,
  "data": {}
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

Example:
```json
{
  "message_type": "TURN",
  "turn_number": 0,
  "game_state": {}
}
```

### TURN_ACK
This message type is sent from **clients** to **netorcai**.

It tells netorcai that the client managed a turn.
For players, it contains the actions it wants to do.

Fields:
- `turn_number` (non-negative integral number):
  The number of the turn the client has managed.
- `actions` (array): Game-dependent content. Must be empty for visualizations.

Example:
```json
{
  "message_type": "TURN_ACK",
  "turn_number": 0,
  "actions": []
}
```

### KICK
This message type is sent from **netorcai** to **clients**.

Kicks a player from a game. After sending this message,
**netorcai** will no longer read messages from the kicked client and is about
to close the socket.

It can be sent for multiple reasons:
- As a negative acknowledge to a [LOGIN](#login) message
- If a message is invalid: Its content is not valid JSON or a message
  field is missing.
- If a client does not follow its
  [expected behavior](#expected_client_behavior).

Fields:
- `kick_reason` (string): The reason why the client has been kicked

Example:
```json
{
  "message_type": "KICK",
  "kick_reason": "Invalid message: Content is not valid JSON"
}
```

### DO_TURN
This message type is sent from **netorcai** to **game logic**.

It tells the game logic to do a new turn.

Fields:
- `player_actions` (array): Game-dependent content.
  Each element of this array exactly correspond to the `actions` field of the
  [TURN_ACK](#turn_ack) message.

Example:
```json
{
  "message_type": "DO_TURN",
  "player_actions": []
}
```

### DO_TURN_ACK
This message type is sent from **game logic** to **netorcai**.

Game logic has computed a new turn and transmits its results.

Fields:
- `winner_player_id` (non-negative integral number or null):
  The unique identifier of the player currently winning the game.
  Can be null if there is no current winner.
- `game_state` (object): Game-dependent content.

Example:
```json
{
  "message_type": "DO_TURN_ACK",
  "game_state": {}
}
```

## Expected client behavior
**netorcai** manages the clients by associating them with a state.
In a given state, a client can only receive and send certain types of messages.

The following figure summarizes the expected behavior of a client:
- ?MSG_TYPE means that a client receives a message of type MSG_TYPE.
- !MSG_TYPE means that a client sends a message of type MSG_TYPE.
