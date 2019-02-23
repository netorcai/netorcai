Client libraries
================

The netorcai architecture is a client-server one.
The netorcai program has the role of a network server while
the other entities (games, players and visualizations) have a client role.

While netorcai clients can be implemented from scratch,
several libraries have been implemented to ease the communication with the netorcai server.
All these libraries are available in the `netorcai organization github repository`_.
Currently, the following libraries have been implemented.

- `netorcai-client-cpp`_
- `netorcai-client-d`_
- `netorcai-client-fortran`_
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
All existing libraries provide the following.

- A high-level :code:`Client` class that manages the network connection.
- Structured types for the various messages of the metaprotocol
  (see :ref:`proto_message_types`).
  Each message is implemented as a :code:`struct` in C++ and D,
  and as :code:`class` in Java and Python.
- Functions to parse the various metaprotocol messages.

The :code:`Client` class is intended to be the main way to send and receive
netorcai messages. This class provides the following methods.

- Various methods to send metaprotocol messages on the network,
  named :code:`send<MESSAGE_TYPE>` (*e.g.*, :code:`sendLogin`).
- Various methods to receive and parse metaprotocol messages from the network,
  named :code:`read<MESSAGE_TYPE>` (*e.g.*, :code:`readLoginAck`).
  **These functions do not return until a message could be read**
  (or if a connection issue has been detected).
- :code:`sendString` and :code:`sendJson`,
  that respectively send a user-defined string
  or a user-defined JSON object on the network.
- :code:`recvString` and :code:`recvJson`,
  that respectively receive a string or a JSON object from the network.
  **These functions do not return until a message could be read**
  (or if a connection issue has been detected).

.. note::

  All these methods can throw exceptions if a network error has been encountered.
  Furthermore, all :code:`read<MESSAGE_TYPE>` methods will throw an exception if an unexpected
  message type has been received (*e.g.*, if the client received a :ref:`proto_KICK`).

Usage examples
~~~~~~~~~~~~~~

As an example, here is a basic player bot in Python.

.. code:: python

    try:
        # Instantiate a client in memory.
        client = Client()

        # Connect the internal socket to netorcai (on the 4242 port of the local machine).
        client.connect("localhost", 4242)

        # Log in to netorcai as a player. The client's nickname is "Example".
        client.send_login("Example", "player")
        client.read_login_ack()

        # Wait for the game to start.
        game_starts = client.read_game_starts()

        # Precalculation can be done here. Here, the initial game state is just printed.
        print(game_starts.initial_game_state)

        # For each turn.
        for i in range(game_starts.nb_turns_max):
            # Wait for the turn to start.
            turn = client.read_turn()
            # Decide what to do. Here, the current game state is just printed and no action is done.
            print(turn.game_state)
            actions = []
            # Send the decided actions to netorcai.
            client.send_turn_ack(turn.turn_number, [])
    except Exception as e:
        print(e)

All libraries have examples in the :code:`examples` directory of their
respective repository. Please refer to them for more examples.

Getting the libraries
~~~~~~~~~~~~~~~~~~~~~

Getting the latest released version is easy for languages that have a standard package index.

- D: Add the :code:`netorcai-client` dependency in your project (`netorcai-client package on DUB`_).
- Java: Not uploaded on the maven repository yet 😽.
- Python: :code:`pip install netorcai` (`netorcai package on PyPI`_)

Otherwise, getting the library from its git repository is pretty straightforward.
Building and installation instructions are in the README of each repository.

Alternatively, some of these libraries are packaged in Nix_ in the netorcaipkgs_ package repository.
Here are some commands to install the libraries.

.. code:: bash

    # Install the C++ client library.
    # Latest release
    nix-env -f https://github.com/netorcai/netorcaipkgs/archive/master.tar.gz -iA netorcai_client_cpp
    # Up-to-date (latest commit)
    nix-env -f https://github.com/netorcai/netorcaipkgs/archive/master.tar.gz -iA netorcai_client_cpp_dev

.. _netorcai organization github repository: https://github.com/netorcai/
.. _netorcaipkgs: https://github.com/netorcai/pkgs
.. _netorcai-client-cpp: https://github.com/netorcai/netorcai-client-cpp
.. _netorcai-client-d: https://github.com/netorcai/netorcai-client-d
.. _netorcai-client-fortran: https://github.com/netorcai/netorcai-client-fortran
.. _netorcai-client-java: https://github.com/netorcai/netorcai-client-java
.. _netorcai-client-python: https://github.com/netorcai/netorcai-client-python
.. _Nix: https://nixos.org/nix/
.. _bindings: https://en.wikipedia.org/wiki/Language_binding
.. _netorcai-client package on DUB: https://code.dlang.org/packages/netorcai-client
.. _netorcai package on PyPI: https://pypi.org/project/netorcai/
