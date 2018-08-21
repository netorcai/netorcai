import std.json;
import std.format;
import std.stdio;

import netorcai;

void main()
{
    auto c = new Client;
    c.connect();

    write("Logging in as a player...");
    c.sendLogin("D-player", "player");
    c.readLoginAck();
    writeln(" done");

    write("Waiting for GAME_STARTS...");
    auto gameStarts = c.readGameStarts();
    writeln(" done");

    try
    {
        for (;;)
        {
            write("Waiting for TURN...");
            auto turn = c.readTurn();
            c.sendTurnAck(turn.turnNumber, `[{"player": "D"}]`.parseJSON);
            writeln(" done");
        }
    }
    catch(Exception e)
    {
        writeln("Caught!", e);
    }
}
