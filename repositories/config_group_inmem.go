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

func groupKey(name string, version int) string {
	return fmt.Sprintf("%s/%d", name, version)
}

// Add implements model.ConfigGroupRepository.
func (r *ConfigGroupInMemRepository) Add(group model.ConfigGroup) error {
	key := groupKey(group.Name, group.Version)
	if _, exists := r.groups[key]; exists {
		return errors.New("config group already exists")
	}
	r.groups[key] = group
	return nil
}

// Get implements model.ConfigGroupRepository.
func (r *ConfigGroupInMemRepository) Get(name string, version int) (model.ConfigGroup, error) {
	key := groupKey(name, version)
	grp, ok := r.groups[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return grp, nil
}

func (r *ConfigGroupInMemRepository) GetAll() ([]model.ConfigGroup, error) {
	result := make([]model.ConfigGroup, 0, len(r.groups))
	for _, g := range r.groups {
		result = append(result, g)
	}
	return result, nil
}

func (r *ConfigGroupInMemRepository) Update(group model.ConfigGroup) error {
	key := groupKey(group.Name, group.Version)
	if _, exists := r.groups[key]; !exists {
		return errors.New("config group not found")
	}
	r.groups[key] = group
	return nil
}

func (r *ConfigGroupInMemRepository) Delete(name string, version int) error {
	key := groupKey(name, version)
	if _, exists := r.groups[key]; !exists {
		return errors.New("config group not found")
	}
	delete(r.groups, key)
	return nil
}


