#!/bin/bash

set -euo pipefail

cdir="$(dirname "$(readlink -f "${0}")")"

function sync {
	rsync -rvzza --progress --info=progress2 "${cdir}"/ "${node}:$(basename "${cdir}")"/
}

for node in node{1,2,3}.indiana.com; do
	sync "${node}"
done
