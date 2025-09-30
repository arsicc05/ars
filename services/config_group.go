package services

import (
	"errors"
	"projekat/model"
)

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

func (s ConfigGroupService) Add(group model.ConfigGroup) error {
	return s.repo.Add(group)
}

func (s ConfigGroupService) Get(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigGroupService) GetAll() ([]model.ConfigGroup, error) {
	return s.repo.GetAll()
}


func (s ConfigGroupService) Delete(name string, version int) error {
	return s.repo.Delete(name, version)
}

func (s ConfigGroupService) CreateGroupWithConfig(groupName string, currentVersion int, config model.GroupConfig) (model.ConfigGroup, error) {
	existingGroup, err := s.repo.Get(groupName, currentVersion)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	
	newVersion := currentVersion + 1
	newGroup := model.NewConfigGroup(groupName, newVersion)
	
	for _, existingConfig := range existingGroup.Configs {
		newGroup.AddConfig(existingConfig)
	}
	
	newGroup.AddConfig(config)
	
	err = s.repo.Add(newGroup)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	
	return newGroup, nil
}

func (s ConfigGroupService) CreateGroupWithoutConfig(groupName string, currentVersion int, configName string) (model.ConfigGroup, error) {
	existingGroup, err := s.repo.Get(groupName, currentVersion)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	
	newVersion := currentVersion + 1
	newGroup := model.NewConfigGroup(groupName, newVersion)
	
	found := false
	for _, existingConfig := range existingGroup.Configs {
		if existingConfig.Name != configName {
			newGroup.AddConfig(existingConfig)
		} else {
			found = true
		}
	}
	
	if !found {
		return model.ConfigGroup{}, errors.New("config not found in group")
	}
	
	err = s.repo.Add(newGroup)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	
	return newGroup, nil
}

func (s ConfigGroupService) GetConfig(groupName string, groupVersion int, configName string) (model.GroupConfig, error) {
	group, err := s.repo.Get(groupName, groupVersion)
	if err != nil {
		return model.GroupConfig{}, err
	}
	config, found := group.GetConfig(configName)
	if !found {
		return model.GroupConfig{}, errors.New("config not found in group")
	}
	return config, nil
}
