package main

import (
        "context"
        "fmt"
        "log"
        "net"
        "net/http"
        "os"
        "os/signal"
        "regexp"
        "strconv"
        "sync"
        "syscall"
        "time"
        "math/rand"

        "github.com/pkg/errors"
        "github.com/spf13/viper"
        "golang.org/x/net/proxy"
        tb "gopkg.in/tucnak/telebot.v2"
)

// Config struct for toml config file
type Config struct {
        ButtonText          string `mapstructure:"button_text"`
        FakeButton          string `mapstructure:"ban_button_text"`
        WelcomeMessage      string `mapstructure:"welcome_message"`
        AfterSuccessMessage string `mapstructure:"after_success_message"`
        AfterFailMessage    string `mapstructure:"after_fail_message"`
        PrintSuccessAndFail string `mapstructure:"print_success_and_fail_messages_strategy"`
        WelcomeTimeout      string `mapstructure:"welcome_timeout"`
        BanDurations        string `mapstructure:"ban_duration"`
        FakeBanDurationMin  int64  `mapstructure:"fake_ban_duration_min"`
        UseSocks5Proxy      string `mapstructure:"use_socks5_proxy"`
        Socks5Address       string `mapstructure:"socks5_address"`
        Socks5Port          string `mapstructure:"socks5_port"`
        Socks5Login         string `mapstructure:"socks5_login"`
        Socks5Password      string `mapstructure:"socks5_password"`
        Captcha              string `mapstructure:"captcha"`
        CaptchaEnable        string `mapstructure:"captcha_enable"`
        AttackMode          string `mapstructure:"attack_mode"`
        AttackModeEnable    string `mapstructure:"attack_mode_enable"`
        CasEnable           string `mapstructure:"cas_enable"`
	CasBanDuration      string `mapstructure:"cas_ban_duration"`

} 

var config Config
var passedUsers = sync.Map{}
var bot *tb.Bot
var tgtoken = "TGTOKEN"
var configPath = "CONFIG_PATH"
var handledUsers = sync.Map{}
var botStates = sync.Map{}
var attackMode = sync.Map{}


func init() {
        err := readConfig()
        if err != nil {
                log.Fatalf("Cannot read config file. Error: %v", err)
        }
}

func main() {
        token, err := getToken(tgtoken)
        if err != nil {
                log.Fatalln(err)
        }
        log.Printf("Telegram Bot Token [%v] successfully obtained from env variable $TGTOKEN\n", token)

        var httpClient *http.Client
        if config.UseSocks5Proxy == "yes" {
                var err error
                httpClient, err = initSocks5Client()
                if err != nil {
                        log.Fatalln(err)
                }
        }

        bot, err = tb.NewBot(tb.Settings{
                Token:  token,
                Poller: &tb.LongPoller{Timeout: 10 * time.Second},
                Client: httpClient,
        })
        if err != nil {
                log.Fatalf("Cannot start bot. Error: %v\n", err)
        }

        bot.Handle(tb.OnUserJoined, challengeUser)

        bot.Handle("/captcha", func(m *tb.Message) {
        if !m.Private() {
                chatMember, err := bot.ChatMemberOf(m.Chat, m.Sender)
                if err != nil {
                        log.Println(err)
                        return
                }
                if chatMember.Role == tb.Creator || chatMember.Role == tb.Administrator {
                        currentState, _ := botStates.LoadOrStore(m.Chat.ID, true)
                        newState := !currentState.(bool)
                        botStates.Store(m.Chat.ID, newState)

             var status string
             if newState {
                 status = fmt.Sprintf("<b>enabled</b>\n%s", config.CaptchaEnable)
             } else {
                 status = "<b>disabled</b>"
             }
                        response := fmt.Sprintf(config.Captcha, status)
                        _, err := bot.Send(m.Chat, response, &tb.SendOptions{
                        ParseMode: tb.ModeHTML,
             })
                        if err != nil {
                                log.Println(err)
                                        }
                                }
                        }
                })



                helpMessage, err := readFileToString("help_message.txt")
if err != nil {
        log.Fatalf("Error reading help message file: %v", err)
}

bot.Handle("/start", func(m *tb.Message) {
        if m.Chat.Type == tb.ChatPrivate {
                _, err := bot.Send(m.Chat, helpMessage, tb.ModeMarkdown)
                if err != nil {
                        log.Println(err)
                }
        }
})

bot.Handle("/help", func(m *tb.Message) {
        if m.Chat.Type == tb.ChatPrivate {
                _, err := bot.Send(m.Chat, helpMessage, tb.ModeMarkdown)
                if err != nil {
                        log.Println(err)
                }
        }
})





        bot.Handle("/attack", func(m *tb.Message) {
     if !m.Private() {
         chatMember, err := bot.ChatMemberOf(m.Chat, m.Sender)
         if err != nil {
             log.Println(err)
             return
         }
         if chatMember.Role == tb.Creator || chatMember.Role == tb.Administrator {
             currentState, _ := attackMode.LoadOrStore(m.Chat.ID, false)
             newState := !currentState.(bool)
             attackMode.Store(m.Chat.ID, newState)

             var status string
             if newState {
                 status = fmt.Sprintf("<b>enabled</b>\n%s", config.AttackModeEnable)
             } else {
                 status = "<b>disabled</b>"
             }
             response := fmt.Sprintf(config.AttackMode, status)
             _, err := bot.Send(m.Chat, response, &tb.SendOptions{
                        ParseMode: tb.ModeHTML,
             })
             if err != nil {
                 log.Println(err)
             }
         }
     }
 })
        bot.Handle("/healthz", func(m *tb.Message) {
                msg := "I'm OK"
                if _, err := bot.Send(m.Chat, msg); err != nil {
                        log.Println(err)
                }
                log.Printf("Healthz request from user: %v\n in chat: %v", m.Sender, m.Chat)
        })

        log.Println("Bot started!")
        go func() {
                bot.Start()
        }()

        signalChan := make(chan os.Signal, 1)
        signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
        <-signalChan
        log.Println("Shutdown signal received, exiting...")
}
func shuffleButtons(buttons []tb.InlineButton) [][]tb.InlineButton {
    r := rand.New(rand.NewSource(time.Now().Unix()))

    humanIndex := r.Intn(2) + 1 // Guarantees that the index of the "I am human" button will be 1 or 2.
    shuffled := make([]tb.InlineButton, len(buttons))

    shuffled[humanIndex] = buttons[0]

    notHumanIndices := []int{0, 1, 2}
    notHumanIndices = append(notHumanIndices[:humanIndex], notHumanIndices[humanIndex+1:]...)

    for i, notHumanIndex := range notHumanIndices {
        shuffled[notHumanIndex] = buttons[i+1]
    }

    return [][]tb.InlineButton{shuffled}
}


func challengeUser(m *tb.Message) {
    // Проверяем, включен ли режим атаки
    attackEnabled, _ := attackMode.Load(m.Chat.ID)
    if attackEnabled != nil && attackEnabled.(bool) {
        banDuration := time.Now().Add(5 * time.Minute).Unix()
        chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: banDuration}
        err := bot.Ban(m.Chat, &chatMember)
        if err != nil {
            log.Println(err)
        }
        log.Printf("User: %v was banned in chat: %v for: 5 minutes (attack mode)", m.UserJoined, m.Chat)
        err = bot.Delete(m)
        if err != nil {
            log.Println(err)
        }
        return
    }

    // Проверяем, был ли пользователь уже забанен
    chatMember, err := bot.ChatMemberOf(m.Chat, m.UserJoined)
    if err != nil {
        log.Println(err)
        return
    }
    if chatMember.RestrictedUntil != 0 {
        // Пользователь уже забанен, пропускаем вывод капчи
        log.Printf("User: %v is already restricted in chat: %v", m.UserJoined, m.Chat)
	        if config.PrintSuccessAndFail == "del" {
                        err := bot.Delete(m)
                        if err != nil {
                            log.Println(err)
                        }
                    }
        return
    }

    // Проверяем, является ли пользователь инициатором
    if m.UserJoined.ID != m.Sender.ID {
        return
    }
    log.Printf("User: %v joined the chat: %v", m.UserJoined, m.Chat)

    // Проверяем, находится ли пользователь в списке забаненных CAS
    if config.CasEnable == "yes" {
        isBannedByCas, casStatus, err := checkUserCas(m.UserJoined.ID)
        if err != nil {
            log.Printf("Error checking user: %v with CAS in chat: %v, error: %v", m.UserJoined, m.Chat, err)
        } else {
            if casStatus == "" {
                log.Printf("User: %v was checked by CAS in chat: %v, user is not in CAS blacklist", m.UserJoined, m.Chat)
            } else {
                log.Printf("User: %v was checked by CAS in chat: %v, status: %v", m.UserJoined, m.Chat, casStatus)
            }
            if isBannedByCas {
                banDuration, e := getCasBanDuration()
                if e != nil {
                    log.Println(e)
                }
                chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: banDuration}
                err := bot.Ban(m.Chat, &chatMember)
                if err != nil {
                    log.Println(err)
                }
                log.Printf("User: %v was banned by CAS in chat: %v", m.UserJoined, m.Chat)
                return
            }
        }
    }

    // Проверяем, прошел ли пользователь уже проверку
    _, passed := passedUsers.Load(m.UserJoined.ID)
    if passed {
        return
    }

    // Применяем ограничения к пользователю
    newChatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: false}}
    err = bot.Restrict(m.Chat, &newChatMember)
    if err != nil {
        log.Println(err)
    }

    // Выводим капчу только для непрошедших пользователей
    if _, handled := handledUsers.Load(m.UserJoined.ID); !handled {
        challengeBtn := tb.InlineButton{
            Unique: "challenge_btn",
            Text:   config.ButtonText,
            Data:   "challenge_btn",
        }
        banBtn := tb.InlineButton{
            Unique: "ban_btn",
            Text:   config.FakeButton,
            Data:   "ban_btn",
        }
        banBtn2 := tb.InlineButton{
            Unique: "ban_btn_2",
            Text:   config.FakeButton,
            Data:   "ban_btn_2",
        }
        shuffledKeys := shuffleButtons([]tb.InlineButton{challengeBtn, banBtn, banBtn2})
        personalizedWelcomeMessage := fmt.Sprintf("%s, %s", m.UserJoined.FirstName, config.WelcomeMessage)
        challengeMsg, err := bot.Reply(m, personalizedWelcomeMessage, &tb.ReplyMarkup{InlineKeyboard: shuffledKeys})
        if err != nil {
            log.Printf("Can't send challenge message: %v", err)
            return
        }

        bot.Handle(&challengeBtn, passChallenge)
        bot.Handle(&banBtn, fakeChallenge)
        bot.Handle(&banBtn2, fakeChallenge)

        n, err := strconv.ParseInt(config.WelcomeTimeout, 10, 64)
        if err != nil {
            log.Println(err)
        }
        time.AfterFunc(time.Duration(n)*time.Second, func() {
            _, passed := passedUsers.Load(m.UserJoined.ID)
            if !passed {
                _, handled := handledUsers.Load(m.UserJoined.ID)
                if !handled {
                    banDuration, e := getBanDuration()
                    if e != nil {
                        log.Println(e)
                    }
                    chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: banDuration}
                    err := bot.Ban(m.Chat, &chatMember)
                    if err != nil {
                        log.Println(err)
                    }

                    if config.PrintSuccessAndFail == "show" {
                        _, err := bot.Edit(challengeMsg, config.AfterFailMessage)
                        if err != nil {
                            log.Println(err)
                        }
                    } else if config.PrintSuccessAndFail == "del" {
                        err := bot.Delete(m)
                        if err != nil {
                            log.Println(err)
                        }
                        err = bot.Delete(challengeMsg)
                        if err != nil {
                            log.Println(err)
                        }
                    }

                    log.Printf("User: %v was banned in chat: %v for: %v minutes", m.UserJoined, m.Chat, config.BanDurations)
                }
                handledUsers.Delete(m.UserJoined.ID)
            }
            passedUsers.Delete(m.UserJoined.ID)
        })
    }
}

func readFileToString(filePath string) (string, error) {
        content, err := os.ReadFile(filePath)
        if err != nil {
                return "", err
        }
        return string(content), nil
}


// passChallenge is used when user passed the validation
func passChallenge(c *tb.Callback) {
    if c.Message.ReplyTo.Sender.ID != c.Sender.ID {
        err := bot.Respond(c, &tb.CallbackResponse{Text: "This button isn't for you"})
        if err != nil {
            log.Println(err)
        }
        return
    }

    passedUsers.Store(c.Sender.ID, struct{}{})

    if config.PrintSuccessAndFail == "show" {
        _, err := bot.Edit(c.Message, config.AfterSuccessMessage)
        if err != nil {
            log.Println(err)
        }
    } else if config.PrintSuccessAndFail == "del" {
        err := bot.Delete(c.Message)
        if err != nil {
            log.Println(err)
        }
         err = bot.Delete(c.Message.ReplyTo)
         if err != nil {
            log.Println(err)
        }
    }

    log.Printf("User: %v passed the challenge in chat: %v", c.Sender, c.Message.Chat)
    newChatMember := tb.ChatMember{User: c.Sender, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: true}}
    err := bot.Promote(c.Message.Chat, &newChatMember)
    if err != nil {
        log.Println(err)
    }
    err = bot.Respond(c, &tb.CallbackResponse{Text: "Validation passed!"})
    if err != nil {
        log.Println(err)
    }
}




func fakeChallenge(c *tb.Callback) {
        if c.Message.ReplyTo.Sender.ID != c.Sender.ID {
                err := bot.Respond(c, &tb.CallbackResponse{Text: "This button isn't for you"})
                if err != nil {
                        log.Println(err)
                }
                return
        }


        banDuration := time.Now().Add(time.Duration(config.FakeBanDurationMin) * time.Minute).Unix()

        chatMember := tb.ChatMember{User: c.Sender, RestrictedUntil: banDuration}
        err := bot.Ban(c.Message.Chat, &chatMember)
        if err != nil {
                log.Println(err)
        }
        err = bot.Respond(c, &tb.CallbackResponse{Text: "Banned"})
        if err != nil {
                log.Println(err)
        }

        if config.PrintSuccessAndFail == "del" {
            err := bot.Delete(c.Message)
            if err != nil {
                log.Println(err)
            }
             err = bot.Delete(c.Message.ReplyTo)
             if err != nil {
                log.Println(err)
            }
        } else if config.PrintSuccessAndFail == "show" {
            _, err := bot.Edit(c.Message, config.AfterFailMessage)
            if err != nil {
                log.Println(err)
            }
        }

        handledUsers.Store(c.Sender.ID, struct{}{})

           log.Printf("User: %v was banned by fake button in chat: %v for: %v minutes", c.Sender, c.Message.Chat, config.FakeBanDurationMin)
}


func readConfig() (err error) {
        v := viper.New()
        path, ok := os.LookupEnv(configPath)
        if ok {
                v.SetConfigName("config")
                v.AddConfigPath(path)
        }
        v.SetConfigName("config")
        v.AddConfigPath(".")

        if err = v.ReadInConfig(); err != nil {
                return err
        }
        if err = v.Unmarshal(&config); err != nil {
                return err
        }
        return
}

func getToken(key string) (string, error) {
        token, ok := os.LookupEnv(key)
        if !ok {
                err := errors.Errorf("Env variable %v isn't set!", key)
                return "", err
        }
        match, err := regexp.MatchString(`^[0-9]+:.*$`, token)
        if err != nil {
                return "", err
        }
        if !match {
                err := errors.Errorf("Telegram Bot Token [%v] is incorrect. Token doesn't comply with regexp: `^[0-9]+:.*$`. Please, provide a correct Telegram Bot Token through env variable TGTOKEN", token)
                return "", err
        }
        return token, nil
}

func getBanDuration() (int64, error) {
        if config.BanDurations == "forever" {
                return tb.Forever(), nil
        }

        n, err := strconv.ParseInt(config.BanDurations, 10, 64)
        if err != nil {
                return 0, err
        }

        return time.Now().Add(time.Duration(n) * time.Minute).Unix(), nil
}


func getCasBanDuration() (int64, error) {
//    return time.Now().Add(time.Duration(config.CasBanDuration) * time.Minute).Unix(), nil

        if config.CasBanDuration == "forever" {
                return tb.Forever(), nil
        }

        n, err := strconv.ParseInt(config.CasBanDuration, 10, 64)
        if err != nil {
                return 0, err
        }

        return time.Now().Add(time.Duration(n) * time.Minute).Unix(), nil


}



func initSocks5Client() (*http.Client, error) {
        addr := fmt.Sprintf("%s:%s", config.Socks5Address, config.Socks5Port)
        dialer, err := proxy.SOCKS5("tcp", addr, &proxy.Auth{User: config.Socks5Login, Password: config.Socks5Password}, proxy.Direct)
        if err != nil {
                return nil, fmt.Errorf("cannot init socks5 proxy client dialer: %w", err)
        }

        httpTransport := &http.Transport{}
        httpClient := &http.Client{Transport: httpTransport}
        dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
                return dialer.Dial(network, address)
        }

        httpTransport.DialContext = dialContext

        return httpClient, nil
}
