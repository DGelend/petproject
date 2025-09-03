package tgbot

import (
	"log"
	"os"
	"something/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sendmessanges chan tgbotapi.Chattable

func StartTg(resp chan models.MessageResponse) {
	//sendmessanges = make(chan tgbotapi.MessageConfig, 100)
	sendmessanges = make(chan tgbotapi.Chattable, 100)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	// пишем в лог название бота
	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	//запускаем горутину которая будет отправлять сообщения
	//go func(b *tgbotapi.BotAPI, ch chan tgbotapi.MessageConfig) {
	go func(b *tgbotapi.BotAPI, ch chan tgbotapi.Chattable) {
		for m := range ch {
			//тут возможно притормозить отпавку сообщений при надобности
			b.Send(m)
		}
	}(bot, sendmessanges)

	for update := range updates {
		if update.Message != nil {
			//	workchanin <- update
			var m models.MessageResponse
			m.MessageText = update.Message.Text
			if update.Message.Contact != nil {
				m.Phone = update.Message.Contact.PhoneNumber
			}
			m.FromId = update.Message.From.ID
			m.Command = update.Message.Command()

			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp <- m
			//go Multiplexmassage(m)
			//go AddMessage(u)
			//sessions.repository.Rep.AddMessage(&m)
			//log.Printf("это данные пользователя %s", update.Message.From)

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyMarkup = numericKeyboard

			//bot.Send(msg)
		}
		if update.CallbackQuery != nil {
			u := tgbotapi.NewCallback(update.CallbackQuery.ID, "Обработка...")
			sendmessanges <- u
			var m models.MessageResponse
			m.Data = update.CallbackQuery.Data
			m.MessageText = update.CallbackQuery.Data
			m.FromId = update.CallbackQuery.Message.Chat.ID
			var response models.Message
			response.MessageText = m.Data
			msg := tgbotapi.NewMessage(m.FromId, response.MessageText)
			sendmessanges <- msg
			resp <- m
		}
	}

}
