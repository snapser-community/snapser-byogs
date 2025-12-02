# Snapser Bedrock Minecraft Server (Agones-aware)

Minecraft: Bedrock Edition dedicated server container for Snapser's
UDP based Game Server Fleet Snap.

- **Edition:** Bedrock
- **Protocol:** UDP (19132)
- **Agones:** Integrated via REST SDK calls from `agones_wrapper.sh`
- **Modes:** Vanilla default; no plugins yet

## Setup

1. Download the official **Minecraft Bedrock Dedicated Server (Linux)** from
   the Minecraft website.
2. Save it in this directory as:

```bash
bedrock-server.zip
```

## Local
### Build
```bash
docker build -t your-registry/snapser-bedrock-agones:latest .
docker push your-registry/snapser-bedrock-agones:latest
```

### Run
```bash
docker run --rm -it \
  -e EULA=TRUE \
  -e SERVER_NAME="Local Bedrock Test" \
  -e AGONES_ENABLED=false \
  -p 19132:19132/udp \
  your-registry/snapser-bedrock-agones:latest
```



## Snapser
### Publish
- Once the setup is complete, run the following command to publish your game server
```bash
# Publish the game server
#  TODO: Replace path with path to your Dockerfile eg: `--path /Users/Ajinkya/Development/SnapserEngine/snapser-community/snapser-byogs/basic/minecraft-bedrock/.`
snapctl byogs publish --tag mc-1 --path $path
```

### Deploy a Fleet
- Create a Snapend with `auth` and `game-server-fleet` snaps.
- Once the Snapend is up go to the Snap Configuration tool and select Game Server Fleet in the left Nav.
- Click on the Create Fleet button.
  - Give your fleet a name, description
  - Select the `mc-1` tag from the list of available tags.
  - Add `19132` as the port for the fleet. Do note: when a fleet server comes up, it will be allocated a random port but you have to tell Snapser the internal port your container is going to listen on.
  - Let the region be the default and keep the default Allocation Settings for now.
  - For Hardware settings, add a min of 0.5 CPU and 2 GB memory
  - IMPORTANT: Add an environment variable `EULA=TRUE`. If you do not do this, the fleet server will not come up. Microsoft requires you to accept the EULA before you can run the server.
  - Click `Save`.
- In a few minutes your fleet should be ready. Go to its Control panel and click on `Allocate'
- You will now see a server with a public IP and port. You can connect to it using the Minecraft Bedrock client.

### Minecraft App
- Open up your Minecraft App on your windows or iOS machine.
- Go to `Play -> Servers - > Add Server`
- Enter the IP and port of the server you just allocated.
- IMPORTANT: Enter the IP and Port from the `Allocated Servers` section of the fleet control panel.
- Click `Save` and then `Join` to connect to the server.

## Future Enhancements
1. Currently, the server is not persistent. In the future we plan to enhance this example to show you how you can persist the world using Snapser Asset/Storage snaps.
