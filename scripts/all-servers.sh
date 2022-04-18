set -euo pipefail

SCRIPT_DIR="$(git rev-parse --show-toplevel)"
cd "${SCRIPT_DIR}"

./scripts/assessment-server.sh &
./scripts/identity-server.sh &

wait < <(jobs -p)
