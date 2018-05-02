load test_helper

# setup is called before each test
setup() {
    killall_netorcai
}

# teardown is called after each test
teardown() {
    killall_netorcai
}

@test "invalid-login-notjson" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-notjson.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-no-messagetype" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-no-messagetype.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-no-role" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-no-role.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-no-nickname" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-no-nickname.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-role" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-bad-role.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-nickname-short" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-bad-nickname-short.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-nickname-long" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-bad-nickname-long.py
    [ "${status}" -eq 0 ]
}

@test "invalid-login-bad-nickname-badchars" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-bad-nickname-badchars.py
    [ "${status}" -eq 0 ]
}
