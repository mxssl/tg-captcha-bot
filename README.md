# Telegram Captcha Bot

This telegram bot validates new users that enter supergroup. Validation works like a simple captcha.

## How it works
0. Promote bot for administrator privileges in your group
1. New user enter the supergroup
2. Bot restricts new user's ability to send messages
3. Bot show welcome message and captcha button to the user
4. Bot waits 30 seconds for the user to press the button
5. Bot bans the user if she/he didn't press the button within 30 seconds

## How to run
0. Obtain bot token from [@BotFather](https://t.me/BotFather)
1. Main method to run this bot is Docker container
2. Install [Docker](https://docs.docker.com/install)
3. Install [Docker Compose](https://docs.docker.com/compose/install)

#### Clone repo
```
git clone https://github.com/mxssl/tg-captcha-bot.git
cd tg-captcha-bot
```

#### Add token from BotFather to env variable in docker-compose.yml
```
version: '3'

services:
  tg-captcha-bot:
    build:
      context: .
      dockerfile: Dockerfile
    image: tg-captcha-bot:latest
    volumes:
      - ./config.toml:/config.toml
    environment:
      - TGTOKEN="your_token"
```

#### Build Docker container
```
docker-compose build
```

#### Run container
```
docker-compose up -d
```

#### Check that everything is OK
```
docker-compose ps
docker-compose logs
```

Add bot to your supergroup and give it administrator privileges.

#### Customize bot
You can change several bot's settings through the configuration file `config.toml`

## Contacts
If you have questions feel free to ask me [@mxssl](https://t.me/mxssl)
