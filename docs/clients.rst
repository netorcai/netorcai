Client libraries
================

The netorcai architecture is a client-server one.
The netorcai program has the role of a network server while
the other entities (games, players and visualizations) have a client role.

Programming netorcai clients can be done completely from scratch,
but libraries have been implemented in several programming languages
to help in their development.
All these libraries should be available in the
`netorcai organization github repository`_.
Currently, the following libraries have been implemented.

- `netorcai-client-cpp`_
- `netorcai-client-d`_
- `netorcai-client-java`_
- `netorcai-client-python`_

Contrary to bindings_, all these libraries are fully implemented
in the target programming language.
The main advantage is that the installation of each library is simplified,
as it can be done directly with the language packaging tools.

Client libraries API
~~~~~~~~~~~~~~~~~~~~

All the client libraries propose the same programming interface.
Inner details may of course vary depending on the programming language,
such as the type used to store collections of items or the
variable/function name depending on the language coding style.
All libraries **must** do the following.

- Provide a high-level :code:`Client` class that manages the network connection.
- Provide structured types for the various messages of the metaprotocol
  (see :ref:`proto_message_types`).
  Each message is implemented as a :code:`struct` in C++ and D,
  and as :code:`class` in Java and Python.
- Provide functions to parse the various metaprotocol messages.

The :code:`Client` class is intended to be the main way to send and receive
netorcai messages. This class provides the following methods.

- Various methods to send metaprotocol messages on the network,
  named :code:`send<MESSAGE_TYPE>` (*e.g.*, :code:`sendLogin`).
- Various methods to receive metaprotocol messages from the network,
  named :code:`recv<MESSAGE_TYPE>` (*e.g.*, :code:`recvLoginAck`).
  **These functions are blocking.**
- :code:`sendString` and :code:`sendJson`,
  that respectively send a user-defined string
  or a user-defined JSON object on the network.
- :code:`recvString` and :code:`recvJson`,
  that respectively receive a string or a JSON object from the network.
  **These two functions are blocking.**

Additionnally, the libraries **may** provide other features.
For example, the C++ API provides a non-blocking API for reading messages.

Usage examples
~~~~~~~~~~~~~~

As an example, here is a basic player bot in Python.

.. code:: python

    try:
        # Instantiate a client in memory.
        client = Client()

        # Connect the internal socket to netorcai.
        client.connect()

        # Log in to netorcai as a player.
        client.send_login("py-player", "player")
        client.read_login_ack()

        # Wait for the game to start.
        game_starts = client.read_game_starts()

        # Do no action on each turn.
        for i in range(game_starts.nb_turns_max):
            turn = client.read_turn()
            client.send_turn_ack(turn.turn_number, [])
    except Exception as e:
        print(e)

All libraries have examples in the :code:`examples` directory of their
respective repository. Please refer to them for more examples.

Getting the libraries
~~~~~~~~~~~~~~~~~~~~~

The most straightforward way to get the up-to-date version of each library is to
clone its git repository and to build and install it from source.
Building and installation instructions are in the README of each repository.

Some of them are also available in the package registry of their language.

- D: `netorcai-client package on DUB`_
- Python: `netorcai package on PyPI`_

Some of them are packaged in Nix_ in the netorcaipkgs_ package repository.

.. code:: bash

    # Install the C++ client library.
    nix-env -f https://github.com/netorcai/netorcaipkgs/archive/master.tar.gz -iA netorcai_client_cpp # latest release
    nix-env -f https://github.com/netorcai/netorcaipkgs/archive/master.tar.gz -iA netorcai_client_cpp_dev # up-to-date

.. _netorcai organization github repository: https://github.com/netorcai/
.. _netorcaipkgs: https://github.com/netorcai/pkgs
.. _netorcai-client-cpp: https://github.com/netorcai/netorcai-client-cpp
.. _netorcai-client-d: https://github.com/netorcai/netorcai-client-d
.. _netorcai-client-java: https://github.com/netorcai/netorcai-client-java
.. _netorcai-client-python: https://github.com/netorcai/netorcai-client-python
.. _Nix: https://nixos.org/nix/
.. _bindings: https://en.wikipedia.org/wiki/Language_binding
.. _netorcai-client package on DUB: https://code.dlang.org/packages/netorcai-client
.. _netorcai package on PyPI: https://pypi.org/project/netorcai/
