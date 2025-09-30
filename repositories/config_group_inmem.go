package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigGroupInMemRepository struct {
	groups map[string]model.ConfigGroup
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemRepository{
		groups: make(map[string]model.ConfigGroup),
	}
}

func (r *ConfigGroupInMemRepository) Add(group model.ConfigGroup) error {
	key := fmt.Sprintf("%s/%d", group.Name, group.Version)
	if _, exists := r.groups[key]; exists {
		return errors.New("config group already exists")
	}
	r.groups[key] = group
	return nil
}

func (r *ConfigGroupInMemRepository) Get(name string, version int) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	group, ok := r.groups[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return group, nil
}

func (r *ConfigGroupInMemRepository) GetAll() ([]model.ConfigGroup, error) {
	groups := make([]model.ConfigGroup, 0, len(r.groups))
	for _, group := range r.groups {
		groups = append(groups, group)
	}
	return groups, nil
}


func (r *ConfigGroupInMemRepository) Delete(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	if _, exists := r.groups[key]; !exists {
		return errors.New("config group not found")
	}
	delete(r.groups, key)
	return nil
}
