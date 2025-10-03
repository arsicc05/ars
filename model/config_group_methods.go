package model

func NewConfigGroup(name string, version int) ConfigGroup {
	return ConfigGroup{
		Name:    name,
		Version: version,
		Configs: make([]GroupConfig, 0),
	}
}

func (cg *ConfigGroup) AddConfig(config GroupConfig) {
	cg.Configs = append(cg.Configs, config)
}

func (cg ConfigGroup) GetConfig(name string) (GroupConfig, bool) {
	for _, config := range cg.Configs {
		if config.Name == name {
			return config, true
		}
	}
	return GroupConfig{}, false
}

func (cg *ConfigGroup) RemoveConfig(name string) bool {
	for i, config := range cg.Configs {
		if config.Name == name {
			cg.Configs = append(cg.Configs[:i], cg.Configs[i+1:]...)
			return true
		}
	}
	return false
}

func (cg ConfigGroup) GetConfigCount() int {
	return len(cg.Configs)
}

func NewGroupConfig(name string) GroupConfig {
	return GroupConfig{
		Name:       name,
		Parameters: make([]ConfigParameter, 0),
		Labels:     make([]Label, 0),
	}
}

func (gc *GroupConfig) AddParameter(key, value string) {
	param := NewConfigParameter(key, value)
	gc.Parameters = append(gc.Parameters, param)
}

func (gc *GroupConfig) AddLabel(key, value string) {
	label := NewLabel(key, value)
	gc.Labels = append(gc.Labels, label)
}

func (gc GroupConfig) GetParameter(key string) (string, bool) {
	for _, param := range gc.Parameters {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}

func (gc GroupConfig) GetLabel(key string) (string, bool) {
	for _, label := range gc.Labels {
		if label.Key == key {
			return label.Value, true
		}
	}
	return "", false
}