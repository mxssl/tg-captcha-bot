package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"runtime"
	"time"
)

// Токен для бота получаем из переменной окружения
var tgToken = os.Getenv("TGTOKEN")

func main() {
	// Настройки обращения к API telegram
	b, err := tb.NewBot(tb.Settings{
		Token:  tgToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	// Если подключиться к API telegram не удалось, тогда крашимся
	if err != nil {
		log.Fatal(err)
		return
	}

	// Описываем кнопку, которая будет отображаться новому посетителю
	inlineBtn := tb.InlineButton{
		Unique: "sad_moon",
		Text:   "Я не спамер!",
	}
	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn},
	}

	// Переменная, в которую будет записываться ID того, кто нажал на кнопку
	var temp int

	// Делаем callback на нажатие кнопки и записываем ID того, кто ее нажал
	b.Handle(&inlineBtn, func(c *tb.Callback) {
		temp = c.Sender.ID
		b.Respond(c, &tb.CallbackResponse{Text: "Доступ разрешен!"})
	})

	/*
	Основная функция с логикой работы бота
	Задача проверить, что тот кто нажал на кнопку == новому посетителю
	*/
	go b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		// Обнуляем значение переменной для ID того, кто нажал на кнопку
		temp = 0

		// Вешаем на каждого нового пользователя рестрикт на отправку сообщений
		newChatMember := tb.ChatMember{User: m.UserJoined,
		RestrictedUntil: tb.Forever(),
		Rights: tb.Rights{CanSendMessages: false},
		}
		b.Restrict(m.Chat, &newChatMember)
		log.Printf("Присоединился пользователь: %v", m.UserJoined)

		// Запоминаем инфо о новом посетителе
		username := m.UserJoined.Username
		firstname := m.UserJoined.FirstName

		// Переменная, в которую будет записываться собранное сообщение с приветствием
		var msg string

		// Собираем сообщение в зависимости есть ли у нового пользователя username
		warning := "Это защита от спама. У вас есть 30 секунд нажать на кнопку. Иначе вы будете забанены!"
		if username != "" {
			msg = fmt.Sprintf("@%v\n"+warning, username)
		} else {
			msg = fmt.Sprintf("%v\n"+warning, firstname)
		}

		// Отправляем собранное сообщение в чат
		botMsg, _ := b.Send(m.Chat, msg, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})

		// Переменная назначается в зависимости есть ли у нового посетителя username
		var u string
		if username != "" {
			u = username
		} else {
			u = firstname
		}

		/*
		30 секунд ждем нажатия нового посетителя на кнопку.
		Если нажатия не произошло вешаем бан
		*/
		var idx int
		for start := time.Now(); ; {
			if idx%30 == 0 {
				if time.Since(start) > 30*time.Second {
					break
				}
			}

			idx++

			/*
			Если тот, кто нажал на кнопку == новый посетитель,
			тогда снимаем с него рестрикт и выходим из цикла
			 */
			if temp == m.UserJoined.ID {
				msgCheckPassed := fmt.Sprintf("@%v\nДобро пожаловать!", u)
				b.Edit(botMsg, msgCheckPassed, tb.ParseMode(tb.ModeMarkdown), tb.NoPreview)
				log.Printf("Пользователь прошел проверку: %v", m.UserJoined)
				newChatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: true}}
				b.Promote(m.Chat, &newChatMember)
				return
			}
		}

		// Если в течении 30ти секунд нажатия не произошло, баним нового посетителя
		chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever()}
		b.Ban(m.Chat, &chatMember)
		banMsg := fmt.Sprintf("@%v не прошел проверку", u)
		b.Edit(botMsg, banMsg)
		log.Printf("Пользователь забанен: %v", m.UserJoined)
	})

	// Команда для проверки работоспособности бота
	b.Handle("/healthz", func(m *tb.Message) {
		runtimeVer := runtime.Version()
		verMsg := fmt.Sprintf(`Я здоров!
Версия Go: %v`, runtimeVer)
		b.Send(m.Chat, verMsg)
	})

	// Запускаем бота
	log.Print("Бот запущен!")
	b.Start()
}

