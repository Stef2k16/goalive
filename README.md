# Goalive
[![Go Report Card](https://goreportcard.com/badge/github.com/Stef2k16/goalive)](https://goreportcard.com/report/github.com/Stef2k16/gosleep)

**Goalive** is a simple tool to monitor health endpoints of your services.

## Functionality
**Goalive** allows to define a list of custom endpoints that should be polled periodically. These endpoints should be
dedicated health endpoints that return a 2xx HTTP status code if the service is running fine.
In case of problems, i.e. HTTP status codes != 2xx or connection issues, notifications can be send via Discord or 
Telegram. To avoid messages en masse for a failing endpoint, notifications are only
send for the first detected failure or if a previously failing endpoint has been fixed.

Alternatively, the most recent status can be requested manually. 
- For Discord, send `!status` to the channel with added bot
- For Telegram, send `/status` to the bot

## Setup
The setup consists of two steps:
- Set up of a Discord or Telegram bot
- Get **Goalive** and adjust the configuration to your needs

### Setup of Discord or Telegram Bot
To deliver notifications, you need to run your own bot.

#### Discord
For the configuration of **Goalive** with Discord, you need a bot token and a channel ID. The bot token is required to 
connect to your Bot. The channel ID determines the channel to deliver notifications to.

How to create a bot and receive a token is e.g. detailed [here](https://www.writebots.com/discord-bot-token/).
To get the channel ID, simply right-click the channel of choice and copy the ID.

#### Telegram
For the configuration of **Goalive** with Telegram, your need a bot token and your User ID.

To set up a new bot and receive a bot token, simply send `\newbot` to _BotFather_. To retrieve your user ID, you can for
example add the _userinfobot_ in Telegram and start it. The bot will then reply with your user ID.

### Get Goalive and Adjust the Configuration to Your Needs
Currently, Goalive can only be built from source. You thus need [Go](https://golang.org/dl/) (version >= 1.16).

To create an executable, clone the repository and run `go build`. The repo contains an example configuration
`example-config.yaml`. You at least have to add your bot token and channel ID (for Discord) / user ID (for Telegram)
to the configuration.
Additionally, you can define the URLs you want to monitor, the polling interval in seconds, and the location of a log 
file.

After your configuration is ready, you can start the monitoring with `./goalive --config path/to/your/config.yaml`
on Linux-based systems or `goalive.exe --config path/to/your/config.yaml` on Windows.