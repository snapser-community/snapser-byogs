#!/usr/bin/env bash
# snapser-bedrock-agones/agones_wrapper.sh
set -euo pipefail

cd /bedrock

SDK_PORT="${AGONES_SDK_HTTP_PORT:-9358}"
SDK_BASE="http://127.0.0.1:${SDK_PORT}"

log() {
  echo "[AGONES] $*"
}

wait_for_sdk() {
  local retries=60
  while (( retries > 0 )); do
    if curl -sSf "${SDK_BASE}/gameserver" >/dev/null 2>&1; then
      log "SDK is available on ${SDK_BASE}"
      return 0
    fi
    log "Waiting for SDK on ${SDK_BASE}..."
    sleep 2
    (( retries-- ))
  done
  log "SDK server never became available"
  return 1
}

send_ready() {
  log "Sending Ready() to Agones"
  curl -sSf -X POST -H "Content-Type: application/json" \
    -d '{}' "${SDK_BASE}/ready" >/dev/null
}

send_health() {
  curl -sSf -X POST -H "Content-Type: application/json" \
    -d '{}' "${SDK_BASE}/health" >/dev/null || \
    log "Health() call failed (will retry next tick)"
}

send_shutdown() {
  log "Sending Shutdown() to Agones"
  curl -sSf -X POST -H "Content-Type: application/json" \
    -d '{}' "${SDK_BASE}/shutdown" >/dev/null || \
    log "Shutdown() call failed"
}

start_bedrock() {
  log "Starting bedrock_server..."
  env LD_LIBRARY_PATH=. ./bedrock_server &
  BEDROCK_PID=$!
  log "bedrock_server PID=${BEDROCK_PID}"
}

wait_for_bedrock_ready() {
  # Simple heuristic: let it boot for a bit.
  # You could improve this by tailing logs or probing UDP if you want.
  log "Waiting for bedrock_server to finish startup..."
  sleep 10
}

main() {
  wait_for_sdk || exit 1
  start_bedrock
  wait_for_bedrock_ready
  send_ready

  # Health loop while child is running
  while kill -0 "${BEDROCK_PID}" 2>/dev/null; do
    send_health
    sleep 10
  done

  # Child exited
  wait "${BEDROCK_PID}" || true
  send_shutdown
}

main
