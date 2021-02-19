#!/bin/bash -xe

tgt="${1}"

# the log.json file has to exist as an array
[[ -f "${tgt}" ]] || echo "[]" >"${tgt}"
