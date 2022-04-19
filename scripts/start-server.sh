set -euo pipefail

if [[ $# -ne 3 ]]; then 
    SCRIPT_NAME=$(basename "${BASH_SOURCE}") 

    echo "Missing expected arguments"
    echo "  Example: ./scripts/start-server.sh \"Assessment API\" \"./assessment-api/cmd/server/server.go\" 8060"
    echo "  Example: ./scripts/start-server.sh \"Identity API\" \"./identity-api/cmd/server/server.go\" 8070"
    exit 1
fi

SERVER_NAME="${1}"
SERVER_COMMAND="${2}"
SERVER_PORT="${3}"

SCRIPT_DIR="$(git rev-parse --show-toplevel)"
cd "${SCRIPT_DIR}"

BUILD_NAME="${PPID}_build_result"

TEMP_DB_DIR="$(mktemp -d)"
function cleanup {
    # kill child processes when this script gets killed
    killall -c "${BUILD_NAME}" &> /dev/null || true
    kill "$(jobs -p)" &> /dev/null
    rm -rf "${TEMP_DB_DIR}"
}

trap 'cleanup' EXIT

BINARY="$(mktemp -d)/${BUILD_NAME}"

function restart_server {
    echo "halting running ${SERVER_NAME} GO server"
    killall -c "${BUILD_NAME}" &> /dev/null || true

    echo "removing old ${SERVER_NAME} GO binary"
    rm "${BINARY}" &> /dev/null || true

    if go build -o "${BINARY}" "${SERVER_COMMAND}"; then
        direnv exec ./ "$(command -v bash)" -c "PORT=${SERVER_PORT} ${BINARY}" &
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
