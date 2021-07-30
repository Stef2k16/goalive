# Goalive
[![Go Report Card](https://goreportcard.com/badge/github.com/Stef2k16/goalive)](https://goreportcard.com/report/github.com/Stef2k16/goalive)
[![Build & Tests](https://github.com/Stef2k16/goalive/actions/workflows/pipeline.yml/badge.svg)](https://github.com/Stef2k16/goalive/actions)

**Goalive** is a simple tool to monitor health endpoints of your services using modern notification
clients.

## Functionality
**Goalive** allows to define a list of custom endpoints that should be polled periodically. These endpoints should be
dedicated health endpoints that return a 2xx HTTP status code if the service is running fine.
In case of problems, i.e. HTTP status codes != 2xx or connection issues, notifications can be send via Discord, 
Telegram, or Slack. To avoid messages en masse for a failing endpoint, notifications are only
send for the first detected failure or if a previously failing endpoint has been fixed.

Alternatively, the most recent status can be requested manually. 
- For Discord, send `!status` to the channel with added bot
- For Telegram, send `/status` to the bot
- For Slack, send `/health` to the channel with the added bot

## Setup
The setup consists of two steps:
- Set up of a Discord or Telegram bot
- Get **Goalive** and adjust the configuration to your needs

### Setup of Discord, Telegram or Slack
To deliver notifications, you need to run your own bot or app.

#### Discord
For the configuration of **Goalive** with Discord, you need a bot token and a channel ID. The bot token is required to 
connect to your Bot. The channel ID determines the channel to deliver notifications to.

How to create a bot and receive a token is e.g. detailed [here](https://www.writebots.com/discord-bot-token/).
To get the channel ID, simply right-click the channel of choice and copy the ID.

#### Telegram
For the configuration of **Goalive** with Telegram, your need a bot token and your User ID.

To set up a new bot and receive a bot token, simply send `\newbot` to _BotFather_. To retrieve your user ID, you can for
example add the _userinfobot_ in Telegram and start it. The bot will then reply with your user ID.

#### Slack
The current implementation of the Slack notification client relies on 
[Socket Mode](https://api.slack.com/apis/connections/socket). It thus requires you to create your own
Slack Application to retrieve both an app token and a bot token. First create an app ([Slack API](https://api.slack.com/apps)).
In your app's settings under _Basic Information_ you can create an _App-Level Token_ which requires `connections:write` as scope.
Additionally, you have to enable _Socket Mode_ in the settings of the same name.

To retrieve a bot token, go to the _OAuth & Permissions_ tab in the app settings. The required
scopes are `chat:write` and `commands`.

Finally, install your app to any of your workspaces and add the created bot to a channel of your choice. You can copy
the channel ID in the channel's details.

If you need more information on how to create a Slack app, checkout, for example, [this](https://api.slack.com/authentication/basics#scopes)
official tutorial. For more information on the different token types, take a look [here](https://api.slack.com/authentication/token-types#bot).

### Get Goalive and Adjust the Configuration to Your Needs.
#### Download the Latest Release
You can download a pre-built executable for Windows (goalive.exe) or Linux (goalive) from 
the [release page](https://github.com/Stef2k16/goalive/releases).

#### Build from Source
If you need an executable for a different platform, Goalive can be built from source with
[Go](https://golang.org/dl/) (version >= 1.16).

To create an executable, clone the repository and run `go build` within the cloned repository. 

#### Create a Configuration
The repo contains an example configuration
`example-config.yaml`. You at least have to add your bot token and channel ID (for Discord) / user ID (for Telegram)
to the configuration. For Slack, you need an app token, a bot token and the channel ID.

Additionally, you can define the URLs you want to monitor, the polling interval in seconds, and the location of a log 
file.

After your configuration is ready, you can start the monitoring with `./goalive --config path/to/your/config.yaml`
on Linux-based systems or `goalive.exe --config path/to/your/config.yaml` on Windows.

## References
The project relies on fantastic API wrappers for Discord, Slack and Telegram:
- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [Slack API in Go](https://github.com/slack-go/slack)
- [Telebot](https://github.com/tucnak/telebot)
