load test_helper

# setup is called before each test
setup() {
    killall_netorcai
}

# teardown is called after each test
teardown() {
    killall_netorcai
}

@test "scenario-login-player-ascii" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-ok-player-ascii.py
    [ "${status}" -eq 0 ]
}

@test "scenario-login-player-arabic" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-ok-player-arabic.py
    [ "${status}" -eq 0 ]
}

@test "scenario-login-player-japanese" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./login-ok-player-japanese.py
    [ "${status}" -eq 0 ]
}

@test "scenario-max-nb-players" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run ./max-nb-players.py
    [ "${status}" -eq 0 ]
}

@test "scenario-max-nb-visus" {
    run_netorcai_wait_listening "" "" --nb-visus-max=4
    [ $? -eq 0 ]

    run ./max-nb-visus.py
    [ "${status}" -eq 0 ]
}

@test "scenario-parallel-same-port" {
    run_netorcai_wait_listening "" ""
    [ $? -eq 0 ]

    run netorcai
    [ "${status}" -ne 0 ]
}

@test "scenario-parallel-different-port" {
    run_netorcai_wait_listening "" "4242" --port=4242
    [ $? -eq 0 ]

    run_netorcai_wait_listening "" "5151" --port=5151
    [ $? -eq 0 ]
}
