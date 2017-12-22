package telegram

import (
	"BadBot/lib"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func MainTtelegram(logs chan string) {
	bot, err := tgbotapi.NewBotAPI(lib.TelegramBotToken)
	if err != nil {
		log.Panic("ошибка подключения бота", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		m := update.Message.Text
		switch {
		case m == "listen":
		loop:
			for{
				logs <- "ok"
				select {
				case u := <-updates:
					if u.Message.Text == "exit" {
						break loop
					}
				case reply := <-logs:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ParseMode = "markdown"
					bot.Send(msg)
				}

			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}
