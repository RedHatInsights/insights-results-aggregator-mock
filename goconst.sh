#!/bin/bash
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

BLUE=$(tput setaf 4)
RED_BG=$(tput setab 1)
GREEN_BG=$(tput setab 2)
NC=$(tput sgr0) # No Color

if ! [ -x "$(command -v goconst)" ]
then
    echo -e "${BLUE}Installing goconst${NC}"
    GO111MODULE=off go get github.com/jgautheron/goconst/cmd/goconst
fi

if [[ $(goconst -min-occurrences=4 ./... | tee /dev/tty | wc -l) -ne 0 ]]
then
    echo "${RED_BG}Duplicated string(s) found${NC}"
    exit 1
else
    echo "${GREEN_BG}No duplicated strings found${NC}"
    exit 0
fi
