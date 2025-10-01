package services

import (
	"errors"
	"projekat/model"
	"strings"
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

// labels string format example: "k1:v1;k2:v2"
func parseLabelsStringToMap(labelsStr string) (map[string]string, error) {
    labels := make(map[string]string)
    if strings.TrimSpace(labelsStr) == "" {
        return labels, nil
    }
    pairs := strings.Split(labelsStr, ";")
    for _, pair := range pairs {
        pair = strings.TrimSpace(pair)
        if pair == "" {
            continue
        }
        kv := strings.SplitN(pair, ":", 2)
        if len(kv) != 2 {
            return nil, errors.New("invalid labels format; expected key:value pairs separated by ';'")
        }
        key := strings.TrimSpace(kv[0])
        value := strings.TrimSpace(kv[1])
        if key == "" || value == "" {
            return nil, errors.New("invalid labels: empty key or value")
        }
        labels[key] = value
    }
    return labels, nil
}

func configMatchesAllLabels(cfg model.GroupConfig, required map[string]string) bool {
    if len(required) == 0 {
        return true
    }
    have := make(map[string]string)
    for _, l := range cfg.Labels {
        have[l.Key] = l.Value
    }
    for k, v := range required {
        hv, ok := have[k]
        if !ok || hv != v {
            return false
        }
    }
    return true
}

func (s ConfigGroupService) FilterConfigsByLabels(groupName string, groupVersion int, labelsStr string) ([]model.GroupConfig, error) {
    group, err := s.repo.Get(groupName, groupVersion)
    if err != nil {
        return nil, err
    }
    labelsMap, err := parseLabelsStringToMap(labelsStr)
    if err != nil {
        return nil, err
    }
    result := make([]model.GroupConfig, 0)
    for _, cfg := range group.Configs {
        if configMatchesAllLabels(cfg, labelsMap) {
            result = append(result, cfg)
        }
    }
    return result, nil
}

func (s ConfigGroupService) CreateGroupWithoutConfigsByLabels(groupName string, currentVersion int, labelsStr string) (model.ConfigGroup, error) {
    existingGroup, err := s.repo.Get(groupName, currentVersion)
    if err != nil {
        return model.ConfigGroup{}, err
    }
    labelsMap, err := parseLabelsStringToMap(labelsStr)
    if err != nil {
        return model.ConfigGroup{}, err
    }

    newVersion := currentVersion + 1
    newGroup := model.NewConfigGroup(groupName, newVersion)

    removedAny := false
    for _, cfg := range existingGroup.Configs {
        if configMatchesAllLabels(cfg, labelsMap) {
            removedAny = true
            continue
        }
        newGroup.AddConfig(cfg)
    }

    if !removedAny {
        return model.ConfigGroup{}, errors.New("no configs matched given labels")
    }

    if err := s.repo.Add(newGroup); err != nil {
        return model.ConfigGroup{}, err
    }

    return newGroup, nil
}
