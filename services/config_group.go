package services

import "projekat/model"

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{repo: repo}
}

func (s ConfigGroupService) Add(group model.ConfigGroup) error {
	return s.repo.Add(group)
}

func (s ConfigGroupService) Get(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigGroupService) Delete(name string, version int) error {
	return s.repo.Delete(name, version)
}

func (s ConfigGroupService) Update(group model.ConfigGroup) error {
	return s.repo.Update(group)
}

func (s ConfigGroupService) AddConfigToGroup(name string, version int, cfg model.GroupConfig) error {
	group, err := s.repo.Get(name, version)
	if err != nil {
		return err
	}
	group.AddConfig(cfg)
	return s.repo.Update(group)
}

func (s ConfigGroupService) RemoveConfigFromGroup(name string, version int, cfgName string) error {
	group, err := s.repo.Get(name, version)
	if err != nil {
		return err
	}
	group.RemoveConfig(cfgName)
	return s.repo.Update(group)
}


