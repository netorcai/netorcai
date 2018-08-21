import std.json;
import std.format;
import std.stdio;

import netorcai;

void main()
{
    auto c = new Client;
    c.connect();

    write("Logging in as a game logic...");
    c.sendLogin("D-gl", "game logic");
    c.readLoginAck();
    writeln(" done");

    write("Waiting for DO_INIT...");
    auto doInit = c.readDoInit();
    c.sendDoInitAck(`{"all_clients": {"gl": "D"}}`.parseJSON);
    writeln(" done");

    foreach (turn; 0..doInit.nbTurnsMax)
    {
        write(format!"Waiting for DO_TURN %d..."(turn));
        auto doTurn = c.readDoTurn();
        c.sendDoTurnAck(`{"all_clients": {"gl": "D"}}`.parseJSON, -1);
        writeln(" done");
    }

    c.close();
}
