set -euo pipefail

SCRIPT_DIR="$(git rev-parse --show-toplevel)"
cd "${SCRIPT_DIR}"

ASSESSMENT_API_BUILD_NAME=assessment_api_build_result

TEMP_DB_DIR="$(mktemp -d)"
function cleanup {
    # kill child processes when this script gets killed
    killall -c "${ASSESSMENT_API_BUILD_NAME}" &> /dev/null || true
    kill "$(jobs -p)" &> /dev/null
    rm -rf "${TEMP_DB_DIR}"
}

trap 'cleanup' EXIT

ASSESSMENT_API_BINARY="$(mktemp -d)/${ASSESSMENT_API_BUILD_NAME}"

function restart_server {
    echo "halting running assessment GO server"
    killall -c "${ASSESSMENT_API_BUILD_NAME}" &> /dev/null || true

    echo "removing old assessment GO binary"
    rm "${ASSESSMENT_API_BINARY}" &> /dev/null || true

    if go build -o "${ASSESSMENT_API_BINARY}" ./assessment-api/cmd/server/server.go; then
        direnv exec ./ "$(command -v bash)" -c "PORT=8060 ${ASSESSMENT_API_BINARY}" &
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
