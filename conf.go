package gconf

//define config data setter,getter
type ConfigRaw interface {
	Set(key string, val interface{}) error
	String(key string) string
	Strings(key string) []string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	OptString(key string, defaultVal string) string
	OptStrings(key string, defaultVal []string) []string
	OptInt(key string, defaultVal int) int
	OptInt64(key string, defaultVal int64) int64
	OptBool(key string, defaultVal bool) bool
	OptFloat(key string, defaultVal float64) float64
	Section(key string) map[string]interface{}
	Serialize() ([]byte, error)
}

//define parse config content
type ConfigParser interface {
	Parse(filename string) (ConfigRaw, error)
	ParseData(data []byte) (ConfigRaw, error)
}
