package service

import (
	//"database/sql"

	//"fmt"
	"log/slog"
	"os"
	"something/logbot"
	"something/models"
	"something/repository"
	"something/repository/storagedbpg"
	"something/tgbot"
	"sync"
	"time"
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//_ "github.com/lib/pq"
)

var TgSessions Sessions

type Sessions struct {
	log         *slog.Logger
	mu          sync.Mutex
	Session     map[int64]SessionData
	ServerStart time.Time
	//db          *sql.DB
	//sendchan   chan tgbotapi.MessageConfig
	Info       models.ListInfo
	Repository repository.SystemRep
	Tgsend     tgbot.SystemMsg
}

type SessionData struct {
	User    models.User
	OldUser bool
	//inchan     chan tgbotapi.Update
	Inchan     chan models.MessageResponse
	Timeupdate time.Time
}

// type ListInfo struct {
// 	UserStatusList []UserStatus
// 	TaxSystemList  []TaxSystem
// 	RegionsList    []Region
// 	OKVEDsList     map[string]models.OKDEDMap
// }

func init() {
	//TgSessions.Log = logbot.Log
	TgSessions.log = logbot.Log.With("component", "sessions")
	TgSessions.log.Info("Запускаю сервер...")
	TgSessions.ServerStart = time.Now()
	TgSessions.Session = make(map[int64]SessionData, 100)
	stordb, _ := storagedbpg.NewStorageDB(os.Getenv("DB_URL_PG"))
	TgSessions.Repository.Rep = repository.NewRepository(stordb)
	var send tgbot.MytgSend
	//	send.Sendchan = make(chan tgbotapi.MessageConfig, 100)
	TgSessions.Tgsend.Send = &send
	UpdateDirectorys()

	TgSessions.Info, _ = TgSessions.Repository.GetInfoList()
	TgSessions.log.Info("Сервер запущен", "timeStart", TgSessions.ServerStart)
	go SystemWork()
}

// func (s *Sessions) Create() {
// 	s.ServerStart = time.Now()
// 	s.Session = make(map[int64]SessionData, 100)
// 	stordb, _ := storagedbpg.NewStorageDB(os.Getenv("DB_URL_PG"))
// 	s.Repository.Rep = repository.NewRepository(stordb)
// 	var send tgbot.MytgSend
// 	//	send.Sendchan = make(chan tgbotapi.MessageConfig, 100)
// 	s.Tgsend.Send = &send

// 	s.Info, _ = s.Repository.GetInfoList()
// }

func (s *Sessions) Check(id int64) bool {
	TgSessions.log.Debug("Проверяем открыта ли сессия", "tgid", id)
	s.mu.Lock()
	_, ok := s.Session[id]
	s.mu.Unlock()
	return ok
}

func (s *Sessions) Add(id int64) error {
	TgSessions.log.Debug("Добавляем сессию", "tgid", id)
	t := time.Now()
	var user models.User
	user.UserId = id
	olduser, err := s.Repository.Rep.LoadUser(&user, s.Info.OKVEDsList)
	//olduser, user, err := checkUserId(id)
	if err != nil {
		return err
	}
	ch := make(chan models.MessageResponse, 100)
	var sd SessionData = SessionData{user, olduser, ch, t}
	s.mu.Lock()
	s.Session[id] = sd
	s.mu.Unlock()
	return nil
}

func (s *Sessions) Del(id int64) {
	TgSessions.log.Debug("Удаляем сессию", "tgid", id)
	s.mu.Lock()
	close(s.Session[id].Inchan)
	delete(s.Session, id)
	s.mu.Unlock()
}

func (s *Sessions) UpdateTimeUpdate(id int64) {
	TgSessions.log.Debug("Обновляем время последнего обращения", "tgid", id)
	s.mu.Lock()
	ses := s.Session[id]
	ses.Timeupdate = time.Now()
	s.Session[id] = ses
	s.mu.Unlock()
}

func (s *Sessions) UpdateSessionUserData(id int64, user *models.User) {
	TgSessions.log.Debug("Обновляем данные пользователя", "tgid", id)
	s.mu.Lock()
	ses := s.Session[id]
	ses.User = *user
	s.Session[id] = ses
	s.mu.Unlock()
}
