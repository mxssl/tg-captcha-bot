<a href="https://github.com/mxssl/tg-captcha-bot/graphs/contributors"><img alt="GitHub contributors" src="https://img.shields.io/github/contributors/mxssl/tg-captcha-bot"></a>
<a href="https://github.com/mxssl/tg-captcha-bot/releases/latest"><img alt="GitHub All Releases" src="https://img.shields.io/github/downloads/mxssl/tg-captcha-bot/total"></a>
<a href="https://hub.docker.com/r/mxssl/tg-captcha-bot"><img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/mxssl/tg-captcha-bot"></a>

# Telegram Captcha Bot

Telegram bot that validates new users that enter supergroup. Validation works like a simple captcha. Bot written in Go (Golang).

The goal of this bot is to be as simple and minimal as possible.

This bot has been tested on several supergroups (2000+ people) for a long time and has shown its effectiveness against spammers.

## Cloud hosted instance of the bot

[@cloud_tg_captcha_bot](https://t.me/cloud_tg_captcha_bot)

## How it works

1. Add the bot to your supergroup
2. Promote the bot for administrator privileges
3. A new user enters your supergroup
4. Bot restricts the user's ability to send messages
5. Bot shows a welcome message and a captcha button to the user
6. If the user doesn't press the button within 30 seconds then the user is banned by the bot

## If you want to run your own instance of the bot

- [Option 1 (the easiest one)](./INSTALL-1.md): docker-compose + already built docker container
- [Option 2](./INSTALL-2.md): docker-compose + build your own docker container
- [Option 3](./INSTALL-3.md): systemd

## Commands

`/healthz` - check that the bot is working correctly

## Configuration

The bot uses a TOML configuration file and environment variables.

### Environment Variables

| Variable      | Required | Description                                                            |
| ------------- | -------- | ---------------------------------------------------------------------- |
| `TGTOKEN`     | Yes      | Telegram bot token from [@BotFather](https://t.me/BotFather)           |
| `CONFIG_PATH` | No       | Directory path containing `config.toml`. Defaults to current directory |

### Configuration File Options

Create a `config.toml` file with the following options:

| Option                        | Type   | Default                                             | Description                                                               |
| ----------------------------- | ------ | --------------------------------------------------- | ------------------------------------------------------------------------- |
| `button_text`                 | string | `"I'm not a robot!"`                                | Text displayed on the captcha button                                      |
| `welcome_message`             | string | `"Hello! This is the spam protection system..."`    | Message sent to new users                                                 |
| `after_success_message`       | string | `"User passed the validation."`                     | Message shown after successful verification                               |
| `after_fail_message`          | string | `"User didn't pass the validation and was banned."` | Message shown after failed verification                                   |
| `success_message_strategy`    | string | `"show"`                                            | Action after success: `"show"` (edit message) or `"del"` (delete message) |
| `fail_message_strategy`       | string | `"del"`                                             | Action after failure: `"show"` (edit message) or `"del"` (delete message) |
| `welcome_timeout`             | string | `"30"`                                              | Integer seconds from 1 to 9223372036 to press the button                    |
| `ban_duration`                | string | `"forever"`                                         | `"forever"` or integer minutes from 1 to 527040 (366 days)                  |
| `delete_join_message_on_fail` | string | `"yes"`                                             | Delete system join/leave messages for failed users: `"yes"` or `"no"`     |
| `use_socks5_proxy`            | string | `"no"`                                              | Enable SOCKS5 proxy: `"yes"` or `"no"`                                    |
| `socks5_address`              | string | `"1.1.1.1"`                                         | SOCKS5 proxy IP address                                                   |
| `socks5_port`                 | string | `"1080"`                                            | SOCKS5 proxy port                                                         |
| `socks5_login`                | string | `"login"`                                           | SOCKS5 proxy username                                                     |
| `socks5_password`             | string | `"password"`                                        | SOCKS5 proxy password                                                     |

Invalid, missing, zero, negative, or out-of-range durations prevent the bot from starting. Use `"forever"` explicitly for permanent bans.

### Example Configuration

```toml
button_text = "I'm not a robot!"
welcome_message = "Hello! Please press the button within 30 seconds or you will be banned!"
after_success_message = "User passed the validation."
after_fail_message = "User didn't pass the validation and was banned."
success_message_strategy = "show"
fail_message_strategy = "del"
welcome_timeout = "30"
ban_duration = "forever"
delete_join_message_on_fail = "no"
use_socks5_proxy = "no"
```

### Default behavior

With the default configuration:

1. New user joins the group
2. Bot restricts the user and sends a welcome message with "I'm not a robot!" button
3. User has 30 seconds to press the button
4. **If user passes**: Challenge message is edited to show "User passed the validation."
5. **If user fails**: User is banned forever, challenge message is deleted
