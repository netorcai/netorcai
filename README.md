[![Build Status](https://img.shields.io/travis/netorcai/netorcai/master.png)](https://travis-ci.org/netorcai/netorcai)
[![Coverage Status](https://img.shields.io/coveralls/github/netorcai/netorcai/master.png)](https://coveralls.io/github/netorcai/netorcai?branch=master)

netorcai
========

![netorcai architecture](./doc/archi.png "netorcai architecture")

netorcai is a network orchestrator for artificial intelligence games.
It splits a classical game server process into two processes, allowing to
develop various games in any language without having to manage all
network-related issues about the clients.
This is done thanks to a [metaprotocol](./doc/metaprotocol.md).

Why?
====
In the context of [Lionel Martin's challenge][challenge lionel martin],
I have been involved in the implementation of multiagent network
games meant to be played by bots.

After implementing several games ([2014][spaceships], [2016][aquar.iom]) I
came to the following conclusions.
- Implenting the network server is tough.
- Handling the clients correctly (errors, fairness, not flooding slow clients...) mostly means that most of the development time is in the network game server, not in the game itself.
- The games in this context are quite specific (fair, turn-based,
visualizable, no big performance constraint), which means the development
effort can be shared regarding the network server.

Installation
============
As netorcai is implemented in Go, it can be built with the `go` command.

Install a recent Go version then run
`go get github.com/netorcai/netorcai/cmd/netorcai` to retrieve the executable
in `${GOPATH}/bin` (if the `GOPATH` environment variable is unset,
it should default to `${HOME}/go` or `%USERPROFILE%\go`).

```bash
    go get github.com/netorcai/netorcai/cmd/netorcai
    ${GOPATH:-${HOME}/go}/bin/netorcai --help
```

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
