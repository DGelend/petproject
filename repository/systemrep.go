package repository

import "something/models"

type SystemRep struct {
	Rep models.MainRepository
}

func NewRepository(repo models.MainRepository) *SystemRep {
	return &SystemRep{Rep: repo}
}

func (s *SystemRep) UpdateUser(user *models.User, ol map[string]models.OKDEDMap) error {
	return s.Rep.UpdateUser(user, ol)
}

func (s *SystemRep) SaveUser(user *models.User) error {
	return s.Rep.SaveUser(user)
}

func (s *SystemRep) GetInfoList() (models.ListInfo, error) {
	return s.Rep.GetInfoList()
}

func (s *SystemRep) UpdateOKVED(o models.OKVED) error {
	return s.Rep.UpdateOKVED(o)
}

func (s *SystemRep) LoadUser(user *models.User, ol map[string]models.OKDEDMap) (bool, error) {
	return s.Rep.LoadUser(user, ol)
}

func (s *SystemRep) AddOKVED(o models.OKVED) error {
	return s.Rep.AddOKVED(o)
}

func (s *SystemRep) DeleteOKVED(str string, v string) error {
	return s.Rep.DeleteOKVED(str, v)
}

func (s *SystemRep) GetOKVEDpart() ([]models.OKVED, error) {
	return s.Rep.GetOKVEDpart()
}

func (s *SystemRep) GetOKVEDs() ([]models.OKVED, error) {
	return s.Rep.GetOKVEDs()
}

func (s *SystemRep) GetOKVEDhead1(p string) ([]models.OKVED, error) {
	return s.Rep.GetOKVEDhead1(p)
}

func (s *SystemRep) GetOKVEDhead2(p, n string) ([]models.OKVED, error) {
	return s.Rep.GetOKVEDhead2(p, n)
}

func (s *SystemRep) GetOKVEDhead3(p, n string) ([]models.OKVED, error) {
	return s.Rep.GetOKVEDhead3(p, n)
}

func (s *SystemRep) GetOKVED() (map[string]string, error) {
	return s.Rep.GetOKVED()
}

func (s SystemRep) AddMessage(m *models.MessageResponse) error {
	return s.Rep.AddMessage(m)
}
