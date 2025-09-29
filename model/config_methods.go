package model

func NewConfig(name string, version int) Config {
	return Config{
		Name:       name,
		Version:    version,
		Parameters: make([]ConfigParameter, 0),
	}
}

func (c *Config) AddParameter(key, value string) {
	param := NewConfigParameter(key, value)
	c.Parameters = append(c.Parameters, param)
}

func (c Config) GetParameter(key string) (string, bool) {
	for _, param := range c.Parameters {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}

// RemoveParameter removes a parameter by key
func (c *Config) RemoveParameter(key string) bool {
	for i, param := range c.Parameters {
		if param.Key == key {
			c.Parameters = append(c.Parameters[:i], c.Parameters[i+1:]...)
			return true
		}
	}
	return false
}

// GetParameterCount returns the number of parameters
func (c Config) GetParameterCount() int {
	return len(c.Parameters)
}