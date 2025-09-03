package storagedbpg

import (
	"database/sql"
	"log/slog"
	"something/logbot"
	"something/models"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

type StorageDB struct {
	db  *sql.DB
	log *slog.Logger
}

func NewStorageDB(dsn string) (*StorageDB, error) {
	var d StorageDB
	d.log = logbot.Log.With("component", "NewStorageDB")
	d.log.Info("Подключение к БД...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		d.log.Error("Ошибка при открытии БД", "error", err)
		return &d, err
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		d.log.Error("Ошибка при подключении к БД", "error", err)
		return &d, err
	}
	d.db = db
	d.log.Info("Успешно подключились к БД")
	return &d, nil
}

// Сохраням сообщение в БД
func (s *StorageDB) AddMessage(u *models.MessageResponse) error {
	s.log.Debug("Сохраняем сообщение в БД", "tgid", u.FromId)
	_, err := s.db.Exec(
		"INSERT INTO messages (tguserid, message, comand, time) VALUES ($1, $2, $3, $4)",
		u.FromId,
		u.MessageText,
		u.Command,
		time.Now(),
	)
	if err != nil {
		s.log.Error("Ошибка при записи сообщения в БД", "tgid", u.FromId, "error", err)
	}
	return err
}

func (s *StorageDB) SaveUser(user *models.User) error {
	s.log.Debug("Сохраняем пользователя в БД", "tgid", user.UserId)
	b, err := s.CheckUserDB(user)
	if err != nil {
		return err
	}
	if b {
		s.log.Warn("Пользователь уже существует и не будет сохранен", "tgid", user.UserId)
	} else {
		err = s.AddUserToDB(user)
		if err != nil {
			return err
		}
		err = s.AddUserOtherOKVEDs(user)
		if err != nil {
			return err
		}
	}

	return nil

}

// Добавляем нового пользователя в БД
func (s *StorageDB) AddUserToDB(user *models.User) error {
	s.log.Debug("Добавляем поьзовательские данные в БД", "tgid", user.UserId)
	_, err := s.db.Exec(
		"INSERT INTO users (tguserid, username, type_id, regioin, taxsystem, phone, email, okved, emploers, registation) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		user.UserId,
		user.Name,
		user.Type,
		user.Region,
		user.TaxSistem,
		user.Phone,
		user.Email,
		user.MainOKVED.Num,
		user.Employees,
		time.Now(),
	)
	if err != nil {
		s.log.Error("Ошибка при добавлении поьзовательских данных в БД", "tgid", user.UserId, "error", err)
	}
	return err
}

// проверяем есть ли пользователь в базе
func (s *StorageDB) CheckUserDB(user *models.User) (bool, error) {
	s.log.Debug("Проверяем наличие пользователя в БД..", "tgid", user.UserId)
	var i int64
	var check = true
	r := s.db.QueryRow(
		`SELECT tguserid
		FROM users
		WHERE tguserid=$1`, user.UserId)
	err := r.Scan(&i)
	if err != nil {
		check = false
		if err == sql.ErrNoRows {
			s.log.Debug("Пользователь не найден в БД", "tgid", user.UserId)
			err = nil
			//user.UserId = id
		} else {
			s.log.Error("Ошибка при проверке пользователя в БД", "tgid", user.UserId, "error", err)
		}

	} else {
		s.log.Debug("Пользователь найден в БД", "tgid", user.UserId)
	}
	return check, err
}

func (s *StorageDB) UpdateUser(user *models.User, ol map[string]models.OKDEDMap) error {
	s.log.Debug("Обновляем пользователя в БД", "tgid", user.UserId)
	err := s.UpdateUserFromDB(user)
	if err != nil {
		return err
	}
	err = s.UpdattOtherOKVED(user, ol)
	return err
}

func (s *StorageDB) UpdattOtherOKVED(user *models.User, ol map[string]models.OKDEDMap) error {
	s.log.Debug("Обновляем дополнительные ОКВЭД пользователя", "tgid", user.UserId)
	old, err := s.GetUserOtherOkved(user, ol)
	if err != nil {
		return err
	}
	if len(old) > 0 {
		for _, o := range user.OKVEDS {
			flag := false
			var index int
			for i, okv := range old {
				if okv == o {
					flag = true
					index = i
					break
				}
			}
			if !flag {
				s.AddOtherOkved(o, user.UserId)

			} else {
				old = append(old[:index], old[index+1:]...)
			}
		}
		for _, o := range old {
			s.DellOtherOkved(o, user.UserId)
		}
	} else {
		for _, o := range user.OKVEDS {
			s.AddOtherOkved(o, user.UserId)
		}
	}
	return nil
}

// обновляем пользователя в БД.
func (s *StorageDB) UpdateUserFromDB(user *models.User) error {
	s.log.Debug("Обновляем пользовательские данные в БД", "tgid", user.UserId)
	_, err := s.db.Exec(`UPDATE users SET username = $1, type_id = $2, regioin = $3, taxsystem = $4, phone = $5, email = $6, okved = $7, emploers =$8
	WHERE tguserid = $9`,
		user.Name,
		user.Type,
		user.Region,
		user.TaxSistem,
		user.Phone,
		user.Email,
		user.MainOKVED.Num,
		user.Employees,
		user.UserId,
	)
	if err != nil {
		s.log.Error("Ошибка при обновлении поьзовательских данных в БД", "tgid", user.UserId, "error", err)
	}
	return err

}

// Добавляем дополнительные коды ОКВЭД пользователя в БД
func (s *StorageDB) AddUserOtherOKVEDs(user *models.User) error {
	s.log.Debug("Добавляем дополнительные ОКВЭД пользователся в БД", "tgid", user.UserId)
	for _, o := range user.OKVEDS {
		err := s.AddOtherOkved(o, user.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}

// Добавляем один дополнительный код ОКВЭД
func (s *StorageDB) AddOtherOkved(o models.OKVED, id int64) error {
	s.log.Debug("Добавляем дополнительный ОКВЭД пользователя в БД", "OKVED", o.Num, "tgid", id)
	_, err := s.db.Exec(
		"INSERT INTO users_other_okved (tguserid, num) VALUES ($1, $2)",
		//876,
		id,
		o.Num,
	)
	if err != nil {
		s.log.Error("Ошибка при добавлении дополнительного ОКВЭД пользователя в БД", "OKVED", o.Num, "tgid", id, "error", err)
	}
	return err
}

// удаляем один дополнительный код ОКВЭД
func (s *StorageDB) DellOtherOkved(o models.OKVED, id int64) error {
	s.log.Debug("Удаляем дополнительный ОКВЭД пользователя", "OKVED", o.Num, "tgid", id)
	_, err := s.db.Exec(`DELETE FROM users_other_okved  
	WHERE tguserid = $1 and num = $2`,
		id,
		o.Num,
	)
	if err != nil {
		s.log.Error("Ошибка при удалении дополнительного ОКВЭД пользователя в БД", "OKVED", o.Num, "tgid", id, "error", err)
	}
	return err
}

// загружаем пользователя из БД
func (s *StorageDB) LoadUser(user *models.User, ol map[string]models.OKDEDMap) (bool, error) {
	s.log.Debug("Загружаем пользователя из БД", "tgid", user.UserId)
	var check = true
	//var user User
	var o models.OKVED
	r := s.db.QueryRow(
		`SELECT tguserid, username, type_id, ut.usertype_name , regioin, r.region_name , taxsystem,t.taxsystemname, phone, email, okved, emploers, registation
		FROM users u join taxsystem t on t.tax_id = u.taxsystem 
		join region r on r.region_id = u.regioin
		join usertype ut on ut.usertype_id = u.type_id
		WHERE tguserid=$1`, user.UserId)
	err := r.Scan(&user.UserId, &user.Name, &user.Type, &user.TypeSt, &user.Region, &user.RegionSt, &user.TaxSistem, &user.TaxSistemSt, &user.Phone, &user.Email, &o.Num, &user.Employees, &user.Registration)
	o.Part = ol[o.Num].Part
	o.Description = ol[o.Num].Description
	user.MainOKVED = o

	if err != nil {
		check = false
		if err == sql.ErrNoRows {
			err = nil
			s.log.Warn("Пользователь не найден в БД", "tgid", user.UserId)
			//user.UserId = id
		} else {
			s.log.Error("Ошибка при загрузке пользователя из БД", "tgid", user.UserId, "error", err)
		}
	} else {
		err = s.LoadUserOtherOkved(user, ol)
	}
	return check, err
}

// получает дополнительные ОКВЭД
func (s *StorageDB) GetUserOtherOkved(user *models.User, ol map[string]models.OKDEDMap) ([]models.OKVED, error) {
	s.log.Debug("Получаем дополнительные ОКВЭД пользователя", "tgid", user.UserId)
	var okvedlist []models.OKVED
	rows, err := s.db.Query(
		`SELECT num FROM users_other_okved WHERE tguserid=$1`, user.UserId,
	)
	if err != nil {
		s.log.Error("Ошибка при получении дополнительных ОКВЭД пользователя", "tgid", user.UserId, "error", err)
		return okvedlist, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		err = rows.Scan(&o.Num)
		if err != nil {
			s.log.Error("Ошибка при получении дополнительных ОКВЭД пользователя", "tgid", user.UserId, "error", err)
			return okvedlist, err
		}
		o.Part = ol[o.Num].Part
		o.Description = ol[o.Num].Description
		okvedlist = append(okvedlist, o)
	}

	return okvedlist, nil
}

// загружаем дополнительные ОКВЭД
func (s *StorageDB) LoadUserOtherOkved(user *models.User, ol map[string]models.OKDEDMap) error {
	s.log.Debug("Выгружаем дополнительные ОКВЭД пользователя", "tgid", user.UserId)
	//var okvedlist = []OKVED
	rows, err := s.db.Query(
		`SELECT num FROM users_other_okved WHERE tguserid=$1`, user.UserId,
	)
	if err != nil {
		s.log.Error("Ошибка при выгрузке дополнительных ОКВЭД пользователя", "tgid", user.UserId, "error", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		err = rows.Scan(&o.Num)
		if err != nil {
			s.log.Error("Ошибка при выгрузке дополнительных ОКВЭД пользователя", "tgid", user.UserId, "error", err)
			return err
		}
		o.Part = ol[o.Num].Part
		o.Description = ol[o.Num].Description
		user.OKVEDS = append(user.OKVEDS, o)
	}

	return nil
}

// загружаем основную информацию из БД
func (s *StorageDB) GetInfoList() (models.ListInfo, error) {
	s.log.Info("Загружаем основную информацию из БД ifoList..")
	var list models.ListInfo
	var userStatusList []models.UserStatus
	rows, err := s.db.Query(`SELECT usertype_id,usertype_name FROM usertype`)
	if err != nil {
		s.log.Error("Ошибка при получении типа поьзователя в БД", "error", err)
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.UserStatus
		rows.Scan(&u.Id, &u.StatusName)
		userStatusList = append(userStatusList, u)
	}
	var ts []models.TaxSystem
	rows, err = s.db.Query(`SELECT tax_id,taxsystemname ,description  FROM taxsystem`)
	if err != nil {
		s.log.Error("Ошибка при получении типа налогооблажения в БД", "error", err)
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var t models.TaxSystem
		rows.Scan(&t.Id, &t.Name, &t.Description)
		ts = append(ts, t)
	}
	var regioin []models.Region
	rows, err = s.db.Query(`SELECT region_id ,region_name FROM region`)
	if err != nil {
		s.log.Error("Ошибка при получении списка регионов в БД", "error", err)
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var r models.Region
		rows.Scan(&r.Id, &r.RegioName)
		regioin = append(regioin, r)
	}
	sort.Slice(userStatusList, func(i, j int) bool { return userStatusList[i].Id < userStatusList[j].Id })
	sort.Slice(ts, func(i, j int) bool { return ts[i].Id < ts[j].Id })
	list.UserStatusList = userStatusList
	list.TaxSystemList = ts
	list.RegionsList = regioin
	OKVEDs := make(map[string]models.OKDEDMap)
	rows, err = s.db.Query(`SELECT part, num, description FROM okved_directory where num like '__.__%' ORDER BY part`)
	if err != nil {
		s.log.Error("Ошибка при получении списка ОКВЭД в БД", "error", err)
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKDEDMap
		var s string
		rows.Scan(&o.Part, &s, &o.Description)
		OKVEDs[s] = o
	}
	list.OKVEDsList = OKVEDs
	s.log.Info("IfoList загружен")
	return list, nil
}

// Получить список ОКВЭД из БД
func (s *StorageDB) GetOKVED() (map[string]string, error) {
	s.log.Debug("Получаем список всех ОКВЭД из БД")
	rezult := make(map[string]string)
	rows, err := s.db.Query(
		`SELECT part, num, description FROM okved_directory ORDER BY part, num`,
	)
	if err != nil {
		s.log.Error("Ошибка при получении списка всех ОКВЭД в БД", "error", err)
		return rezult, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		rows.Scan(&o.Part, &o.Num, &o.Description)
		str := o.Part + o.Num
		rezult[str] = o.Description
	}
	return rezult, err
}

// Записать ОКВЭД в БД
func (s *StorageDB) AddOKVED(o models.OKVED) error {
	s.log.Debug("Записываю ОКВЭД в БД", "OKVED", o.Num)
	_, err := s.db.Exec(`INSERT INTO okved_directory (part,num,description)
	values ($1, $2, $3)`,
		o.Part,
		o.Num,
		o.Description,
	)
	if err != nil {
		s.log.Error("Ошибка при записи ОКВЭД в БД", "OKVED", o.Num, "error", err)
	}
	return err
}

// обновить ОКВЭД
func (s *StorageDB) UpdateOKVED(o models.OKVED) error {
	s.log.Debug("Обновляю ОКВЭД в БД", "OKVED", o.Num)
	_, err := s.db.Exec(`UPDATE okved_directory SET description = $1 
	WHERE part = $2 AND num = $3`,
		o.Description,
		o.Part,
		o.Num,
	)
	if err != nil {
		s.log.Error("Ошибка при обновлении ОКВЭД в БД", "OKVED", o.Num, "error", err)
	}
	return err

}

// Удалить ОКВЭД
func (s *StorageDB) DeleteOKVED(str string, v string) error {
	s.log.Debug("Удаляю ОКВЭД из БД", "OKVED", v)
	r := []rune(str)
	part := string(r[:1])
	num := string(r[1:])

	_, err := s.db.Exec(`DELETE FROM okved_directory  
	WHERE part = $1 AND num = $2 AND description = $3`,
		part,
		num,
		v,
	)
	if err != nil {
		s.log.Error("Ошибка при удалении ОКВЭД в БД", "OKVED", v, "error", err)
	}
	return err

}

// Ищем все разделы кодов ОКВЭД
func (s *StorageDB) GetOKVEDpart() ([]models.OKVED, error) {
	s.log.Debug("Получаю разделы ОКВЭД в БД")
	rezult := make([]models.OKVED, 0, 30)
	rows, err := s.db.Query(
		`SELECT part, description FROM okved_directory Where num = '' ORDER BY part`,
	)
	if err != nil {
		s.log.Error("Ошибка при получении разделов ОКВЭД в БД", "error", err)
		return rezult, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		rows.Scan(&o.Part, &o.Description)
		rezult = append(rezult, o)

	}
	return rezult, err
}

// ищем в БД все окведы выбранного раздела у которых формат две цифры дд
func (s *StorageDB) GetOKVEDhead1(p string) ([]models.OKVED, error) {
	s.log.Debug("Получаю список ОКВЭД в БД у которых формат две цифры дд в выбнанном разделе", "part", p)
	rezult := make([]models.OKVED, 0, 30)
	rows, err := s.db.Query(
		`select part, num, description FROM okved_directory where part = $1 AND num like '__' ORDER BY part`,
		p,
	)
	if err != nil {
		s.log.Error("Ошибка при получении списка ОКВЭД формата дд в выбранном разделе", "part", p, "error", err)
		return rezult, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		rows.Scan(&o.Part, &o.Num, &o.Description)
		rezult = append(rezult, o)

	}
	return rezult, nil
}

// ищем в БД все окведы выбранного раздела у которых формат две цифры.число
func (s *StorageDB) GetOKVEDhead2(p, n string) ([]models.OKVED, error) {
	s.log.Debug("Получаю список ОКВЭД в БД у которых формат дд.д в выбнанном разделе", "part", p, "префикс", n)
	rezult := make([]models.OKVED, 0, 30)
	rows, err := s.db.Query(
		`SELECT part, num, description FROM okved_directory where part = $1 AND num like $2 || '._' ORDER BY part`,
		p,
		n,
	)
	if err != nil {
		s.log.Error("Ошибка при получении списка ОКВЭД формата дд.д в выбранном разделе", "part", p, "префикс", n, "error", err)
		return rezult, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		rows.Scan(&o.Part, &o.Num, &o.Description)
		rezult = append(rezult, o)

	}
	return rezult, nil
}

// ищем в БД все окведы выбранного раздела у которых формат dd.dd
func (s *StorageDB) GetOKVEDhead3(p, n string) ([]models.OKVED, error) {
	s.log.Debug("Получаю список ОКВЭД в БД у которых формат дд.дд в выбнанном разделе", "part", p, "префикс", n)
	rezult := make([]models.OKVED, 0, 30)
	rows, err := s.db.Query(
		`select part, num, description FROM okved_directory where part = $1 AND num like $2 || '_' ORDER BY part`,
		p,
		n,
	)
	if err != nil {
		s.log.Error("Ошибка при получении списка ОКВЭД формата дд.дд в выбранном разделе", "part", p, "префикс", n, "error", err)
		return rezult, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		rows.Scan(&o.Part, &o.Num, &o.Description)
		rezult = append(rezult, o)

	}
	return rezult, nil
}

// Получить все ОКВЭД формата дд.дд
func (s *StorageDB) GetOKVEDs() ([]models.OKVED, error) {
	s.log.Debug("Получаю список ОКВЭД в БД у которых формат дд.дд")
	rezult := make([]models.OKVED, 0, 3000)
	rows, err := s.db.Query(
		`SELECT part, num, description FROM okved_directory WHERE num LIKE '__.__%'`,
	)
	if err != nil {
		s.log.Error("Ошибка при получении списка ОКВЭД в БД у которых формат дд.дд", "error", err)
		return rezult, err
	}
	defer rows.Close()
	for rows.Next() {
		var o models.OKVED
		rows.Scan(&o.Part, &o.Num, &o.Description)
		rezult = append(rezult, o)

	}
	return rezult, nil
}
