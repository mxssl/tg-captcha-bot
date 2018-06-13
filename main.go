package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
)

var passedUsers = make(map[int]struct{})
var bot *tb.Bot

func main() {
	var err error
	bot, err = tb.NewBot(tb.Settings{
		Token:  os.Getenv("TGTOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle(tb.OnUserJoined, challengeUser)
	bot.Handle(tb.OnCallback, passChallenge)

	log.Print("Started listening")
	bot.Start()
}

func challengeUser(m *tb.Message) {
	if m.UserJoined.ID != m.Sender.ID {
		return //invited users are not challengeable
	}
	log.Printf("%v joined the %v chat", m.UserJoined, m.Chat)
	newChatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: false}}
	bot.Restrict(m.Chat, &newChatMember)

	inlineKeys := [][]tb.InlineButton{{tb.InlineButton{
		Unique: "challenge_btn",
		Text:   "Я не спамер!",
	}}}
	challengeMsg, _ := bot.Reply(m, "Это защита от спама. У вас есть 30 секунд, чтобы нажать на кнопку, иначе вы будете забанены!", &tb.ReplyMarkup{InlineKeyboard: inlineKeys})

	time.AfterFunc(30*time.Second, func() {
		_, passed := passedUsers[m.UserJoined.ID]
		if !passed {
			chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever()}
			bot.Ban(m.Chat, &chatMember)
			bot.Delete(challengeMsg)
			log.Printf("%v was banned in %v", m.UserJoined, m.Chat)
		}
		delete(passedUsers, m.UserJoined.ID)
	})
}

func passChallenge(c *tb.Callback) {
	if c.Message.ReplyTo.Sender.ID != c.Sender.ID {
		bot.Respond(c, &tb.CallbackResponse{Text: "Эта кнопка не для вас"})
		return
	}
	passedUsers[c.Sender.ID] = struct{}{}
	bot.Edit(c.Message, "Добро пожаловать!")
	log.Printf("%v passed the challenge in %v", c.Sender, c.Message.Chat)
	newChatMember := tb.ChatMember{User: c.Sender, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: true}}
	bot.Promote(c.Message.Chat, &newChatMember)
	bot.Respond(c, &tb.CallbackResponse{Text: "Доступ разрешен!"})
}