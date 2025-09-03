package controller

import (
	"fmt"
	"net/mail"

	//"fmt"
	//"log"

	"something/models"
	"strconv"
	"strings"
)

type TgMessage struct {
}

// type Comand struct {
// 	ComandName  string `json:"ComandName"`
// 	Description string `json:"Description"`
// }

// type CurrentProcess struct {
// 	Name string
// 	Step int
// }

// type Comands []Comand

// type UserStatus struct {
// 	id         int
// 	statusName string
// }

// type Region struct {
// 	id        int
// 	regioName string
// }

// type TaxSystem struct {
// 	id          int
// 	name        string
// 	description string
// }

// type OKVED struct {
// 	part        string
// 	num         string
// 	description string
// }

// func (o *OKVED) ShowOKVED() string {
// 	str := fmt.Sprintf("%s %s", o.num, o.description)
// 	return str
// }

/*
func initcomands(p string) (com Comands, err error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return com, err
	}
	err = json.Unmarshal(data, &com)
	if err != nil {
		return com, err
	}
	return com, nil
}
*/

func RunRegistration(id int64) {
	var m models.Message
	m.MessageText = "Привет! Я ваш виртуальный бухгалтер. Давайте познакомимся! Хотите пройти регистрацию?"
	//str := "Привет! Я ваш виртуальный бухгалтер. Давайте познакомимся! Хотите пройти регистрацию?"
	var b1, b2 models.TgButton
	b1.Text = "Да"
	b2.Text = "нет"
	m.Buttons = append(m.Buttons, b1, b2)
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	/*
		 	msg := tgbotapi.NewMessage(id, str)
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Да"),
					tgbotapi.NewKeyboardButton("Нет"),
				),
			)

			sessions.sendchan <- msg
	*/
	//m.MessageText = "К сожалению я Вас не понял. Ответьте да или нет"
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		switch strings.ToLower(u.MessageText) {
		case "да":
			break LOOPFOR
		case "нет":
			CancelRegistration(id)
			return
		default:
			m.MessageText = "К сожалению я Вас не понял. Ответьте да или нет"
			serv.sessions.Tgsend.Send.SendMessage(m, id)
			//msg := tgbotapi.NewMessage(id, "К сожалению я Вас не понял. Ответьте да или нет")
			//sessions.sendchan <- msg
		}
	}
	RegistationNewUser(id)

}

func CancelRegistration(id int64) {
	var m models.Message
	m.MessageText = "Вы можете пройти регистрацию позже. До новых встречь!"
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	serv.sessions.Del(id)
	return
}

func RegistationNewUser(id int64) {
	var user models.User
	var cancel bool
	user.UserId = id
	user.Name = regName(id)
	user.Phone, cancel = regPhone(id)
	if cancel {
		CancelRegistration(id)
		return
	}
	user.Email = regMail(id)
	//fmt.Println("User у нас такой ", user)
	us := regUserstatus(id, &serv.sessions.Info)
	user.Type = us.Id
	user.TypeSt = us.StatusName
	ts := regTaxSystem(id, &serv.sessions.Info)
	user.TaxSistem = ts.Id
	user.TaxSistemSt = ts.Name
	//fmt.Println("User у нас такой ", user)
	r := regRegion(id, &serv.sessions.Info)
	user.Region = r.Id
	user.RegionSt = r.RegioName
	//fmt.Println("User у нас такой ", user)
	okv, okvlist := regOKVEDlist(id)
	user.OKVEDS = okvlist
	user.MainOKVED = okv
	//fmt.Println("User у нас такой ", user)
	emp := 0
	if user.Type == 3 || user.Type == 2 {
		emp = regEmployees(id)
	}
	user.Employees = emp
	//fmt.Println("User emploers у нас такой ", user.Employees)
	StartChangeUserInfo(&user, id, true)

}

// Выбор имени пользователя -
func regName(id int64) string {
	var m models.Message
	m.MessageText = "Введите Ваше имя"
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	u := <-serv.sessions.Session[id].Inchan
	return u.MessageText
}

// Выбор статуса пользователя -
func regUserstatus(id int64, listInfo *models.ListInfo) models.UserStatus {
	var i int
	var m models.Message
	var err error
	m.MessageText = "Выберете статус:\n"
	var bs []models.TgButton
	for index, l := range listInfo.UserStatusList {
		m.MessageText += fmt.Sprint(index+1, ". ", l.StatusName, "\n")
		var b models.TgButton
		b.Text = strconv.Itoa(index + 1)
		bs = append(bs, b)
	}
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)

	for u := range serv.sessions.Session[id].Inchan {
		i, err = strconv.Atoi(u.MessageText)
		if err != nil || i <= 0 || i > len(listInfo.UserStatusList) {
			m.MessageText = "К сожалению я Вас не понял. Введите цифру из списка"
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		} else {
			break
		}
	}
	return listInfo.UserStatusList[i-1]
}

// Выбор системы налогооблажения пользователя -
func regTaxSystem(id int64, listInfo *models.ListInfo) models.TaxSystem {
	var m models.Message
	var i int
	var err error
	m.MessageText = "Выберете систему налогооблажения:\n"
	var bs []models.TgButton
	for index, l := range listInfo.TaxSystemList {
		m.MessageText += fmt.Sprint(index+1, ". ", l.Name, "\n")
		var b models.TgButton
		b.Text = strconv.Itoa(index + 1)
		bs = append(bs, b)
	}
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	for u := range serv.sessions.Session[id].Inchan {
		i, err = strconv.Atoi(u.MessageText)
		if err != nil || i <= 0 || i > len(listInfo.UserStatusList) {
			m.MessageText = "К сожалению я Вас не понял. Введите цифру из списка"
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		} else {
			break
		}
	}
	return listInfo.TaxSystemList[i-1]
}

// Выбор региона во время процесса регистрации -
func regRegion(id int64, listInfo *models.ListInfo) models.Region {
	var m models.Message
	var j int
	var err error
	var regs []models.Region
	m.MessageText = "Выберете Регион"
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	for u := range serv.sessions.Session[id].Inchan {
		utext := u.MessageText
		for _, s := range listInfo.RegionsList {
			if strings.Contains(strings.ToLower(s.RegioName), strings.ToLower(utext)) {
				regs = append(regs, s)
			}
		}
		if len(regs) == 0 {
			m.MessageText = "Регион не найден, пропроуйте еще раз. Например - 'Москва'"
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		} else if len(regs) > 10 {
			m.MessageText = "Поиск вернул слишком много совпадений. Попробуйте еще раз."
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		} else {
			break
		}
	}
	if len(regs) == 1 {
		m.MessageText = fmt.Sprint("Выбран регион ", regs[0].RegioName)
		serv.sessions.Tgsend.Send.SendMessage(m, id)
		return regs[0]
	}
	m.MessageText = "Выберете регион из списка совпадений:\n"
	var bs []models.TgButton
	for index, l := range regs {
		var b models.TgButton
		m.MessageText += fmt.Sprint(index+1, ". ", l.RegioName, "\n")
		b.Text = strconv.Itoa(index + 1)
		bs = append(bs, b)
	}
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	for u := range serv.sessions.Session[id].Inchan {
		j, err = strconv.Atoi(u.MessageText)
		if err != nil || j <= 0 || j > len(regs) {
			m.MessageText = "К сожалению я Вас не понял. Введите цифру из списка"
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		} else {
			break
		}
	}
	m.MessageText = fmt.Sprint("Выбран регион ", regs[j-1].RegioName)
	m.Buttons = make([]models.TgButton, 0)
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	return regs[j-1]

}

// загрузка или обновление кодов ОКВЭД
/*
func UpdateOKVED(url string) error {
	// Get the data
	var add, del, update int
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Body.Read())
	st, err := io.ReadAll(resp.Body)
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(st)
	str := string(out)
	//fmt.Println(str)
	old, err := serv.sessions.Repository.GetOKVED()
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
				serv.sessions.Repository.UpdateOKVED(o)
				update++
				//RegionsList
			}
		} else {
			//fmt.Println("добавили ", val, str, o.num, o.part, o.description)

			serv.sessions.Repository.AddOKVED(o)
			add++
		}
		delete(old, str)
	}
	//fmt.Println(len(old))
	for key, val := range old {
		del++
		serv.sessions.Repository.DeleteOKVED(key, val)
	}
	//fmt.Println("Загрузка ОКВЭД завершена. Добавлено:", add, "обновлено:", update, "удалено:", del, " записей")
	return nil
}
*/
// main OKVED и список OKVED -
func regOKVEDlist(id int64) (models.OKVED, []models.OKVED) {
	var m models.Message
	var OkvedList []models.OKVED
	var b models.TgButton
	//var bs []tgbotapi.KeyboardButton
	var bs []models.TgButton
	//var br [][]tgbotapi.KeyboardButton
	m.MessageText = `Выберете вид экономической деятельности (ОКВЭД). ОКВЭД должен содержать минимум 4 цифры. Можно выбрать несколько, но один из них должен быть выбран в качестве основного. Если вы знаете нужный код или описание введите их сразу и мы поищем в справочнике. Если нет, мы поможем его найти.
	Итак начнем:`
	//msg := tgbotapi.NewMessage(id, str)
	//sessions.sendchan <- msg
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	o := regOkvedlistProcess(id)
	OkvedList = append(OkvedList, o)
	m.MessageText = `Выбрать еще один ОКВЭД?`
	b.Text = "да"
	//b = tgbotapi.NewKeyboardButton("да")
	bs = append(bs, b)
	b.Text = "нет"
	//b = tgbotapi.NewKeyboardButton("нет")
	bs = append(bs, b)
	m.Buttons = bs
	m.CounInLine = 10
	//msg = tgbotapi.NewMessage(id, str)
	//br = KeyboardLines(bs, 10)
	//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(br...)
	//sessions.sendchan <- msg
	serv.sessions.Tgsend.Send.SendMessage(m, id)
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		switch strings.ToLower(u.MessageText) {
		case "да":
			o := regOkvedlistProcess(id)
			OkvedList = append(OkvedList, o)
			//sessions.sendchan <- msg
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		case "нет":
			break LOOPFOR
		default:
			var m1 models.Message
			m1.MessageText = "К сожалению я Вас не понял. Ответьте да или нет"
			//msg1 := tgbotapi.NewMessage(id, "К сожалению я Вас не понял. Ответьте да или нет")
			//sessions.sendchan <- msg1
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		}
	}
	//fmt.Println("!!!!!!!!!!! Это оквед лист", OkvedList)
	okv := SearchOKVEDFomSlice(OkvedList, id, "Из списка выбранных кодов ОКВЭД выберите главный, нажатием соответствующей кнопки:\n")
	i := 0
	for index, value := range OkvedList {
		if value == okv {
			i = index
			break
		}
	}
	OkvedList = append(OkvedList[:i], OkvedList[i+1:]...)
	return okv, OkvedList
}

// выбор кода ОКВЭД при регистрации пользователя (одного кода) -
func regOkvedlistProcess(id int64) models.OKVED {
	var okv models.OKVED
	parts, err := serv.sessions.Repository.GetOKVEDpart()
	if err != nil {
		ErrSessionOut(id)
		//fmt.Println(err)
	}
	OKVEDs, err := serv.sessions.Repository.GetOKVEDs()
	if err != nil {
		ErrSessionOut(id)
		//fmt.Println(err)
	}
	var part, head string
	var OKVEDSeach []models.OKVED
	//fmt.Println(parts)
	// := 1
	//var part string
	//var orvedlist []string
	var m models.Message
	var bs []models.TgButton
	//var br [][]tgbotapi.KeyboardButton
	//первый этап выбираем раздел ОКВЭД
	m.MessageText = `Введите код или название или выберите нужный раздел используя кнопки:
	`
	for _, val := range parts {
		m.MessageText += strings.Join([]string{"<b>", val.Part, "</b> - ", val.Description, "\n"}, "")
		var b models.TgButton
		b.Text = val.Part
		bs = append(bs, b)
	}

	//br = append(br, bs)
	m.CounInLine = 10
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)

	// br = KeyboardLines(bs, 10)
	// msg := tgbotapi.NewMessage(id, str)
	// msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(br...)
	// sessions.sendchan <- msg

LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		OKVEDSeach = make([]models.OKVED, 0, 50)
		part = ""
		utext := u.MessageText
		for _, s := range parts {
			if s.Part == utext {
				part = s.Part
				break LOOPFOR
			}
		}
		if part == "" {
			for _, s := range OKVEDs {
				//fmt.Println(s.num, s.description, " значения")
				if strings.Contains(strings.ToLower(s.Num), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
					continue
				}
				if strings.Contains(strings.ToLower(s.Description), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
				}
			}
		}
		if len(OKVEDSeach) == 0 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск не выдал ни одного результата. Воспользуйтесь кнопками чтобы выбрать нужный раздел или попробуйте другую строку поиска`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else if len(OKVEDSeach) > 20 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск вернул слишком много результатов. Попробуйте еще раз:`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else {
			okv = SearchOKVEDFomSlice(OKVEDSeach, id, "Выберете порядковый номер нужного ОКВЭД из списка совпадений, нажатием соответствующей кнопки:\n")
			return okv
		}
	}
	//Второй этап выбираем общеий ОКВЭД формата две цифры внутри выбранного раздела
	parts, err = serv.sessions.Repository.GetOKVEDhead1(part)
	if err != nil {
		ErrSessionOut(id)
		//fmt.Println(err)
	}
	bs = nil
	m.MessageText = `Выберете подходящий варинт используя кнопки:
	`
	for _, val := range parts {
		m.MessageText += strings.Join([]string{"<b>", val.Num, "</b> - ", val.Description, "\n"}, "")
		var b models.TgButton
		b.Text = val.Num
		bs = append(bs, b)
	}
	m.Buttons = bs
	m.CounInLine = 5
	serv.sessions.Tgsend.Send.SendMessage(m, id)
LOOPFOR2:
	for u := range serv.sessions.Session[id].Inchan {
		OKVEDSeach = make([]models.OKVED, 0, 50)
		head = ""
		utext := u.MessageText
		for _, s := range parts {
			if s.Num == utext {
				head = s.Num
				break LOOPFOR2
			}
		}
		if head == "" {
			for _, s := range OKVEDs {
				//fmt.Println(s.Num, s.Description, " значения")
				if strings.Contains(strings.ToLower(s.Num), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
					continue
				}
				if strings.Contains(strings.ToLower(s.Description), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
				}
			}
		}
		if len(OKVEDSeach) == 0 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск не выдал ни одного результата. Воспользуйтесь кнопками чтобы выбрать нужный раздел или попробуйте другую строку поиска`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else if len(OKVEDSeach) > 20 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск вернул слишком много результатов. Попробуйте еще раз:`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else {
			okv = SearchOKVEDFomSlice(OKVEDSeach, id, "Выберете порядковый номер нужного ОКВЭД из списка совпадений, нажатием соответствующей кнопки:\n")
			return okv
		}
	}
	//Третий этап выбираем общеий ОКВЭД формата три цифры внутри выбранного раздела
	parts, err = serv.sessions.Repository.GetOKVEDhead2(part, head)
	if err != nil {
		ErrSessionOut(id)
		//fmt.Println(err)
	}
	bs = nil
	m.MessageText = `Выберете подходящий варинт используя кнопки:
	`
	for _, val := range parts {
		m.MessageText += strings.Join([]string{"<b>", val.Num, "</b> - ", val.Description, "\n"}, "")
		var b models.TgButton
		b.Text = val.Num
		bs = append(bs, b)
	}
	m.Buttons = bs
	m.CounInLine = 5
	serv.sessions.Tgsend.Send.SendMessage(m, id)
LOOPFOR3:
	for u := range serv.sessions.Session[id].Inchan {
		OKVEDSeach = make([]models.OKVED, 0, 50)
		head = ""
		utext := u.MessageText
		for _, s := range parts {
			if s.Num == utext {
				head = s.Num
				break LOOPFOR3
			}
		}
		if head == "" {
			for _, s := range OKVEDs {
				if strings.Contains(strings.ToLower(s.Num), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
					continue
				}
				if strings.Contains(strings.ToLower(s.Description), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
				}
			}
		}
		if len(OKVEDSeach) == 0 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск не выдал ни одного результата. Воспользуйтесь кнопками чтобы выбрать нужный раздел или попробуйте другую строку поиска`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else if len(OKVEDSeach) > 20 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск вернул слишком много результатов. Попробуйте еще раз:`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else {
			okv = SearchOKVEDFomSlice(OKVEDSeach, id, "Выберете порядковый номер нужного ОКВЭД из списка совпадений, нажатием соответствующей кнопки:\n")
			//fmt.Println("обнаружен ", okv)
			return okv
		}
	}
	//Четвертый этап этап выбираем ОКВЭД
	//fmt.Println("Это head!!!!", head)
	parts, err = serv.sessions.Repository.GetOKVEDhead3(part, head)
	if err != nil {
		ErrSessionOut(id)
		//fmt.Println(err)
	}
	bs = nil
	m.MessageText = `Выберете подходящий варинт используя кнопки:
			`
	for _, val := range parts {
		m.MessageText += strings.Join([]string{"<b>", val.Num, "</b> - ", val.Description, "\n"}, "")
		var b models.TgButton
		b.Text = val.Num
		bs = append(bs, b)
	}
	m.CounInLine = 5
	m.Buttons = bs
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	for u := range serv.sessions.Session[id].Inchan {
		OKVEDSeach = make([]models.OKVED, 0, 50)
		head = ""
		utext := u.MessageText
		for _, s := range parts {
			if s.Num == utext {
				head = s.Num
				return s
			}
		}
		if head == "" {
			for _, s := range OKVEDs {
				if strings.Contains(strings.ToLower(s.Num), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
					continue
				}
				if strings.Contains(strings.ToLower(s.Description), strings.ToLower(utext)) {
					OKVEDSeach = append(OKVEDSeach, s)
				}
			}
		}
		if len(OKVEDSeach) == 0 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск не выдал ни одного результата. Воспользуйтесь кнопками чтобы выбрать нужный раздел или попробуйте другую строку поиска`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else if len(OKVEDSeach) > 20 {
			var m1 models.Message
			m1.MessageText = `К сожалению поиск вернул слишком много результатов. Попробуйте еще раз:`
			serv.sessions.Tgsend.Send.SendMessage(m1, id)
		} else {
			okv = SearchOKVEDFomSlice(OKVEDSeach, id, "Выберете порядковый номер нужного ОКВЭД из списка совпадений, нажатием соответствующей кнопки:\n")
			//fmt.Println("обнаружен ", okv)
			return okv
		}
	}

	return okv

}

// выбор ОКВЭД из списка поиска подходящих значений -
func SearchOKVEDFomSlice(okved []models.OKVED, id int64, m string) models.OKVED {
	if len(okved) == 1 {
		return okved[0]
	}
	var msg models.Message
	var j int
	var err error
	msg.MessageText = m
	for index, l := range okved {
		msg.MessageText += fmt.Sprint("<b>", index+1, ". </b>", l.Num, " - ", l.Description, "\n")
		var b models.TgButton
		b.Text = strconv.Itoa(index + 1)
		msg.Buttons = append(msg.Buttons, b)
	}
	msg.CounInLine = 10
	serv.sessions.Tgsend.Send.SendMessage(msg, id)
	for u := range serv.sessions.Session[id].Inchan {
		j, err = strconv.Atoi(u.MessageText)
		if err != nil || j <= 0 || j > len(okved) {
			msg.MessageText = "К сожалению я Вас не понял. Введите цифру из списка"
			serv.sessions.Tgsend.Send.SendMessage(msg, id)
		} else {
			//user.Type = listInfo.UserStatusList[i-1].id
			break
		}
	}
	return okved[j-1]

}

// регистрация количества сотрудников организации -
func regEmployees(id int64) int {
	var j int
	var m models.Message
	var err error
	m.MessageText = `Сколько у вас сотрудников?  
Введите число или 0, если работаете один.`
	serv.sessions.Tgsend.Send.SendMessage(m, id)
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		j, err = strconv.Atoi(u.MessageText)
		if err != nil {
			m.MessageText = `Значение должно быть числом.
			Введите число или 0, если работаете один.`
			serv.sessions.Tgsend.Send.SendMessage(m, id)
			continue
		}
		if j < 0 {
			m.MessageText = `Значение должно быть неотрицательным числом.
			Введите число или 0, если работаете один.`
			serv.sessions.Tgsend.Send.SendMessage(m, id)
			continue
		}
		break LOOPFOR
	}
	return j
}

// регистрация Email -
func regMail(id int64) string {
	var email string
	var m models.Message
	m.MessageText = `Для регистрации необходимо указать Email для важных уведомлений:`
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	m.MessageText = `Формат адреса не корректный, попробуйте еще раз`
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		email = u.MessageText
		_, err := mail.ParseAddress(email)
		if err != nil {
			serv.sessions.Tgsend.Send.SendMessage(m, id)
		} else {
			break LOOPFOR
		}
	}
	return email
}

// изменение телефона -
func changePhone(id int64) string {
	var m models.Message
	var b1 models.TgButton
	var phone string
	m.MessageText = `Необходимо указать телефон для важных уведомлений:
			Формат: 79XX...
			Нажмите "Поделиться телефоном" или введите его вручную.
			`
	b1.Text = "Поделиться телефоном"
	b1.TypeButton = "contact"
	m.Buttons = append(m.Buttons, b1)
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	m.MessageText = `Формат телефона не верный:
			Формат: 79XX...
			всего 11 символов
			Нажмите "Поделиться телефоном" или введите его вручную.
			`
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		if u.MessageText == "" {
			phone = u.Phone
			break LOOPFOR
		}

		phone = u.MessageText
		if checkFormatPhone(phone) {
			break LOOPFOR
		}
		serv.sessions.Tgsend.Send.SendMessage(m, id)

	}
	return phone
}

// регистрация телефона
func regPhone(id int64) (string, bool) {
	var m models.Message
	var b1, b2 models.TgButton
	var phone string
	var cancelreg bool
	m.MessageText = `Для регистрации необходимо указать телефон для важных уведомлений:
			Формат: 79XX...
			Нажмите "Поделиться телефоном" или введите его вручную.
			`

	b1.Text = "Поделиться телефоном"
	b1.TypeButton = "contact"
	m.Buttons = append(m.Buttons, b1)
	b2.Text = "Прекратить регистрацию"
	m.Buttons = append(m.Buttons, b2)
	serv.sessions.Tgsend.Send.SendMessage(m, id)
	m.MessageText = `Формат телефона не верный:
	Формат: 79XX...
	всего 11 символов
	Нажмите "Поделиться телефоном" или введите его вручную.
	`
LOOPFOR:
	for u := range serv.sessions.Session[id].Inchan {
		if u.MessageText == "Прекратить регистрацию" {
			cancelreg = true
			break LOOPFOR
		}
		if u.MessageText == "" {
			phone = u.Phone
			break LOOPFOR
		}
		phone = u.MessageText
		if checkFormatPhone(phone) {
			break LOOPFOR
		}
		serv.sessions.Tgsend.Send.SendMessage(m, id)
	}
	return phone, cancelreg
}
