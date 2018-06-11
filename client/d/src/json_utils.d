import std.json, std.conv, std.format, std.algorithm, std.exception;

/// Reads a boolean from a JSONValue (handling several ways to store it)
bool getBool(in JSONValue v)
{
    if (v.type == JSON_TYPE.TRUE)
        return true;
    else if (v.type == JSON_TYPE.FALSE)
        return false;
    else if (v.type == JSON_TYPE.INTEGER)
    {
        enforce([0, 1].canFind(v.integer),
                format!"Cannot deduce boolean value from integer %d"(v.integer));
        return to!bool(v.integer);
    }
    else if (v.type == JSON_TYPE.UINTEGER)
    {
        enforce([0, 1].canFind(v.uinteger),
                format!"Cannot deduce boolean value from uinteger %d"(v.uinteger));
        return to!bool(v.uinteger);
    }
    else
    {
        enforce(0, "Cannot read bool value from JSONValue " ~ v.toString);
        return bool.init;
    }
}

unittest
{
    auto json = parseJSON(`{"a":true,
                            "b":false,
                            "c":0,
                            "d":1,
                            "e":2,
                            "f":42.51}`);

    assert(json["a"].getBool == true);
    assert(json["b"].getBool == false);
    assert(json["c"].getBool == false);
    assert(json["d"].getBool == true);

    auto e = collectException(json["e"].getBool);
    assert(e, "Should NOT be able to retrieve boolean value from JSON integer 2");

    e = collectException(json["f"].getBool);
    assert(e, "Should NOT be able to retrieve boolean value from JSON double 42.51");
}

/// Reads an integer from a JSONValue (handling several ways to store it)
int getInt(in JSONValue v)
{
    if (v.type == JSON_TYPE.INTEGER)
        return to!int(v.integer);
    else if (v.type == JSON_TYPE.UINTEGER)
        return to!int(v.uinteger);
    else
    {
        enforce(0, "Cannot read int value from JSONValue " ~ v.toString);
        return int.init;
    }
}

unittest
{
    auto json = parseJSON(`{"a":-10,
                            "b":37,
                            "c":42000,
                            "d":123456789,
                            "e":true,
                            "f":42.51}`);

    assert(json["a"].getInt == -10);
    assert(json["b"].getInt == 37);
    assert(json["c"].getInt == 42_000);
    assert(json["d"].getInt == 123_456_789);

    auto e = collectException(json["e"].getInt);
    assert(e, "Should NOT be able to retrieve integer value from JSON bool");

    e = collectException(json["f"].getInt);
    assert(e, "Should NOT be able to retrieve integer value from JSON double 42.51");
}

/// Reads a double from a JSONValue (handling several ways to store it)
double getDouble(in JSONValue v)
{
    if (v.type == JSON_TYPE.FLOAT)
        return v.floating;
    else if (v.type == JSON_TYPE.INTEGER)
        return to!double(v.integer);
    else if (v.type == JSON_TYPE.UINTEGER)
        return to!double(v.uinteger);
    else
    {
        enforce(0, "Cannot read double value from JSONValue " ~ v.toString);
        return double.init;
    }
}

unittest
{
    auto json = parseJSON(`{"a":-10,
                            "b":37,
                            "c":42000,
                            "d":123456789,
                            "e":true,
                            "f":42.51}`);

    assert(json["a"].getDouble == -10);
    assert(json["b"].getDouble == 37);
    assert(json["c"].getDouble == 42_000);
    assert(json["d"].getDouble == 123_456_789);

    immutable auto e = collectException(json["e"].getDouble);
    assert(e, "Should NOT be able to retrieve double value from JSON bool");

    assert(json["f"].getDouble == 42.51);
}
