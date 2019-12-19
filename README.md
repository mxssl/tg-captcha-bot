[![Go Report Card](https://goreportcard.com/badge/github.com/mxssl/tg-captcha-bot)](https://goreportcard.com/report/github.com/mxssl/tg-captcha-bot)

# Telegram Captcha Bot

Telegram bot that validates new users that enter supergroup. Validation works like a simple captcha. Bot written in Go (Golang).

This bot has been tested on several supergroups (2000+ people) for a long time and has shown its effectiveness against spammers.

## Cloud hosted instance of the bot

[@cloud_tg_captcha_bot](https://t.me/cloud_tg_captcha_bot)

## How it works

1. Add the bot to your supergroup
2. Promote the bot for administrator privileges
3. A new user enters your supergroup
4. Bot restricts the user's ability to send messages
5. Bot shows a welcome message and a captcha button to the user
6. If the user doesn't press the button within 30 seconds then the user is banned by the bot

## If you want to run your own instance of the bot

- [Option 1 (the easiest one)](./INSTALL-1.md): docker-compose + already built docker container
- [Option 2](./INSTALL-2.md): docker-compose + build your own docker container
- [Option 3](./INSTALL-3.md): systemd

## Commands

`/healthz` - check that the bot is working correctly

## Сustomization

You can change several bot's settings through the configuration file `config.toml`

## Contacts

If you have questions feel free to ask me in TG [@mxssl](https://t.me/mxssl)
