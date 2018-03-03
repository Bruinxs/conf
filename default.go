package conf

var ini = NewIniConfig()

//ParseFile open and read the ini file
func ParseFile(file string) error {
	_, err := ini.ParseFile(file)
	return err
}

func Exists(key string) bool {
	return ini.Exists(key)
}

func String(key string) string {
	return ini.String(key)
}

func Strings(key string) []string {
	return ini.Strings(key)
}

func Int(key string) (int, error) {
	return ini.Int(key)
}

func Int64(key string) (int64, error) {
	return ini.Int64(key)
}

func Bool(key string) (bool, error) {
	return ini.Bool(key)
}

func Float(key string) (float64, error) {
	return ini.Float(key)
}

func DefaultString(key string, defaultVal string) string {
	return ini.DefaultString(key, defaultVal)
}
func DefaultStrings(key string, defaultVal []string) []string {
	return ini.DefaultStrings(key, defaultVal)
}

func DefaultInt(key string, defaultVal int) int {
	return ini.DefaultInt(key, defaultVal)
}

func DefaultInt64(key string, defaultVal int64) int64 {
	return ini.DefaultInt64(key, defaultVal)
}

func DefaultBool(key string, defaultVal bool) bool {
	return ini.DefaultBool(key, defaultVal)
}

func DefaultFloat(key string, defaultVal float64) float64 {
	return ini.DefaultFloat(key, defaultVal)
}
