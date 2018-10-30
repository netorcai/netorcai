[![Build Status](https://img.shields.io/travis/netorcai/netorcai/master.png)](https://travis-ci.org/netorcai/netorcai)
[![Coverage Status](https://img.shields.io/coveralls/github/netorcai/netorcai/master.png)](https://coveralls.io/github/netorcai/netorcai?branch=master)

A network orchestrator for artificial intelligence games.

Why?
====
In the context of [Lionel Martin's challenge][challenge lionel martin],
I have been involved in the implementation of multiagent network
games meant to be played by bots.
This was very interesting, as such projects gather multiple tasks
such as game design, network protocol design, game server and visualization
implementation.

I however came to the following conclusions after implementing multiple games
([2014][spaceships], [2016][aquar.iom]).  
First, implementing the network server is tough.
Handling the clients correctly
(errors, fairness, not punishing slow clients too much...)
mostly means that you will spend most of your time in the
network part of the game server.  
Second, in this context we want fair turn-based games that can be displayed
easily on a projector.
The game and network protocol design is therefore limited.
Reimplementing the network part for each new game could therefore be avoided.

Main idea
=========
Netorcai proposes to separate a classical game server into two components:
- Game logic:
  - How to apply clients' actions
  - How to compute a turn
  - What to send to the clients
- Network orchestration:
  - Manage calls to the game logic
  - Manage communications with the clients
  - Manage the state of each client

These two components (and the player and visualization clients) are
instantiated as processes and communicate thanks to a
[metaprotocol][metaprotocol].

Frequent questions / issues
===========================

Running netorcai in my scripts gives an ioctl error
---------------------------------------------------
Try using the `--simple-prompt` option.

Running netorcai in background does not work in my scripts
----------------------------------------------------------
Try launching netorcai via `nohup`.

[//]: =========================================================================
[challenge lionel martin]: https://www.univ-orleans.fr/iut-orleans/informatique/intra/concours/
[aquar.iom]: https://github.com/mpoquet/aquar.iom
[spaceships]: https://github.com/mpoquet/concoursiuto2015
[metaprotocol]: ./doc/metaprotocol.md
