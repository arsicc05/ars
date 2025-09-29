package model

func NewLabel(key, value string) Label {
	return Label{
		Key:   key,
		Value: value,
	}
}

func (l Label) IsEmpty() bool {
	return l.Key == "" || l.Value == ""
}

func (l Label) String() string {
	return l.Key + ":" + l.Value
}