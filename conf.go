package conf

import (
	"io"
)

//Config define configuration setter,getter
type Config interface {
	Exists(key string) bool
	String(key string) string
	Strings(key string) []string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string
	DefaultStrings(key string, defaultVal []string) []string
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
}

//Parser parse raw data to Config
type Parser interface {
	ParseFile(file string) (Config, error)
	Parse(reader io.Reader) (Config, error)
}
