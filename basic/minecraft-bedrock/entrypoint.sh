#!/usr/bin/env bash
# snapser-bedrock-agones/entrypoint.sh
set -euo pipefail

cd /bedrock

# Gate on EULA, even though Bedrock doesn't use eula.txt like Java
if [[ "${EULA:-FALSE}" != "TRUE" && "${EULA:-FALSE}" != "true" ]]; then
  echo "You must set EULA=TRUE to run this server (see README)."
  exit 1
fi

# Ensure server.properties exists
if [[ ! -f server.properties ]]; then
  echo "server.properties not found, creating a minimal one..."
  cat <<EOF > server.properties
server-name=Snapser Bedrock Server
gamemode=survival
difficulty=normal
allow-cheats=false
max-players=10
online-mode=true
white-list=false
view-distance=32
tick-distance=4
player-idle-timeout=30
server-port=19132
server-portv6=19133
level-name=world
EOF
fi

# Helper to set/update a key in server.properties
set_prop() {
  local key="$1"
  local value="$2"

  if grep -qE "^${key}=" server.properties; then
    sed -i "s|^${key}=.*|${key}=${value}|" server.properties
  else
    echo "${key}=${value}" >> server.properties
  fi
}

# Apply env-driven config overrides
set_prop "server-name" "${SERVER_NAME:-Snapser Bedrock Server}"
set_prop "max-players" "${MAX_PLAYERS:-10}"
set_prop "gamemode" "${GAMEMODE:-survival}"
set_prop "difficulty" "${DIFFICULTY:-normal}"
set_prop "server-port" "19132"
set_prop "server-portv6" "19133"

echo "====================================================="
echo " Minecraft Bedrock Dedicated Server"
echo "  Name       : $(grep '^server-name=' server.properties | cut -d= -f2-)"
echo "  Gamemode   : $(grep '^gamemode=' server.properties | cut -d= -f2-)"
echo "  Difficulty : $(grep '^difficulty=' server.properties | cut -d= -f2-)"
echo "  Max players: $(grep '^max-players=' server.properties | cut -d= -f2-)"
echo "  Port (UDP) : 19132"
echo "  Agones     : ${AGONES_ENABLED:-true}"
echo "====================================================="

if [[ "${AGONES_ENABLED:-true}" == "true" || "${AGONES_ENABLED:-true}" == "TRUE" ]]; then
  exec /agones_wrapper.sh
else
  echo "Starting bedrock_server without Agones integration..."
  exec env LD_LIBRARY_PATH=. ./bedrock_server
fi
