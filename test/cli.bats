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

# Arguments test
