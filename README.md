# Lotus Jewel Client
Lotus Jewel Client is an unofficial tool that makes it easy to add [buttplug.io](https://buttplug.io/) support to Zero Ranger (and maybe other GameMaker games).

## Resources
[Latest release](https://github.com/Ripazhakgggdkp/lotus-jewel-client/releases)

## Prerequisites
This tool requires you to download and install [Intiface Central](https://intiface.com/central/) to command devices.

## Usage
This tool runs a server in port `25565` and relays HTTP messages to Intiface Central, a cross-platform hub that lets you easily connect any buttplug.io supported device. This makes it really simple to stick HTTP Calls via a tool like [UMT](https://github.com/krzys-h/UndertaleModTool) anywhere you want your device to do something. 

### API
`GET /connect/`

Connects the device to Intiface Central. Required before calling anything else.

`GET /vibrate/{intensity}`

Starts a vibration with the specified intensity (from 0 to 1). Will automatically stop after one second.

`GET /stop/`

Stops any vibration immediately. 

## Example UMT Mod
 
![imagen](https://user-images.githubusercontent.com/3671809/209412034-3187694d-a180-4f78-89b4-07c3ae8373b1.png)
 
`http_get("http://localhost:25565/connect")` will connect the game to Intiface Central.

Similarly, `http_get("http://localhost:25565/vibrate/0.3")` will command the first Intiface registered device to start vibrating at intensity 0.3.

## Limitations
Vibration cutoff duration is not configurable.
