# Telegram Captcha Bot

Telegram bot that validates new users that enter supergroup. Validation works like a simple captcha. Bot written in Go (Golang).

This bot has been tested on several large supergroups (1500+ people) for a long time and has shown its effectiveness against spammers.

## How it works
0. Add a bot to your supergroup
1. Promote the bot for administrator privileges
2. A new user enters your supergroup
3. Bot restricts the user's ability to send messages
4. Bot shows a welcome message and a captcha button to the user
5. If the user doesn't press the button within 30 seconds then the user is banned by the bot

## How to run
0. Obtain bot token from [@BotFather](https://t.me/BotFather)
1. The main method to run this bot is Docker container
2. Install [Docker](https://docs.docker.com/install)
3. Install [Docker Compose](https://docs.docker.com/compose/install)

## Instructions 
0. Clone the repo
```
git clone https://github.com/mxssl/tg-captcha-bot.git
cd tg-captcha-bot
```

1. Add a token from BotFather to env variable in docker-compose.yml
```yaml
version: '3'

services:
  tg-captcha-bot:
    build:
      context: .
      dockerfile: Dockerfile
    image: tg-captcha-bot:latest
    volumes:
      - ./config.toml:/config.toml
    restart: unless-stopped
    environment:
      - TGTOKEN=your_token
```

2. Build a Docker container
```
docker-compose build
```

3. Run the container
```
docker-compose up -d
```

4. Check that the bot started correctly
```
docker-compose ps
docker-compose logs
```

5. Add the bot to your supergroup and give it administrator privileges

## Сustomization
You can change several bot's settings through the configuration file `config.toml`

## Contacts
If you have questions feel free to ask me in TG [@mxssl](https://t.me/mxssl)
