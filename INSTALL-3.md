# systemd

## Prerequisites

Obtain bot token from [@BotFather](https://t.me/BotFather)

## Instructions

1. Clone the repo

```bash
git clone https://github.com/mxssl/tg-captcha-bot.git
cd tg-captcha-bot
```

2. Download the v1.1.19 bot binary and move it to the needed directory. The commands below are for Linux amd64; for Linux arm64, replace `linux_amd64` with `linux_arm64` in both archive names.

```bash
wget https://github.com/mxssl/tg-captcha-bot/releases/download/v1.1.19/tg-captcha-bot_1.1.19_linux_amd64.tar.gz

tar xvzf tg-captcha-bot_1.1.19_linux_amd64.tar.gz

mv tg-captcha-bot /usr/local/bin/tg-captcha-bot

chmod +x /usr/local/bin/tg-captcha-bot
```

3. Move bot's config to needed path

```bash
mkdir -p /etc/tg-captcha-bot
cp config.toml /etc/tg-captcha-bot/config.toml
```

Before starting, edit `/etc/tg-captcha-bot/config.toml` using the [configuration reference](README.md#configuration). Starting with v1.1.19, invalid durations prevent startup: `welcome_timeout` must be an integer from 1 to 9223372036 seconds, and `ban_duration` must be `"forever"` or an integer from 1 to 527040 minutes (366 days). Keep these TOML values quoted, for example `welcome_timeout = "30"` and `ban_duration = "10"`.

4. Create systemd unit file `/etc/systemd/system/tg-captcha-bot.service`

```bash
[Unit]
Description=tg-captcha-bot
Wants=network-online.target
After=network-online.target

[Service]
Environment="TGTOKEN=your_token"
Environment="CONFIG_PATH=/etc/tg-captcha-bot"
Type=simple
ExecStart=/usr/local/bin/tg-captcha-bot

Restart=always
RestartSec=3s

[Install]
WantedBy=multi-user.target
```

5. Reload configuration and restart service

```bash
systemctl daemon-reload
systemctl restart tg-captcha-bot.service
```

6. Check service status

```bash
systemctl status tg-captcha-bot.service
```

7. Check logs

```bash
journalctl -u tg-captcha-bot.service
```

8. Add the bot to your supergroup and give it administrator privileges
