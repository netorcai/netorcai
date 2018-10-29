# Changelog
All notable changes to this project will be documented in this file.  
The format is based on [Keep a Changelog][changelog].

netorcai adheres to [Semantic Versioning][semver].  
Its public API includes:
- netorcai's program command-line interface.
- netorcai's metaprotocol.

[//]: =========================================================================
## [Unreleased]

[//]: =========================================================================
## [1.1.0] - 2018-10-29
### Added
- New CLI command `--simple-prompt`, that forces the use of the basic prompt.

[//]: =========================================================================
## [1.0.1] - 2018-10-23
### Changed
- The repository has moved to https://github.com/netorcai/netorcai.

[//]: =========================================================================
## [1.0.0] - 2018-06-11
### Added (program):
- The metaprotocol is now fully implemented.  
  netorcai is now heavily tested under continuous integration,
  all coverable code should now be covered.
- New `--delay-turns` command-line option to specify
  the minimum number of milliseconds between two consecutive turns.
- New interactive prompt.

### Changed (metaprotocol):
- `GAME_STARTS`:
  - The `data` field has been renamed `initial_game_state`.
  - `player_id`: the "null" player_id is now represented as -1 (was `null`).
  - New `milliseconds_between_turns` field (minimum amount of milliseconds
    between two consecutive turns).
  - New `players_info` array used to forward information about the
    players to visualization clients.
- `GAME_ENDS`:
  - - The `data` field has been renamed `game_state`.
  - `winner_player_id`: the "null" player_id is now represented as -1
    (was `null`).
- `TURN`:
  - New `players_info` array used to forward information about the
    players to visualization clients.
- The `DO_FIRST_TURN` message type has been renamed `DO_INIT`
- New `DO_INIT_ACK` message (game logic initialization).
- `DO_TURN_ACK` : New `winner_player_id` field, which represents the current
  leader of the game (if any).

### Fixed:
- Various fixes, as the metaprotocol was not implemented yet — and therefore
  not tested.

[//]: =========================================================================
## 0.1.0 - 2018-05-02
- First released version.

[//]: =========================================================================
[changelog]: http://keepachangelog.com/en/1.0.0/
[semver]: http://semver.org/spec/v2.0.0.html

[Unreleased]: https://github.com/netorcai/netorcai/compare/v1.1.0...master
[1.1.0]: https://github.com/netorcai/netorcai/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/netorcai/netorcai/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/netorcai/netorcai/compare/v0.1.0...v1.0.0
