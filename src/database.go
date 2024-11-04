package src

import "maps"

type Database struct {
	Data map[string]map[string]interface{}
}

func (d *Database) Set(key string, value string) {
	temp := map[string]interface{}{
		"VALUE": value,
	}

	d.Data[key] = temp
}

func (d *Database) SetSetting(key string, settingKey string, value interface{}) {
	temp := map[string]interface{}{
		settingKey: value,
	}

	maps.Copy(d.Data[key], temp)
}

func (d *Database) Unset(key string) {
	delete(d.Data, key)
}

func (d *Database) Get(key string) string {
	value, ok := d.Data[key]["VALUE"]

	if !ok {
		return ""
	}

	return value.(string)
}
