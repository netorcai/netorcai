load test_helper

# setup is called before each test
setup() {
    killall_netorcai
}

# teardown is called after each test
teardown() {
    killall_netorcai
}

@test "cli-noargs" {
    run_netorcai_wait_listening "" ""
    [ "$?" -eq 0 ]
}

@test "cli-help" {
    run netorcai --help
    [ "${status}" -eq 0 ]
}

@test "cli-h" {
    run netorcai -h
    [ "${status}" -eq 0 ]
}

@test "cli-version" {
    run netorcai --version
    [ "${status}" -eq 0 ]
    [ $(echo "${lines[0]}" | grep -o -E 'v[0-9]+\.[0-9]+\.[0-9]+.*') == "${lines[0]}" ]
}

# Verbosity tests
@test "cli-ok-verbose" {
    run_netorcai_wait_listening "" "" --verbose
    [ "$?" -eq 0 ]
}

@test "cli-ok-quiet" {
    run_netorcai_wait_listening "" "" --quiet
    [ "$?" -eq 0 ]
}

@test "cli-ok-debug" {
    run_netorcai_wait_listening "" "" --debug
    [ "$?" -eq 0 ]
}

@test "cli-ok-jsonlogs" {
    run_netorcai_wait_listening "" "" --json-logs
    [ "$?" -eq 0 ]
}

##################
# Arguments test #
##################
# --nb-players-max
@test "cli-nb-players-max-not-integer" {
    run netorcai --nb-players-max=meh
    [ "${status}" -ne 0 ]
}

@test "cli-nb-players-max-toosmall" {
    run netorcai --nb-players-max=0
    [ "${status}" -ne 0 ]
}

@test "cli-nb-players-max-small" {
    run_netorcai_wait_listening "" "" --nb-players-max=1
    [ "$?" -eq 0 ]
}

@test "cli-nb-players-max-big" {
    run_netorcai_wait_listening "" "" --nb-players-max=1024
    [ "$?" -eq 0 ]
}

@test "cli-nb-players-max-toobig" {
    run netorcai --nb-players-max=1025
    [ "${status}" -ne 0 ]
}

# --port
@test "cli-port-not-integer" {
    run netorcai --port=meh
    [ "${status}" -ne 0 ]
}

@test "cli-port-toosmall" {
    run netorcai --port=0
    [ "${status}" -ne 0 ]
}

@test "cli-port-small" {
    run_netorcai_wait_listening "" "1025" --port=1025
    [ "$?" -eq 0 ]
}

@test "cli-port-big" {
    run_netorcai_wait_listening "" "65535" --port=65535
    [ "$?" -eq 0 ]
}

@test "cli-port-toobig" {
    run netorcai --port=65536
    [ "${status}" -ne 0 ]
}

# --nb-turns-max
@test "cli-nb-turns-max-not-integer" {
    run netorcai --nb-turns-max=meh
    [ "${status}" -ne 0 ]
}

@test "cli-nb-turns-max-toosmall" {
    run netorcai --nb-turns-max=0
    [ "${status}" -ne 0 ]
}

@test "cli-nb-turns-max-small" {
    run_netorcai_wait_listening "" "" --nb-turns-max=1
    [ "$?" -eq 0 ]
}

@test "cli-nb-turns-max-big" {
    run_netorcai_wait_listening "" "" --nb-turns-max=65535
    [ "$?" -eq 0 ]
}

@test "cli-nb-turns-max-toobig" {
    run netorcai --nb-turns-max=65536
    [ "${status}" -ne 0 ]
}

# --nb-visus-max
@test "cli-nb-visus-max-not-integer" {
    run netorcai --nb-visus-max=meh
    [ "${status}" -ne 0 ]
}

@test "cli-nb-visus-max-toosmall" {
    run netorcai --nb-visus-max=-1
    [ "${status}" -ne 0 ]
}

@test "cli-nb-visus-max-small" {
    run_netorcai_wait_listening "" "" --nb-visus-max=0
    [ "$?" -eq 0 ]
}

@test "cli-nb-visus-max-big" {
    run_netorcai_wait_listening "" "" --nb-visus-max=1024
    [ "$?" -eq 0 ]
}

@test "cli-nb-visus-max-toobig" {
    run netorcai --nb-visus-max=1025
    [ "${status}" -ne 0 ]
}

# --delay-first-turn
@test "cli-delay-first-turn-not-float" {
    run netorcai --delay-first-turn=meh
    [ "${status}" -ne 0 ]
}

@test "cli-delay-first-turn-toosmall" {
    run netorcai --delay-first-turn=49.999
    [ "${status}" -ne 0 ]
}

@test "cli-delay-first-turn-small" {
    run_netorcai_wait_listening "" "" --delay-first-turn=50
    [ "$?" -eq 0 ]
}

@test "cli-delay-first-turn-big" {
    run_netorcai_wait_listening "" "" --delay-first-turn=10000
    [ "$?" -eq 0 ]
}

@test "cli-delay-first-turn-toobig" {
    run netorcai --delay-first-turn=10000.001
    [ "${status}" -ne 0 ]
}
