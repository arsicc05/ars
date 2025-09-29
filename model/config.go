package model

type ConfigParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Config struct {
	Name       string            `json:"name"`
	Version    int               `json:"version"`
	Parameters []ConfigParameter `json:"parameters"`
}

type GroupConfig struct {
	Name       string            `json:"name"`
	Parameters []ConfigParameter `json:"parameters"`
	Labels     []Label          `json:"labels"`
}

type ConfigGroup struct {
	Name         string        `json:"name"`
	Version      int           `json:"version"`
	Configs      []GroupConfig `json:"configs"`
}

type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version int) (Config, error)
	GetAll() ([]Config, error)
	Update(config Config) error
	Delete(name string, version int) error
}

type ConfigGroupRepository interface {
	Add(group ConfigGroup) error
	Get(name string, version int) (ConfigGroup, error)
	GetAll() ([]ConfigGroup, error)
	Update(group ConfigGroup) error
	Delete(name string, version int) error
}
