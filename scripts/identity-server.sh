set -euo pipefail

SCRIPT_DIR="$(git rev-parse --show-toplevel)"
cd "${SCRIPT_DIR}"

IDENTITY_API_BUILD_NAME=identity_api_build_result

TEMP_DB_DIR="$(mktemp -d)"
function cleanup {
    # kill child processes when this script gets killed
    killall -c "${IDENTITY_API_BUILD_NAME}" &> /dev/null || true
    kill "$(jobs -p)" &> /dev/null
    rm -rf "${TEMP_DB_DIR}"
}

trap 'cleanup' EXIT

IDENTITY_API_BINARY="$(mktemp -d)/${IDENTITY_API_BUILD_NAME}"

function restart_server {
    echo "halting running identity GO server"
    killall -c "${IDENTITY_API_BUILD_NAME}" &> /dev/null || true

    echo "removing old identity GO binary"
    rm "${IDENTITY_API_BINARY}" &> /dev/null || true

    if go build -o "${IDENTITY_API_BINARY}" ./identity-api/cmd/server/server.go; then
        direnv exec ./ "$(command -v bash)" -c "PORT=8070 ${IDENTITY_API_BINARY}" &
    fi
}

restart_server

# 0.5 second latency between builds of go
fswatch -l 0.5 -o \
    -e ".*" \
    -i "^.*[.]go$" \
    -i "^.*[.]env$" \
    -i "^.*[.]env_secret$" \
    -i "^.*[.]tmpl$" \
    -i "^.*[.]html$" \
    . | while read -r; do

    restart_server
done
