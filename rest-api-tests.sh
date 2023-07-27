#!/usr/bin/env bash
# Copyright 2020 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COLORS_RED='\033[0;31m'
COLORS_RESET='\033[0m'
VERBOSE_OUTPUT=false

echo bash version is:
bash --version

if [[ $* == *verbose* ]] || [[ -n "${VERBOSE}" ]]; then
    # print all possible logs
    VERBOSE_OUTPUT=true
fi

function start_service() {
    if [ "$NO_SERVICE" = true ]; then
        echo "Not starting service"
        return
    fi

    echo "Starting a service"
    ./insights-results-aggregator-mock ||
        echo -e "${COLORS_RED}service exited with error${COLORS_RESET}" &
    # shellcheck disable=2181
    if [ $? -ne 0 ]; then
        echo "Could not start the service"
        exit 1
    fi
}

function test_rest_api() {
    start_service
    sleep 1

    echo "Building REST API tests utility"
    if go build -o rest-api-tests tests/rest_api_tests.go; then
        echo "REST API tests build ok"
    else
        echo "Build failed"
        return 1
    fi
    sleep 1

    curl http://localhost:8080/api/insights-results-aggregator/v2/ || {
        echo -e "${COLORS_RED}server is not running(for some reason)${COLORS_RESET}"
        exit 1
    }

    OUTPUT=$(./rest-api-tests 2>&1)
    EXIT_CODE=$?

    if [ "$VERBOSE_OUTPUT" = true ]; then
        echo "$OUTPUT"
    else
        echo "$OUTPUT" | grep -v -E "^Pass "
    fi

    return $EXIT_CODE
}

function cleanup() {
    print_descendent_pids() {
        pids=$(pgrep -P "$1")
        echo "$pids"
        for pid in $pids; do
            print_descendent_pids "$pid"
        done
    }

    echo Exiting and killing all children...

    children=$(print_descendent_pids $$)

    # disable the message when you send stop signal to child processes
    set +m

    for pid in $(echo -en "$children"); do
        # nicely asking a process to commit suicide
        if ! kill -PIPE "$pid" &>/dev/null; then
            # we even gave them plenty of time to think
            sleep 1
        fi
    done

    # restore the message back since we want to know that process wasn't stopped correctory
    # set -m

    for pid in $(echo -en "$children"); do
        # murdering those who're alive
        kill -9 "$pid" &>/dev/null
    done

    sleep 1
}
trap cleanup EXIT

echo -e "------------------------------------------------------------------------------------------------"

make build
test_rest_api
exit $?
