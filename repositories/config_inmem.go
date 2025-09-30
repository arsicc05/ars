package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

// Add implements model.ConfigRepository.
func (c *ConfigInMemRepository) Add(config model.Config) error {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	if _, exists := c.configs[key]; exists {
		return errors.New("config already exists")
	}
	c.configs[key] = config
	return nil
}

// Get implements model.ConfigRepository.
func (c *ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}

func (c *ConfigInMemRepository) GetAll() ([]model.Config, error) {
	result := make([]model.Config, 0, len(c.configs))
	for _, cfg := range c.configs {
		result = append(result, cfg)
	}
	return result, nil
}

func (c *ConfigInMemRepository) Update(config model.Config) error {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	if _, exists := c.configs[key]; !exists {
		return errors.New("config not found")
	}
	c.configs[key] = config
	return nil
}

func (c *ConfigInMemRepository) Delete(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	if _, exists := c.configs[key]; !exists {
		return errors.New("config not found")
	}
	delete(c.configs, key)
	return nil
}
