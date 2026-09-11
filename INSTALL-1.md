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

2. Replace `docker-compose.yml` with the following and set `TGTOKEN` to the token from BotFather

```yaml
services:
  tg-captcha-bot:
    image: mxssl/tg-captcha-bot:v1.1.19
    volumes:
      - ./config.toml:/config.toml
    restart: unless-stopped
    environment:
      TGTOKEN: <your_telegram_bot_token_here>
```

Before starting, edit `config.toml` using the [configuration reference](README.md#configuration). Starting with v1.1.19, invalid durations prevent startup: `welcome_timeout` must be an integer from 1 to 9223372036 seconds, and `ban_duration` must be `"forever"` or an integer from 1 to 527040 minutes (366 days). Keep these TOML values quoted, for example `welcome_timeout = "30"` and `ban_duration = "10"`.

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
