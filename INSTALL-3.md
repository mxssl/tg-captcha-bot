# systemd

## Prerequisites

Obtain bot token from [@BotFather](https://t.me/BotFather)

## Instructions

1. Clone the repo

```bash
git clone https://github.com/mxssl/tg-captcha-bot.git
cd tg-captcha-bot
```

2. Download bot binary and move it to needed directory

```bash
wget https://github.com/mxssl/tg-captcha-bot/releases/download/v1.1.6/tg-captcha-bot_1.1.6_linux_amd64.tar.gz

tar xvzf tg-captcha-bot_1.1.4_linux_amd64.tar.gz

mv tg-captcha-bot /usr/local/bin/tg-captcha-bot

chmod +x /usr/local/bin/tg-captcha-bot
```

3. Move bot's config to needed path

```bash
mkdir -p /etc/tg-captcha-bot
cp config.toml /etc/tg-captcha-bot/config.toml
```

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
