# Snapser - Game Server Fleet Example - Go

A very simple "toy" game server created to demo and test running a UDP and/or
TCP game server on Snapser.

## Pre-requisite
Snapser uses Agones under the hood. This game server has the Agones SDK integrated. To learn how to deploy your edited version of go server to gcp, please check out this link: [Edit Your First Game Server (Go)](https://agones.dev/site/docs/getting-started/edit-first-gameserver-go/), or also look at the [`Makefile`](./Makefile).

## Setup
To create Game servers in Snapser you need to first add the **Game Server Fleet** control plane to your Snapser infrastructure.
1. Go to **Create Snapend**.
2. Under snaps pick the ones you want but also select the **Game Server Fleet** snap. Then proceed to the final step and hit **Snap it**.
3. Once your Snapend comes up, click on the cluster which will take you to your Snapend home page.
4. Under Quick links you will see the **Snaps Configuration** tool. Clicking on it, will take you to admin tools for all your Snaps.
5. Here click on **Game Server Fleets** in the left Nav bar and then you will see a button to **Add a fleet**.

Note: This is where you will create your Snapser hosted game server fleet. But before that you will need to upload your game server image. So proceed now to the Configuration Step 1.

## Configuration
## Step 1 - Upload your game server image
1. You will need to download the Snapser CLI tool for this and have docker running locally.
2. You will then use the CLI tool to upload this code to your own game server private image repository on Snapser.
```
snapctl byogs publish --tag "v0.0.1" --path <path_to_root_of_this_repo>
```
3. Once your code is uploaded go back to the web browser to move to Step 2.

## Step 2 - Create a fleet
1. Now coming back to the Admin tool for Fleet creation, give your fleet a name and a description.
2. Under Game Server Image, pick the tag that you just uploaded. If you are following these instructions it will be **v0.0.1**.
3. Next, setup dev, staging and prod settings for your BYOSnap. You will
at least need to setup one of the three. We recommend you to add the dev settings to start.
4. Here you will select, Max and buffer server counts for your fleet, CPU, Memory, which you can keep as defaults. Additionally, select **7654** as the Ingress Port and then hit Create Fleet.

Note: The fleet may take a minute or two to come up.

## Allocating Servers
You can allocate game servers via
1. Manually in the Web app, via the Fleet control plane.
2. The Matchmaker snap.
3. The Lobbies snap.

## Interacting with the Fleet

Once you have a game server up. You will receive a public URL and a port. With this you can start interacting with your
game server. Presently, Snapser game servers only support UDP.

When the server receives a text packet, it will send back "ACK:<text content>"
for UDP as an echo or "ACK TCP:<text content>" for TCP.

There are some text commands you can send the server to affect its behavior:

| Command             | Behavior                                                                                 |
| ------------------- | ---------------------------------------------------------------------------------------- |
| "EXIT"              | Causes the game server to exit cleanly calling `os.Exit(0)`                              |
| "UNHEATHY"          | Stopping sending health checks                                                           |
| "GAMESERVER"        | Sends back the game server name                                                          |
| "READY"             | Marks the server as Ready                                                                |
| "ALLOCATE"          | Allocates the game server                                                                |
| "RESERVE"           | Reserves the game server after the specified duration                                    |
| "WATCH"             | Instructs the game server to log changes to the resource                                 |
| "LABEL"             | Sets the specified label on the game server resource                                     |
| "CRASH"             | Causes the game server to exit / crash immediately                                       |
| "ANNOTATION"        | Sets the specified annotation on the game server resource                                |
| "PLAYER_CAPACITY"   | With one argument, gets the player capacity; with two arguments sets the player capacity |
| "PLAYER_CONNECT"    | Connects the specified player to the game server                                         |
| "PLAYER_DISCONNECT" | Disconnects the specified player from the game server                                    |
| "PLAYER_CONNECTED"  | Returns true/false depending on whether the specified player is connected                |
| "GET_PLAYERS"       | Returns a list of the connected players                                                  |
| "PLAYER_COUNT"      | Returns a count of the connected players                                                 |
