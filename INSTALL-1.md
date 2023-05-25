# docker-compose: use already built docker container image

## Prerequisites

- Obtain bot token from [@BotFather](https://t.me/BotFather)
- Install [Docker](https://docs.docker.com/install)

## Instructions

1. Clone the repo

```bash
git clone https://github.com/momai/tg-captcha-bot.git
cd tg-captcha-bot
```

2. Add a token from BotFather to env variable
```mv .env.sample .env```

3. Run the container

```bash
docker compose up -d
```

4. Check that the bot started correctly

```bash
docker compose ps
docker compose logs
```

6. Add the bot to your supergroup and give it administrator privileges
