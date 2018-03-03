package conf

var def = NewIniConfig()

//ParseFile open and read the ini file
func ParseFile(file string) error {
	_, err := def.ParseFile(file)
	return err
}

func Exists(key string) bool {
	return def.Exists(key)
}

func String(key string) string {
	return def.String(key)
}

func Strings(key string) []string {
	return def.Strings(key)
}

func Int(key string) (int, error) {
	return def.Int(key)
}

func Int64(key string) (int64, error) {
	return def.Int64(key)
}

func Bool(key string) (bool, error) {
	return def.Bool(key)
}

func Float(key string) (float64, error) {
	return def.Float(key)
}

func DefaultString(key string, defaultVal string) string {
	return def.DefaultString(key, defaultVal)
}
func DefaultStrings(key string, defaultVal []string) []string {
	return def.DefaultStrings(key, defaultVal)
}

func DefaultInt(key string, defaultVal int) int {
	return def.DefaultInt(key, defaultVal)
}

func DefaultInt64(key string, defaultVal int64) int64 {
	return def.DefaultInt64(key, defaultVal)
}

func DefaultBool(key string, defaultVal bool) bool {
	return def.DefaultBool(key, defaultVal)
}

func DefaultFloat(key string, defaultVal float64) float64 {
	return def.DefaultFloat(key, defaultVal)
}
