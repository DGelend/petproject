package main

import (
	"something/controller"
	"something/models"
	"something/service"
	"something/tgbot"
)

// константы максимальное количество вокеров и время жизни сессии.

var sessions service.Sessions

// версия для публикации
func main() {

	//создаем структуру которая хранит сессии

	//sessions.Create()
	//запускаем обработчик сообщений tg
	respchan := make(chan models.MessageResponse, 100)
	go tgbot.StartTg(respchan)
	go controller.Multiplexmassage(respchan)
	//UpdateOKVED("https://classifikators.ru/assets/downloads/okved/okved.csv")

	end := make(chan struct{})

	<-end
}
