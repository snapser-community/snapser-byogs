# Simple Unity Example

This is a very simple "unity server" that doesn't do much other than show how the SDK works in Unity.

## Prerequisites
### A. Install Unity
This example is working on
```
Unity Editor: Unity 2022.3.2f1 or later
OS: Windows 10 Pro or MacOS
```

### B. Install Unity modules
1. Make sure your Unity version has `Linux Build Support (Mono)` installed.
1. Make sure your Unity version has `Linux Dedicated Server Build Support` installed.
  <Note>
    You can install both using the Unity Hub.
    - Open Unity Hub
    - Go to the Installs tab
    - Click the three dots ⋯ next to your Unity version
    - Click Add Modules
  </Note>

### C. Build Server Settings
1. Make sure you have picked `Dedicated Server` as the platform and the target as `Linux`.


### Building a Server
* Open this folder with UnityEditor.
* Click on the `Build Tool/Build Server` menu item in the menu bar.
  * The Builds are created in a `Builds/Server` Folder.

### Publishing the server
```bash
snapctl byogs publish --tag $tag --path $codePath
```

## Calling an SDK API via a Client

### Building a Client
* Open this folder with UnityEditor.
* Click on the `Build Tool/Build Client` menu item in the menu bar.

### How to use a Client
* Run `Builds/Client/UnitySimpleClient.exe`.
* Set `Address` and `Port` text fields to GameServer's one. You can see these with the following command.
    ```
    $ kubectl get gs
    NAME                        STATE   ADDRESS         PORT   NODE       AGE
    unity-simple-server-z7nln   Ready   192.168.*.*     7854   node-name  1m
    ```
* Click on the `Change Server` Button.
* Set any text to a center text filed and click the `Send` button.
  * The Client will send the text to the Server.

  When a Server receives a text, it will send back "Echo : $text" as an echo.
  And an SDK API will be executed in a Server by the following rules.

    | Sending Text(Client) | SDK API(Server) |
    | ---- | ---- |
    | Allocate | Allocate() |
    | Label $1 $2 | SetLabel($1, $2) |
    | Annotation $1 $2 | SetAnnotation($1, $2) |
    | Shutdown | Shutdown() |
    | GameServer | GameServer() |
