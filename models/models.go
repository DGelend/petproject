package models

import (
	"fmt"
	"time"
)

type MainRepository interface {
	UpdateUser(user *User, ol map[string]OKDEDMap) error
	LoadUser(user *User, ol map[string]OKDEDMap) (bool, error)
	SaveUser(user *User) error
	GetInfoList() (ListInfo, error)
	// AddUser(user *User, ol map[string]OKDEDMap) error
	UpdateOKVED(o OKVED) error
	AddOKVED(o OKVED) error
	DeleteOKVED(str string, v string) error
	GetOKVEDpart() ([]OKVED, error)
	GetOKVEDs() ([]OKVED, error)
	GetOKVEDhead1(p string) ([]OKVED, error)
	GetOKVEDhead2(p, n string) ([]OKVED, error)
	GetOKVEDhead3(p, n string) ([]OKVED, error)
	GetOKVED() (map[string]string, error)
	AddMessage(u *MessageResponse) error
}

type ListInfo struct {
	UserStatusList []UserStatus
	TaxSystemList  []TaxSystem
	RegionsList    []Region
	OKVEDsList     map[string]OKDEDMap
}

type UserStatus struct {
	Id         int
	StatusName string
}

type Region struct {
	Id        int
	RegioName string
}

type TaxSystem struct {
	Id          int
	Name        string
	Description string
}

type User struct {
	UserId       int64
	Name         string
	Type         int
	TypeSt       string
	Region       int
	RegionSt     string
	TaxSistem    int
	TaxSistemSt  string
	Phone        string
	Email        string
	MainOKVED    OKVED
	OKVEDS       []OKVED
	Employees    int
	Registration time.Time
}

type Users []User

type OKVED struct {
	Part        string
	Num         string
	Description string
}

type OKDEDMap struct {
	Part        string
	Description string
}

type MessageResponse struct {
	MessageText string
	Phone       string
	FromId      int64
	Command     string
	Data        string
}

func (u User) ShowUserInfo() string {
	str := fmt.Sprintf("- Имя: %s \n- Телефон: %s \n- Email: %s \n- Статус: %s\n- Регион: %s \n- Налоги: %s \n%s", u.Name, u.Phone, u.Email, u.TypeSt, u.RegionSt, u.TaxSistemSt, u.ShowOKVEDList())
	if u.Type == 3 || u.Type == 2 {
		str = str + fmt.Sprintf("- Сотрудники: %d\n", u.Employees)
	}
	return str
}

func (u User) ShowListProperties() (string, int) {
	str := fmt.Sprintf("%d - Имя\n%d - Телефон\n%d - Email\n%d - Статус\n%d - Регион\n%d - Налоги\n%d - ОКВЭД\n", 1, 2, 3, 4, 5, 6, 7)
	i := 7
	if u.Type == 3 || u.Type == 2 {
		str = str + fmt.Sprintf("%d- Сотрудники\n", 8)
		i = 8
	}
	return str, i
}

func (u User) ShowOKVEDList() string {
	str := "- Основной ОКВЭД: " + u.MainOKVED.ShowOKVED() + "\n"
	if len(u.OKVEDS) > 0 {
		str += "- Дополнительные ОКВЭД:\n"
		for _, o := range u.OKVEDS {
			str += " ○ " + o.ShowOKVED() + "\n"
		}
	}
	return str
}

func (o *OKVED) ShowOKVED() string {
	str := fmt.Sprintf("%s %s", o.Num, o.Description)
	return str
}

// текст кнопки и его тип
type TgButton struct {
	Text       string
	Data       string
	TypeButton string
}

// сообщение отправляемое в телеграм (текст сообщения срез кнопок и колччество рядов кнопок)
type Message struct {
	MessageText string
	Buttons     []TgButton
	CounInLine  int
	InlineKB    bool
}

// type MessageResponse struct {
// 	MessageText string
// 	Phone       string
// 	FromId      int64
// 	Command      string
// }

type BotSender interface {
	SendMessage(m Message, id int64) error
}
