# snapser-bedrock-agones/Dockerfile
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Dependencies for Bedrock + curl for Agones REST
RUN apt-get update && apt-get install -y \
  unzip \
  libcurl4-openssl-dev \
  ca-certificates \
  tzdata \
  curl \
  && rm -rf /var/lib/apt/lists/*

# Create non-root user and working dir
RUN useradd -m -d /bedrock bedrock
WORKDIR /bedrock

# You download this from the official Minecraft Bedrock Dedicated Server page
COPY bedrock-server.zip /tmp/bedrock-server.zip

# Unpack Bedrock server
RUN unzip /tmp/bedrock-server.zip -d /bedrock && \
  rm /tmp/bedrock-server.zip

# Default server.properties (can be overridden/persisted via volume)
COPY server.properties /bedrock/server.properties

# Entrypoint + Agones wrapper
COPY entrypoint.sh /entrypoint.sh
COPY agones_wrapper.sh /agones_wrapper.sh
RUN chmod +x /entrypoint.sh /agones_wrapper.sh && \
  chown bedrock:bedrock /entrypoint.sh /agones_wrapper.sh

# Default config via env
ENV EULA=FALSE \
  SERVER_NAME="Snapser Bedrock Server" \
  MAX_PLAYERS=10 \
  GAMEMODE=survival \
  DIFFICULTY=normal \
  AGONES_ENABLED=true

# Bedrock default UDP port
EXPOSE 19132/udp

USER bedrock

ENTRYPOINT ["/entrypoint.sh"]
