package service

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"something/models"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

const (
	maxwokercount         = 4
	inactiveTimeout       = 5 * time.Minute
	timeUpdateDirectoryes = 5 * time.Minute
	temecleansessions     = 1 * time.Minute
)

func SystemWork() {
	TgSessions.log.Info("Загрузка служебной горутины.")
	cleanTicker := time.NewTicker(temecleansessions)
	updateDirectoryesTicker := time.NewTicker(timeUpdateDirectoryes)
	for {
		select {
		case <-cleanTicker.C:
			{
				cleanSessions()
			}
		case <-updateDirectoryesTicker.C:
			{
				UpdateDirectorys()
			}
		}
	}
}

func cleanSessions() {
	var n = len(TgSessions.Session)
	var i int
	if n > 0 {
		str := fmt.Sprint("Запускаю удаление сессий. Время жизни сессии:", inactiveTimeout.Minutes(), "минут. Текущее количество сессий:", n)
		TgSessions.log.Info(str)
		for key, s := range TgSessions.Session {
			if time.Now().Sub(s.Timeupdate) > inactiveTimeout {
				TgSessions.Del(key)
				i++
			}
		}
	}
	str := fmt.Sprint("Удалено:", i, " сессий")
	TgSessions.log.Info(str)
}

// обновление срравочной информации
func UpdateDirectorys() {
	updateOKVED("https://classifikators.ru/assets/downloads/okved/okved.csv")
}

// загрузка или обновление кодов ОКВЭД
func updateOKVED(url string) error {
	TgSessions.log.Info("Запускаю обновление кодов ОКВЭД.")
	// Get the data
	var add, del, update int
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	st, err := io.ReadAll(resp.Body)
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(st)
	str := string(out)
	//fmt.Println(str)
	old, err := TgSessions.Repository.GetOKVED()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("длинна олда равна ", len(old))
	reader := bufio.NewReader(strings.NewReader(str))
	for {
		line, err := reader.ReadString('\n')
		var o models.OKVED
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("ошибка чтения: %v\n", err)
			}
			break
		}
		s := strings.Split(line, ";")
		o.Part = strings.Trim(s[0], `"`)
		o.Description = strings.Trim(strings.TrimSuffix(s[2], "\r\n"), `"`)
		st := strings.Trim(strings.Trim(s[1], `"`), " ")
		o.Num = st
		str := o.Part + st
		val, ok := old[str]

		if ok {
			if val != o.Description {
				TgSessions.Repository.UpdateOKVED(o)
				update++
				//RegionsList
			}
		} else {
			//fmt.Println("добавили ", val, str, o.num, o.part, o.description)

			TgSessions.Repository.AddOKVED(o)
			add++
		}
		delete(old, str)
	}
	//fmt.Println(len(old))
	for key, val := range old {
		del++
		TgSessions.Repository.DeleteOKVED(key, val)
	}
	srv := fmt.Sprint("Загрузка ОКВЭД завершена. Добавлено:", add, "обновлено:", update, "удалено:", del, " записей")
	TgSessions.log.Info(srv)
	return nil
}
