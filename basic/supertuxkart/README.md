# BYOGS Super Tux Kart Tutorial

Super Tux Kart is a 3D open-source arcade racers with a variety of characters, tracks and modes to play. This SuperTuxKart example shows how to set up, deploy, and manage a SuperTuxKart game server on a Snapser Game Server Fleet.

This tutorial highlights how you can use the Snapser CLI tool and custom game server code and deploy it to Snapser global fleet.

## Application
- This code example is originally developed by the (Agones)[https://agones.dev/site/docs/examples/supertuxkart/] team. This example wraps the SuperTuxKart server with a Go binary, and introspects the log file to provide the event hooks for the SDK integration.

## Pre-Requisites

### A. Understanding Agones
Agones is an open source platform, for deploying, hosting, scaling, and orchestrating dedicated game servers for large scale multiplayer games, built on top of the industry standard, distributed system platform Kubernetes.

Snapser uses Agones under the hood and replaces bespoke or proprietary cluster management and game server scaling solutions with a managed solution - so that you can focus on the important aspects of building a multiplayer game, rather than developing the infrastructure to support it.

Built with both Cloud and on-premises infrastructure in mind, Snapser fleets can adjust its strategies as needed for Fleet management, autoscaling, and more to ensure the resources being used to host dedicated game servers are cost optimal for the environment that they are in.

<Note>
  Every Game server that is hosted on Snapser fleets needs to have the (Agones SDK)[https://0-8-0.agones.dev/site/docs/guides/client-sdks/] installed.
</Note>

### B. Snapctl Setup
You need to have a valid setup for Snapctl, which is Snapsers CLI tool. Please follow the step by step (tutorial)[https://snapser.com/docs/guides/tutorials/setup-snapctl] if you do not have Snapctl installed on your machine. You can run the following command to confirm you have a valid snapctl setup.

```bash
# Validate if your snapctl setup is correct
snapctl validate
```

### C. Docker
Make sure Docker engine is running on your machine. Open up Docker desktop and settings. Also, please make sure the setting **Use containerd for pulling and storing images** is **disabled**. You can find this setting in the Docker Desktop settings.

## Resources
All the files that are required by the Snapctl are under this folder
- **Dockerfile**: BYOSnap needs a Dockerfile. Snapser uses this file to containerize your application and deploy it.
- **server_config.xml**: Configuration file to start the server.


## Tutorial
### Step 1: Create your game server code
This tutorial already comes with an example of a running game and the required, game server code that works on UDP.

### Step 2 Build
- Build your server using Snapctl
```bash
# Deploy your game server to Snapser
# $tag = Tag should be in the format vX.Y.Z eg: v1.0.0
# $pathToCodeRoot = Path to the root of this code
snapctl byogs publish --tag $tag --path $pathToCodeRoot
```
<Note>
It should be noted that every subsequent `byogs publish` will need to have a different tag. Most studios just use semantic versioning as their tags.
</Note>

### Step 3: Create your cluster
#### Automated Setup
- Run `python snapend_create.py $companyId $gameId $tag` which first updates your `snapser-resources/snapser-snapend-manifest.json` file and then deploys it to Snapser via Snapctl.
- At the end, you will have a new Snapend running with an Auth & GSF Snaps with a preconfigured fleet. Keep a note of your `snapendId` as you will need this for the next stage.

#### Manual Setup
##### A. Create a Snapend
- Go to your Game on the Web portal.
- Click on **Create a Snapend**.
- Give your Snapend a name and hit Continue.
- Pick **Authentication** and **Game Server Fleets** snaps.
- Now keep hitting **Continue** till you reach the Review stage and then click **Snap it**.
- Your custom cluster should be up in about 2-4 minutes.
- Now, go into your Snapend and then click on the **Snapend Configuration**.
- Click on the **Game Server Fleet** Snap and then click on the **Fleets** configuration tool.
- Click on the **+ Fleet** button to start creating your first fleet.

##### B. Create a Game Server Fleet
1. Name: Give your fleet any name.
2. Description: Give your fleet any description.
3. Image Selector: Pick the tag that we just uploaded in the previous step.
4. Server Ports: Add **8080** as the Ingress port and leave Debug port as 0.
5. Configuration: This section allows you to configure setups for dev, stage and prod environments at once. But we will just select the environment tab that represents the current Snapend (Shown with a Checkmark)
  - Regions: Snapser supports multiple regions. Pick one or more regions from this list.
  - Allocation Settings: Here you can pick your min, max and buffer server counts. For this tutorial, lets stick to 1 for all three.
  - Hardware Settings: This is where you can pick the memory and CPU requirements for your server. For this tutorial, keep the defaults.
  - Optional Settings: This is where you can add custom commands, arguments and environment variables. Leave this as blank for the tutorial.
6. Scroll back up and hit save. This will tell Snapsers global compute to spin up a fleet and have it ready for you to allocate.

### Step 4: Allocate a Server
1. On the Fleet Configuration tool, you will see the status of your new fleet.
2. Once the servers are ready, you will see a green check next to the region chip. Once you see that, its time to see the Fleet Control plane.
3. Hover over, the fleet and you will see a menu, where you can click the **Control Plane** of the fleet.
4. This will take you to a view that shows you all the servers that are available to you for allocation.
5. For this tutorial, click the **+ Allocate** button on the top of the table and Snapser will assign a URL and port to your server.
6. You will see the server state change to **Allocated**.
7. Copy the **$url** and **$port** of this server as you will need it for testing.

It should be noted that, you do not have to manually allocate servers in Snapser. The Snapser fleet snap has integrations with Snapser snaps like Matchmaking, Parties and Lobbies. This allows you to automate your server allocation workflows.

### Step 5: Testing
- Download the SuperTuxKart (client)[https://supertuxkart.net/Download].
- Once its done downloading, start the app and pick **Online**.
- Then click on **Enter server address**.
- This is where you should enter the $url:$port. Eg: ec2-44-244-0-103.us-west-2.compute.amazonaws.com:7023 (You may have to manually enter this string in the Tux Kart game. Tux Kart game client, does not support clipboard paste on certain platforms.)
- You will now be taken to a lobby. You can ask some of your friends to download this game client and actually play the game with you.
- Once everyone is in (or even if its just you), you can now click on Start race.
- There you have it you are actually playing Super Tux Kart on a game server orchestrated by Snapser fleet.

It should be noted, that as soon as the game ends, Snapser fleet will receive a command from the Agones SDK telling it to shut down the server.
