run_netorcai_wait_listening() {
    # Usage:
    #   run_netorcai_wait_listening
    #   run_netorcai_wait_listening netorcai
    #   run_netorcai_wait_listening netorcai 4242
    #   run_netorcai_wait_listening netorcai 4242 --json-logs --debug
    #   run_netorcai_wait_listening "" "" --json-logs --debug

    # Arguments
    netorcai_cmd=$1
    port=$2
    shift 2
    run_options=$@

    # Default values
    if [ -z ${netorcai_cmd} ]; then
        netorcai_cmd=netorcai
    fi

    if [ -z ${port} ]; then
        port=4242
    fi

    # Run netorcai
    ${netorcai_cmd} ${run_options} 3>/dev/null &

    # Wait for netorcai socket to be opened
    nb_tries=0
    while ! nc -z localhost ${port}; do
        sleep 0.1 # wait for 0.1 seconds before checking again
        ((nb_tries++))
        if [ ${nb_tries} -gt 50 ]; then
            echo "Netorcai did not open socket after waiting 5 s. Aborting."
            return 1
        fi
    done

    return 0
}
