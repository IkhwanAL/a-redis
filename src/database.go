package src

type Database struct {
	Data map[string]interface{}
}

func (d *Database) Set(key string, value string) {
	d.Data[key] = value
}

func (d *Database) Unset(key string) {
	delete(d.Data, key)
}

func (d *Database) Get(key string) string {
	value, ok := d.Data[key]

	if !ok {
		return ""
	}

	return value.(string)
}
