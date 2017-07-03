package conf

//define config data setter,getter
type Config interface {
	Set(key string, val interface{}) error
	String(key string) string
	StrSlice(key string) []string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	StringDef(key string, defaultVal string) string
	StrSliceDef(key string, defaultVal []string) []string
	IntDef(key string, defaultVal int) int
	Int64Def(key string, defaultVal int64) int64
	BoolDef(key string, defaultVal bool) bool
	FloatDef(key string, defaultVal float64) float64
	Section(key string) map[string]interface{}
	Serialize() ([]byte, error)
}

//define parse config content
type ConfigParser interface {
	Parse(filename string) (Config, error)
	ParseData(data []byte) (Config, error)
}
