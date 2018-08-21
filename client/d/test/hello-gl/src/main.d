import std.json;
import std.format;
import std.stdio;

import netorcai;

void main()
{
    auto c = new Client;
    c.connect();
    scope(exit) c.close();

    write("Logging in as a game logic... "); stdout.flush();
    c.sendLogin("D-gl", "game logic");
    c.readLoginAck();
    writeln("done");

    write("Waiting for DO_INIT... "); stdout.flush();
    auto doInit = c.readDoInit();
    c.sendDoInitAck(`{"all_clients": {"gl": "D"}}`.parseJSON);
    writeln("done");

    foreach (turn; 0..doInit.nbTurnsMax)
    {
        write(format!"Waiting for DO_TURN %d... "(turn)); stdout.flush();
        auto doTurn = c.readDoTurn();
        c.sendDoTurnAck(`{"all_clients": {"gl": "D"}}`.parseJSON, -1);
        writeln("done");
    }
}
