# docker-compose: build your own docker container

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
      TGTOKEN: <your_telegram_bot_token_here>
```

Before starting, edit `config.toml` using the [configuration reference](README.md#configuration). Starting with v1.1.19, invalid durations prevent startup: `welcome_timeout` must be an integer from 1 to 9223372036 seconds, and `ban_duration` must be `"forever"` or an integer from 1 to 527040 minutes (366 days). Keep these TOML values quoted, for example `welcome_timeout = "30"` and `ban_duration = "10"`.

3. Build a Docker container

```bash
docker compose build
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
