.. _changelog:

Changelog
=========

All notable changes to this project will be documented in this file.
The format is based on `Keep a Changelog`_.
netorcai adheres to `Semantic Versioning`_ and its public API includes the following.

- netorcai’s program command-line interface.
- netorcai’s metaprotocol.

........................................................................................................................

Unreleased
----------

- `Commits since v1.2.0 <https://github.com/netorcai/netorcai/compare/v1.2.0...master>`_

Changed (breaks metaprotocol)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

- The `CONTENT_SIZE` message header is now 32-bit (was 16-bit).
  `CONTENT_SIZE` value must be smaller than 1 Kio_ for the first message,
  and smaller than 16 Mio_ for other messages.
- Protocol version handshake added in :ref:`proto_LOGIN` and :ref:`proto_LOGIN_ACK`.
  As a result, old clients will not be able to log in anymore because their metaprotocol version is unknown.

Added
~~~~~

- New CLI command ``--fast``, which allows to start a new turn as soon as
  all players have decided what to do — instead of relying on a timer.
- Special players can now be connected to netorcai.
  The game logic knows which bots are *special*, which allows implementing
  game modes with asymetric game rules.
  As an example, the ghosts of a bomberman game could be implemented in a
  special player bot which could have different actions than bombermen.

  - New CLI command ``--nb-splayers-max``, to define the maximum number of special players.
  - :ref:`proto_DO_INIT` and :ref:`proto_GAME_STARTS` messages now contain a ``nb_special_players`` field.

Fixed
~~~~~

- Several issues around netorcai's termination have been fixed.

  - Kicking clients induced a non-compressible time delay to limit the loss of messages.
    This time delay has been removed **and** the last messages sent by netorcai should **not** be lost anymore.
  - Data races could occur in the sending of the last messages to clients.

........................................................................................................................

v1.2.0
------

- `Commits since v1.1.0 <https://github.com/netorcai/netorcai/compare/v1.1.0...v1.2.0>`_

Added
~~~~~

- New CLI command ``--autostart``,
  that automatically starts the game when all clients (and one game logic) are connected.
  The expected clients are those defined by ``--nb-players-max`` and ``--nb-visus-max``.

Changed
~~~~~~~

- Client libraries are now hosted on `netorcai's organization github repository`_.
- Documentation is now on `netorcai's readthedocs`_.

Fixed
~~~~~

- All players always remained connected in the ``players_info`` array of :ref:`proto_GAME_STARTS` and :ref:`proto_TURN` messages.
  Now, the ``is_connected`` field of disconnected players should be set to ``false``.

........................................................................................................................

v1.1.0
------

- `Commits since v1.0.1 <https://github.com/netorcai/netorcai/compare/v1.0.1...v1.1.0>`_
- Release date: 2018-10-29

Added
~~~~~

-  New CLI command ``--simple-prompt``, that forces the use of the basic prompt.

........................................................................................................................

v1.0.1
------

- `Commits since v1.0.0 <https://github.com/netorcai/netorcai/compare/v1.0.0...v1.0.1>`_
- Release date: 2018-10-23

Changed
~~~~~~~

-  The repository has moved to https://github.com/netorcai/netorcai.

........................................................................................................................

v1.0.0
------

- `Commits since v0.1.0 <https://github.com/netorcai/netorcai/compare/v0.1.0...v1.0.0>`_
- Release date: 2018-06-11

Added (program):
~~~~~~~~~~~~~~~~

- The metaprotocol is now fully implemented.
  netorcai is now heavily tested under continuous integration,
  all coverable code should now be covered.
- New ``--delay-turns`` command-line option to specify the minimum
  number of milliseconds between two consecutive turns.
- New interactive prompt.

Changed (metaprotocol):
~~~~~~~~~~~~~~~~~~~~~~~

- :ref:`proto_GAME_STARTS`

   - The ``data`` field has been renamed ``initial_game_state``.
   - ``player_id``: The “null” player_id is now represented as -1
     (was JSON's ``null``).
   - New ``milliseconds_between_turns`` field
     (minimum amount of milliseconds between two consecutive turns).
   - New ``players_info`` array used to forward information about the
     players to visualization clients.

- :ref:`proto_GAME_ENDS`

  - The ``data`` field has been renamed ``game_state``.
  - ``winner_player_id``: The “null” player_id is now represented as -1
    (was JSON's ``null``).

- :ref:`proto_TURN`

  - New ``players_info`` array used to forward information about the
    players to visualization clients.

- :ref:`proto_DO_TURN_ACK`

  - New ``winner_player_id`` field,
    which represents the current leader of the game (if any).

- The ``DO_FIRST_TURN`` message type has been renamed :ref:`proto_DO_INIT`
- New :ref:`proto_DO_INIT_ACK` message (game logic initialization).

Fixed:
~~~~~~

- Various fixes, as the metaprotocol was not implemented yet — and therefore not tested.

........................................................................................................................

v0.1.0
------

- First released version.
- Release date: 2018-05-02

.. _Keep a Changelog: http://keepachangelog.com/en/1.0.0/
.. _Semantic Versioning: http://semver.org/spec/v2.0.0.html
.. _netorcai's organization github repository: https://github.com/netorcai
.. _netorcai's readthedocs: https://netorcai.readthedocs.io
.. _Kio: https://en.wikipedia.org/wiki/Binary_prefix#Adoption_by_IEC,_NIST_and_ISO
.. _Mio: https://en.wikipedia.org/wiki/Binary_prefix#Adoption_by_IEC,_NIST_and_ISO
