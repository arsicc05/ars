package services

import (
	"projekat/model"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Add(config model.Config) error {
	return s.repo.Add(config)
}

func (s ConfigService) Get(name string, version int) (model.Config, error) {
	return s.repo.Get(name, version)
}

func (s ConfigService) GetAll() ([]model.Config, error) {
	return s.repo.GetAll()
}



func (s ConfigService) Delete(name string, version int) error {
	return s.repo.Delete(name, version)
}

