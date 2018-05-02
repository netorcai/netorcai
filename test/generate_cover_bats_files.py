#!/usr/bin/env python3
import re

#############
# Functions #
#############

def generate_bats_file(input_filename, output_filename):
    options_to_bypass = ['--help', '-h', '--version',
                         '--port', '--nb-turns-max', '--nb-players-max',
                         '--nb-visus-max', '--delay-first-turn',
                         '--verbose', '--quiet', '--debug', '--json-logs']

    with open(input_filename, "r") as in_file:
        content = [x.rstrip() for x in in_file.readlines()]

        with open(output_filename, "w") as out_file:
            count = 0
            for line in content:
                if 'run_netorcai_wait_listening "" "' in line:
                    line = re.sub(
                        r'run_netorcai_wait_listening "" ("\d*")',
                        r'run_netorcai_wait_listening "netorcai.cover" \1 '
                        '-test.coverprofile=${BATS_TEST_NAME}' + str(count) +
                        '.covout',
                        line)
                    count = count + 1
                if 'run netorcai' in line:
                    line = re.sub(
                        'run netorcai',
                        'run netorcai.cover '
                        '-test.coverprofile=${BATS_TEST_NAME}' + str(count) +
                        '.covout',
                        line)
                    count = count + 1
                elif '''[ "${status}" -ne 0 ]''' in line:
                    line = re.sub("""\[ "\${status}" -ne 0 \]""",
                                  '''[ "$status" -eq 0 ]''', line)

                for option in options_to_bypass:
                    line = re.sub(option + '\\b',
                                  '__bypass' + option, line)

                out_file.write("{}\n".format(line))

##########
# Script #
##########

# Input files definition
NETORCAI_FILES = ["invalid-client.bats",
                  "cli.bats"
                 ]

for robintest_file in NETORCAI_FILES:
    generate_bats_file(robintest_file, re.sub(
        ".bats", "-cover.bats", robintest_file))
