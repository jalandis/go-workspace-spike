set -euo pipefail

SCRIPT_DIR="$(git rev-parse --show-toplevel)"
cd "${SCRIPT_DIR}"

./scripts/start-server.sh "Assessment API" "./assessment-api/cmd/server/server.go" 8060 &
./scripts/start-server.sh "Identity API" "./identity-api/cmd/server/server.go" 8070 &

wait < <(jobs -p)
