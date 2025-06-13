# docker-compose: use already built docker container image

## Prerequisites

- Obtain bot token from [@BotFather](https://t.me/BotFather)
- Install [Docker](https://docs.docker.com/install)

## Instructions

1. Clone the repo

```bash
git clone https://github.com/mxssl/tg-captcha-bot.git
cd tg-captcha-bot
```

2. Add a token from BotFather to env variable in docker-compose.yml

```yaml
version: "3"

services:
  tg-captcha-bot:
    image: mxssl/tg-captcha-bot:v1.1.13
    volumes:
      - ./config.toml:/config.toml
    restart: unless-stopped
    environment:
      TGTOKEN: <your_telegram_bot_token_here>
```

3. Pull the container

```bash
docker compose pull
```

4. Run the container

```bash
docker compose up -d
```

5. Check that the bot started correctly

```bash
docker compose ps
docker compose logs
```

6. Add the bot to your supergroup and give it administrator privileges
