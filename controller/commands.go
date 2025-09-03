package controller

import (
	"log/slog"
	"something/logbot"
	"something/models"
	"something/service"
	"strconv"
	"strings"
	//"fmt"
	//"log"
)

var serv srv

type srv struct {
	log      *slog.Logger
	sessions *service.Sessions
}

func init() {
	serv.log = logbot.Log.With("component", "controller")
	serv.sessions = &service.TgSessions
}

// Если сессия пользователя уже открыта то пересылаем обновление во входящий канал этой сессии если нет то запускаем старт работы бота в отдельной горутине
func Multiplexmassage(ch chan models.MessageResponse) {
	for u := range ch {
		serv.sessions.Repository.Rep.AddMessage(&u)
		if !serv.sessions.Check(u.FromId) {
			err := serv.sessions.Add(u.FromId)
			if err != nil {
				ErrSessionAdd(u.FromId)
				//log.Panic(err)
			} else {
				go StartWork(u.FromId)
			}
		} else {
			serv.sessions.UpdateTimeUpdate(u.FromId)
			serv.sessions.Session[u.FromId].Inchan <- u
		}
	}
}

// помощь список команд
func RunHelp(id int64) {
	var m models.Message
	//m.MessageText = "Это help"
	m.MessageText = `/userinfo - Отображает информацию о пользователе, Есть возможность поменять данные.`
	//str := strings.Join([]string{"Приветствую ", sessions.Session[id].User.Name, " я готов к работе и жду команд. Для вывода всех доступных Вам команд - введите /help"}, "")
	//msg := tgbotapi.NewMessage(id, str)
	//sessions.sendchan <- msg
	serv.sessions.Tgsend.Send.SendMessage(m, id)
}

func WaitComand(id int64) {
	var m, m1 models.Message
	m.MessageText = `Ожидаю команду. Команда начинается с "/". Для вывода всех доступных Вам команд - введите /help`
	for u := range serv.sessions.Session[id].Inchan {
		if u.Command != "" {
			switch strings.ToLower(u.Command) {
			case "help":
				RunHelp(id)
				//m2.MessageText = "ожидаю новую команду команду /help - список всех команд"
				//serv.sessions.Tgsend.Send.SendMessage(m2, id)
				continue
			case "userinfo":
				user := serv.sessions.Session[id].User
				StartChangeUserInfo(&user, id, false)
				//m2.MessageText = "ожидаю новую команду команду /help - список всех команд"
				//serv.sessions.Tgsend.Send.SendMessage(m2, id)
				continue
			default:
				m1.MessageText = "Неизвестная команда. Для вывода всех доступных Вам команд - введите /help"
				serv.sessions.Tgsend.Send.SendMessage(m1, id)
				continue
			}

		}
		serv.sessions.Tgsend.Send.SendMessage(m, id)
	}

}

// приветствие
func RunHello(id int64) {
	var m models.Message
	m.MessageText = strings.Join([]string{"Приветствую ", serv.sessions.Session[id].User.Name, " я готов к работе и жду команд. Для вывода всех доступных Вам команд - введите /help"}, "")
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	WaitComand(id)
}

func StartWork(id int64) {
	if serv.sessions.Session[id].OldUser {
		RunHello(id)
	} else {
		RunRegistration(id)
	}
}

// func StartWork(id int64) {
// 	switch serv.sessions.Session[id].Process.Name {
// 	case "hello":
// 		var m models.Message
// 		m.MessageText = "Тут мы запускам функцию Hello"
// 		serv.sessions.Tgsend.Send.SendMessage(m, id)
// 		RunHello(id)
// 	case "registration":
// 		RunRegistration(id)
// 	}
// }

// показывает информацию пользователя и запускает процесс изменеия информации
func StartChangeUserInfo(user *models.User, id int64, newuser bool) {
	var m models.Message
	var err error
	m.InlineKB = true
	m.MessageText = "Проверьте корректность данных:\n"
	m.MessageText += user.ShowUserInfo()
	m.MessageText += "Верно? (да или нет)"
	var b models.TgButton
	b.Text = "Да"
	b.Data = "Да"
	var bs []models.TgButton
	bs = append(bs, b)
	b.Text = "Нет"
	b.Data = "Нет"
	bs = append(bs, b)
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		switch strings.ToLower(u.MessageText) {
		case "да":
			serv.sessions.UpdateSessionUserData(id, user)
			if newuser {
				err = serv.sessions.Repository.SaveUser(user)
				if err != nil {
					ErrSessionOut(id)
				}
				RunHello(id)
			} else {
				err = serv.sessions.Repository.UpdateUser(user, serv.sessions.Info.OKVEDsList)
				if err != nil {
					ErrSessionOut(id)
				}
			}
			break LOOPFOR
		case "нет":
			ChangeUserInfo(user, id, newuser)
			return
		default:
			m.MessageText = "К сожалению я Вас не понял. Ответьте да или нет"
			//sessions.sendchan <- msg
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		}
	}
	//fmt.Println(str)
}

func ChangeUserInfo(user *models.User, id int64, newuser bool) {
	var m models.Message
	var bs []models.TgButton
	var n int
	m.MessageText, n = user.ShowListProperties()
	m.MessageText = "Что необходимо поменять?:\n" + m.MessageText
	m.MessageText += strconv.Itoa(n+1) + " - Ничего менять не нужно"
	for i := 1; i <= n+1; i++ {
		var b models.TgButton
		b.Text = strconv.Itoa(i)
		//b := tgbotapi.NewKeyboardButton(strconv.Itoa(i))
		bs = append(bs, b)
	}
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		switch strings.ToLower(u.MessageText) {
		case "1":
			user.Name = regName(id)
			break LOOPFOR
		case "2":
			user.Phone = changePhone(id)
			break LOOPFOR
		case "3":
			user.Email = regMail(id)
			break LOOPFOR
		case "4":
			us := regUserstatus(id, &serv.sessions.Info)
			user.Type = us.Id
			user.TypeSt = us.StatusName
			if us.Id == 2 || us.Id == 3 {
				emp := regEmployees(id)
				user.Employees = emp
			} else {
				user.Employees = 0
			}
			break LOOPFOR
		case "5":
			r := regRegion(id, &serv.sessions.Info)
			user.Region = r.Id
			user.RegionSt = r.RegioName
			break LOOPFOR
		case "6":
			ts := regTaxSystem(id, &serv.sessions.Info)
			user.TaxSistem = ts.Id
			user.TaxSistemSt = ts.Name
			break LOOPFOR
		case "7":
			okv, okvlist := regOKVEDlist(id)
			user.OKVEDS = okvlist
			user.MainOKVED = okv
			break LOOPFOR
		case "8":
			if n == 8 {
				emp := regEmployees(id)
				user.Employees = emp
				break LOOPFOR
			} else {
				break LOOPFOR
			}
		case "9":
			if n == 8 {
				break LOOPFOR
			}
			fallthrough
		default:
			m.MessageText = "К сожалению я Вас не понял. введите номер или нажмите соответствующую кнопку"
			//sessions.sendchan <- msg
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		}
	}
	StartChangeUserInfo(user, id, newuser)
}

func ErrSessionOut(id int64) {
	var m models.Message
	m.MessageText = `Система временно не работает. Проблемы на стороне сервера или идут плановые технические работы. Попобуйте позже.`
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	serv.sessions.Del(id)
}

func ErrSessionAdd(id int64) {
	var m models.Message
	m.MessageText = `Система временно не работает. Проблемы на стороне сервера или идут плановые технические работы. Попобуйте позже.`
	serv.sessions.Tgsend.Send.SendMessage(m, id)
}
