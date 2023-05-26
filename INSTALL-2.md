# docker-compose: build your own docker container

## Prerequisites

- Obtain bot token from [@BotFather](https://t.me/BotFather)
- Install [Docker](https://docs.docker.com/install)

## Instructions

1. Clone the repo

```bash
git clone https://github.com/momai/tg-captcha-bot.git
cd tg-captcha-bot
```

2. remove ```docker-compose.yml``` and rename 
```bash
mv docker-compose.local.yml docker-compose..yml
```

3. To rename the file env.sample to .env, you can use the following command:
```mv .env.sample .env```
After renaming the file, open the newly created .env file in a text editor and add your Telegram bot token to it.


4. Build a Docker container

```bash
docker compose build
```

5. Run the container

```bash
docker compose up -d
```

6. Check that the bot started correctly

```bash
docker compose ps
docker compose logs
```

6. Add the bot to your supergroup and give it administrator privileges
