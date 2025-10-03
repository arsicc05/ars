package model

func NewConfigParameter(key, value string) ConfigParameter {
	return ConfigParameter{
		Key:   key,
		Value: value,
	}
}

func (cp ConfigParameter) IsEmpty() bool {
	return cp.Key == "" || cp.Value == ""
}

func (cp ConfigParameter) String() string {
	return cp.Key + "=" + cp.Value
}