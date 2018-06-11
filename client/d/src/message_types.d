import std.json;
import std.exception;

import json_utils;

/// Stores information about one player
struct PlayerInfo
{
    int playerID; /// The player unique identifier (in [0..nbPlayers[)
    string nickname; /// The player nickname
    string remoteAddress; /// The player socket remote address
    bool isConnected; /// Whether the player is currently connected or not
}

/// Parses a player information (in GAME_STARTS and GAME_ENDS messages)
PlayerInfo parsePlayerInfo(JSONValue o)
{
    PlayerInfo info;
    info.playerID = o["player_id"].getInt;
    info.nickname = o["nickname"].str;
    info.remoteAddress = o["remoteAddress"].str;
    info.isConnected = o["is_connected"].getBool;

    return info;
}

/// Parses several player information (in GAME_STARTS and GAME_ENDS messages)
PlayerInfo[] parsePlayersInfo(JSONValue[] array)
{
    PlayerInfo[] infos;
    infos.length = array.length;

    foreach (i, o; array)
    {
        infos[i] = o.parsePlayerInfo;
    }

    return infos;
}

/// Content of a LOGIN_ACK metaprotocol message
struct LoginAckMessage
{
    // ¯\_(ツ)_/¯
}

/// Content of a GAME_STARTS metaprotocol message
struct GameStartsMessage
{
    int playerID; /// Caller's player identifier. players: [0..nbPlayers[. visu: -1
    int nbPlayers; /// Number of players in the game
    int nbTurnsMax; /// Maximum number of turns. Game can finish before it
    double msBeforeFirstTurn; /// Time before the first TURN is sent (in ms)
    double msBetweenTurns; /// Time between two consecutive TURNs (in ms)
    PlayerInfo[] playersInfo; /// (only for visus) Information about the players
    JSONValue initialGameState; /// Game-dependent object.
}

/// Content of a GAME_ENDS metaprotocol message
struct GameEndsMessage
{
    int winnerPlayerID; /// Unique identifier of the player that won the game. Or -1.
    JSONValue gameState; /// Game-dependent object.
}

/// Content of a TURN metaprotocol message
struct TurnMessage
{
    int turnNumber; /// In [0..nbTurnsMax[
    PlayerInfo[] playersInfo; /// (only for visus) Information about the players
    JSONValue gameState; /// Game-dependent object.
}
