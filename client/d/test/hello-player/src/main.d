import std.json;
import std.format;
import std.stdio;

import netorcai;

void main()
{
    auto c = new Client;
    c.connect();
    scope(exit) c.close();

    write("Logging in as a player... "); stdout.flush();
    c.sendLogin("D-player", "player");
    c.readLoginAck();
    writeln("done");

    write("Waiting for GAME_STARTS... "); stdout.flush();
    auto gameStarts = c.readGameStarts();
    writeln("done");

    try
    {
        foreach (i; 1..gameStarts.nbTurnsMax)
        {
            write("Waiting for TURN... "); stdout.flush();
            auto turn = c.readTurn();
            c.sendTurnAck(turn.turnNumber, `[{"player": "D"}]`.parseJSON);
            writeln("done");
        }

        write("Waiting for GAME_ENDS..."); stdout.flush();
        auto gameEnds = c.readGameEnds();
        writeln("done");
    }
    catch(Exception e)
    {
        writeln("Failure: ", e);
    }
}
