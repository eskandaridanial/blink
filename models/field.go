package models

type Field struct {
	Key   string
	Value any
}

func String(key, val string) Field {
	return Field{Key: key, Value: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

func Bool(key string, val bool) Field {
	return Field{Key: key, Value: val}
}

func Map(key string, val map[string]any) Field {
	return Field{Key: key, Value: val}
}
