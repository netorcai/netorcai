Rationale
=========

In the context of `Lionel Martin's challenge`_,
I have been involved in the implementation of multiagent network
games meant to be played by bots.

After implementing several games (spaceships_ in 2014, `aquar.iom`_ in 2016) I
came to the following conclusions.

- Implenting the network server is tough.
- Handling the clients correctly (errors, fairness, not spamming slow clients...)
  mostly means that most of the development time is in the network game server,
  not in the game itself.
- The games in this context are quite specific
  (fair, turn-based, visualizable, no big performance constraint),
  which means the development effort can be shared regarding the network server.

.. _Lionel Martin's challenge: https://www.univ-orleans.fr/iut-orleans/informatique/intra/concours/
.. _spaceships: https://github.com/mpoquet/concoursiuto2015
.. _aquar.iom: https://github.com/mpoquet/aquar.iom
