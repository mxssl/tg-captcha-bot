# docker-compose: build your own docker container

## Prerequisites

1. Obtain bot token from [@BotFather](https://t.me/BotFather)
2. Install [Docker](https://docs.docker.com/install)
3. Install [Docker Compose](https://docs.docker.com/compose/install)

## Instructions

1. Clone the repo

```bash
git clone https://github.com/mxssl/tg-captcha-bot.git
cd tg-captcha-bot
```

2. Add a token from BotFather to env variable in docker-compose.yml

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

3. Build a Docker container

```bash
docker-compose build
```

4. Run the container

```bash
docker-compose up -d
```

5. Check that the bot started correctly

```bash
docker-compose ps
docker-compose logs
```

6. Add the bot to your supergroup and give it administrator privileges
