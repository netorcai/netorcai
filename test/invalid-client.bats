# setup is called before each test
setup() {
    # Kill any running netorcai
    killall netorcai 2>/dev/null || true

    # Run netorcai
    netorcai --json-logs 3>/dev/null &

    # Wait for netorcai socket to be opened
    while ! nc -z localhost 4242; do
        sleep 0.1 # wait for 0.1 seconds before checking again
    done
}

# teardown is called after each test
teardown() {
    killall netorcai 2>/dev/null || true
}

@test "invalid-login-notjson" {
    run ./login-notjson.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-no-messagetype" {
    run ./login-no-messagetype.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-no-role" {
    run ./login-no-role.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-no-nickname" {
    run ./login-no-nickname.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-role" {
    run ./login-bad-role.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-nickname-short" {
    run ./login-bad-nickname-short.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-nickname-long" {
    run ./login-bad-nickname-long.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-nickname-badchars" {
    run ./login-bad-nickname-badchars.py
    [ "${status}" -eq 0 ]
}
