#!/bin/bash -ex

export PATH=$PATH:/tmp/ramdisk/
export NOVENV=1

function prepare_venv() {
    # shellcheck disable=SC1091
    virtualenv -p python3 venv && source venv/bin/activate && python3 "$(command -v pip3)" install -r requirements.txt
}

[ "$NOVENV" == "1" ] || prepare_venv || exit 1

# shellcheck disable=SC2068
PYTHONDONTWRITEBYTECODE=1 python3 "$(command -v behave)" --tags=-skip -D dump_errors=true @feature_list.txt $@

