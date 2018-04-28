Network protocol description
============================

This protocol is based on TCP and is *mostly* textual, as all messages are
composed by two parts:
1. SIZE, a 32-bit little-endian unsigned integer corresponding to the message
   size in octets (excluding the 4 octets needed to store SIZE).
2. CONTENT, an UTF-8 string of SIZE octets without terminating character.

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
