unitsize(1cm);

real margin=1mm;
object logic = draw("game logic", box, (0,0), margin);
object broker = draw("netorcai", box, (0,-2), margin);

string[] client_names = {"player 1", "player 2", "visu 1"};

for (int i = 0; i < client_names.length; ++i)
{
    object o = draw(client_names[i], box, (-3+3*i,-4), 0);
    add(new void(picture pic, transform t)
    {
        draw(pic, point(broker,S,t)..point(o,N,t));
    });
}

add(new void(picture pic, transform t)
{
    draw(pic, point(logic,S,t)..point(broker,N,t));
});
