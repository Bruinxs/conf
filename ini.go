package conf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//DefaultSection define default section name
const DefaultSection = "default"

//IniConfig is ini file config
type IniConfig struct {
	dict map[string]string
}

var (
	_ Config = new(IniConfig)
	_ Parser = new(IniConfig)
)

//NewIniConfig return a empty Config
func NewIniConfig() *IniConfig {
	return &IniConfig{dict: map[string]string{}}
}

//ParseFile open and read the ini file
func (ini *IniConfig) ParseFile(file string) (Config, error) {
	bys, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ini.Parse(bytes.NewBuffer(bys))
}

//Parse ini data from reader
func (ini *IniConfig) Parse(reader io.Reader) (Config, error) {
	err := ini.parse(DefaultSection, 0, reader)
	if err != nil {
		return nil, err
	}
	return ini, nil
}

func (ini *IniConfig) parse(section string, deep int, reader io.Reader) error {
	if deep > 20 {
		return ErrLoadLoop
	}

	scaner := bufio.NewScanner(reader)
	for scaner.Scan() {
		line := strings.Trim(scaner.Text(), " ")
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' || line[0] == ';' {
			continue
		}

		if line[0] == '@' {
			bys, err := runOrder(ini, line)
			if err != nil && err != ErrValueEmpty {
				return err
			} else if err == nil {
				err = ini.parse(section, deep+1, bytes.NewBuffer(bys))
				if err != nil {
					return err
				}
			}
			continue
		}

		if line[0] == '[' {
			if line[len(line)-1] != ']' {
				return fmt.Errorf("unclose section '%v'", line)
			}
			section = strings.Trim(line, "[]")
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("no assign value '%v'", line)
		}
		key := kv[0]
		if section != DefaultSection {
			key = section + "." + key
		}
		ini.dict[key] = assignVariable(ini, kv[1])
	}
	return nil
}

func (ini *IniConfig) Exists(key string) bool {
	_, ok := ini.dict[key]
	return ok || os.Getenv(key) != ""
}

//String get value, key is default section key or such as 'section.key' format key
func (ini *IniConfig) String(key string) string {
	str, ok := ini.dict[key]
	if ok {
		return str
	}
	return os.Getenv(key)
}

func (ini *IniConfig) Strings(key string) []string {
	str := ini.String(key)
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}

func (ini *IniConfig) Int(key string) (int, error) {
	return strconv.Atoi(ini.String(key))
}

func (ini *IniConfig) Int64(key string) (int64, error) {
	return strconv.ParseInt(ini.String(key), 10, 64)
}

func (ini *IniConfig) Bool(key string) (bool, error) {
	return strconv.ParseBool(ini.String(key))
}

func (ini *IniConfig) Float(key string) (float64, error) {
	return strconv.ParseFloat(ini.String(key), 64)
}

func (ini *IniConfig) DefaultString(key string, defaultVal string) string {
	str := ini.String(key)
	if str == "" {
		return defaultVal
	}
	return str
}
func (ini *IniConfig) DefaultStrings(key string, defaultVal []string) []string {
	strs := ini.Strings(key)
	if len(strs) == 0 {
		return defaultVal
	}
	return strs
}

func (ini *IniConfig) DefaultInt(key string, defaultVal int) int {
	if !ini.Exists(key) {
		return defaultVal
	}
	val, err := ini.Int(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (ini *IniConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if !ini.Exists(key) {
		return defaultVal
	}
	val, err := ini.Int64(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (ini *IniConfig) DefaultBool(key string, defaultVal bool) bool {
	if !ini.Exists(key) {
		return defaultVal
	}
	val, err := ini.Bool(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (ini *IniConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if !ini.Exists(key) {
		return defaultVal
	}
	val, err := ini.Float(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (ini *IniConfig) PrintAll() {
	for key, val := range ini.dict {
		fmt.Println(key, "=", val)
	}
}
