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

function test_rest_api() {
    echo "Building REST API tests utility"
    if go build -o rest-api-tests tests/rest_api_tests.go; then
        echo "REST API tests build ok"
    else
        echo "Build failed"
        return 1
    fi
    sleep 1

    echo "Check if Mock service REST API is available"
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

function prepare_code_coverage() {
    echo "Preparing code coverage environment"
    mkdir coverage
    export GOCOVERDIR=coverage/
}

function code_coverage_report() {
    echo "Preparing code coverage report"
    go tool covdata merge -i=coverage/ -o=.
    go tool covdata textfmt -i=. -o=coverage.txt
    go tool cover -func=coverage.txt
}

echo -e "------------------------------------------------------------------------------------------------"

prepare_code_coverage
test_rest_api

# shut down server gracefully
curl -X PUT http://localhost:8080/api/insights-results-aggregator/v2/exit

code_coverage_report
exit $?
