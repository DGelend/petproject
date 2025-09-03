package tgbot

import (
	"something/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MytgSend struct {
}

type SystemMsg struct {
	Send models.BotSender
}

// func NewSender(send models.BotSender) *SystemMsg {
// 	return &SystemMsg{Send: send}
// }

// создаем InLine клавиатуру
func CreateInlineKeyBoad(m models.Message) tgbotapi.InlineKeyboardMarkup {
	var bs []tgbotapi.InlineKeyboardButton
	var br [][]tgbotapi.InlineKeyboardButton
	inline := 5
	for _, bt := range m.Buttons {
		switch bt.TypeButton {
		case "URL":
			{
				b := tgbotapi.NewInlineKeyboardButtonURL(bt.Text, bt.Data)
				bs = append(bs, b)

			}
		default:
			{
				b := tgbotapi.NewInlineKeyboardButtonData(bt.Text, bt.Data)
				bs = append(bs, b)

			}
		}
	}
	if m.CounInLine > 0 {
		inline = m.CounInLine
	}
	//var keyboard tgbotapi.ReplyKeyboardMarkup
	var keyboard tgbotapi.InlineKeyboardMarkup
	if len(bs) > 0 {
		br = KeyboardLinesInline(bs, inline)
		//kb := tgbotapi.NewInlineKeyboardMarkup(br...)
		keyboard = tgbotapi.NewInlineKeyboardMarkup(br...)
		//keyboard = tgbotapi.NewReplyKeyboard(br...)
		//keyboard.OneTimeKeyboard = true
		//keyboard.ResizeKeyboard = true
	}
	return keyboard
}

// создаем обычную клавиатуру
func CreateKeyBoad(m models.Message) tgbotapi.ReplyKeyboardMarkup {
	var bs []tgbotapi.KeyboardButton
	var br [][]tgbotapi.KeyboardButton
	inline := 5
	for _, bt := range m.Buttons {
		switch bt.TypeButton {
		case "contact":
			{
				b := tgbotapi.NewKeyboardButtonContact(bt.Text)
				bs = append(bs, b)

			}
		default:
			{
				b := tgbotapi.NewKeyboardButton(bt.Text)
				bs = append(bs, b)

			}
		}
	}
	if m.CounInLine > 0 {
		inline = m.CounInLine
	}
	var keyboard tgbotapi.ReplyKeyboardMarkup
	if len(bs) > 0 {
		br = KeyboardLines(bs, inline)
		keyboard = tgbotapi.NewReplyKeyboard(br...)
		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true
	}
	return keyboard
}

// Создаем сообщение для бота
func (s *MytgSend) SendMessage(m models.Message, id int64) {
	msg := tgbotapi.NewMessage(id, m.MessageText)
	if len(m.Buttons) == 0 {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	} else {
		if m.InlineKB {
			keyboard := CreateInlineKeyBoad(m)
			msg.ReplyMarkup = keyboard
		} else {
			keyboard := CreateKeyBoad(m)
			msg.ReplyMarkup = keyboard
		}
		//msg.ReplyMarkup = keyboard
	}
	msg.ParseMode = "HTML"
	sendmessanges <- msg
}

// при необходимости создаем несколько рядов обчных кнопок
func KeyboardLines(l []tgbotapi.KeyboardButton, n int) [][]tgbotapi.KeyboardButton {
	var result [][]tgbotapi.KeyboardButton
	len := len(l)
	if len > n {
		for i := 0; i < len; {
			result = append(result, l[i:i+n])
			i = i + n
			if len > i+n {
				continue
			} else {
				result = append(result, l[i:len])
				break
			}
		}
	} else {
		result = append(result, l)
	}
	//fmt.Println(result)
	return result
}

// при необходимости создаем несколько рядов inline кнопок
func KeyboardLinesInline(l []tgbotapi.InlineKeyboardButton, n int) [][]tgbotapi.InlineKeyboardButton {
	var result [][]tgbotapi.InlineKeyboardButton
	len := len(l)
	if len > n {
		for i := 0; i < len; {
			result = append(result, l[i:i+n])
			i = i + n
			if len > i+n {
				continue
			} else {
				result = append(result, l[i:len])
				break
			}
		}
	} else {
		result = append(result, l)
	}
	//fmt.Println(result)
	return result
}
