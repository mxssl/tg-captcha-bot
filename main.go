package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Config struct for toml config file
type Config struct {
	ButtonText          string `mapstructure:"button_text"`
	WelcomeMessage      string `mapstructure:"welcome_message"`
	AfterSuccessMessage string `mapstructure:"after_success_message"`
	AfterFailMessage    string `mapstructure:"after_fail_message"`
	PrintSuccessAndFail string `mapstructure:"print_success_and_fail_messages_strategy"`
}

var config Config
var passedUsers = make(map[int]struct{})
var bot *tb.Bot

func init() {
	// Read config file
	err := readConfig()
	if err != nil {
		log.Fatalf("Cannot read config file. Error: %v", err)
	}
}

func main() {
	var err error
	bot, err = tb.NewBot(tb.Settings{
		Token:  os.Getenv("TGTOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Cannot start bot. Error: %v", err)
		return
	}

	bot.Handle(tb.OnUserJoined, challengeUser)
	bot.Handle(tb.OnCallback, passChallenge)

	log.Println("Bot started!")
	bot.Start()
}

func challengeUser(m *tb.Message) {
	if m.UserJoined.ID != m.Sender.ID {
		return
	}
	log.Printf("User: %v joined the chat: %v", m.UserJoined, m.Chat)
	newChatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: false}}
	bot.Restrict(m.Chat, &newChatMember)

	inlineKeys := [][]tb.InlineButton{{tb.InlineButton{
		Unique: "challenge_btn",
		Text:   config.ButtonText,
	}}}
	challengeMsg, _ := bot.Reply(m, config.WelcomeMessage, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})

	time.AfterFunc(30*time.Second, func() {
		_, passed := passedUsers[m.UserJoined.ID]
		if !passed {
			chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever()}
			bot.Ban(m.Chat, &chatMember)

			if config.PrintSuccessAndFail == "show" {
				bot.Edit(challengeMsg, config.AfterFailMessage)
			} else if config.PrintSuccessAndFail == "del" {
				bot.Delete(challengeMsg)
			}

			log.Printf("User: %v was banned in chat: %v", m.UserJoined, m.Chat)
		}
		delete(passedUsers, m.UserJoined.ID)
	})
}

// passChallenge is used when user passed the validation
func passChallenge(c *tb.Callback) {
	if c.Message.ReplyTo.Sender.ID != c.Sender.ID {
		bot.Respond(c, &tb.CallbackResponse{Text: "This button isn't for you"})
		return
	}
	passedUsers[c.Sender.ID] = struct{}{}

	if config.PrintSuccessAndFail == "show" {
		bot.Edit(c.Message, config.AfterSuccessMessage)
	} else if config.PrintSuccessAndFail == "del" {
		bot.Delete(c.Message)
	}

	log.Printf("User: %v passed the challenge in chat: %v", c.Sender, c.Message.Chat)
	newChatMember := tb.ChatMember{User: c.Sender, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: true}}
	bot.Promote(c.Message.Chat, &newChatMember)
	bot.Respond(c, &tb.CallbackResponse{Text: "Validation passed!"})
}

// readConfig is used for config unmarshall
func readConfig() (err error) {
	v := viper.New()
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
