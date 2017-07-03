package conf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type IniConfig struct {
	dict map[string]interface{}
}

var (
	_ Config       = new(IniConfig)
	_ ConfigParser = new(IniConfig)
)

func NewIniConfig() *IniConfig {
	return &IniConfig{dict: map[string]interface{}{}}
}

func (this *IniConfig) envH(str string) string {
	str = strings.Trim(str, " ")
	if strings.Index(str, "${") > -1 {
		pattern := regexp.MustCompile("\\$\\{.*\\}")
		return pattern.ReplaceAllStringFunc(str, func(match string) string {
			keys := strings.Split(strings.Trim(match, " ${}"), ",")
			for _, k := range keys {
				k = strings.Trim(k, " ")
				val, ok := this.dict[k]
				if ok {
					return fmt.Sprintf("%v", val)
				}
			}
			return ""
		})
	}
	return str
}

func (this *IniConfig) parseCommand(command string) error {
	command = this.envH(command)
	if strings.HasPrefix(command, "l(") && strings.HasSuffix(command, ")") && len(command) > 3 {
		url := command[2 : len(command)-1]
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			_, err = this.ParseData(data)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(url, "file://") {
			return errors.New("illegal url")
		} else {
			_, err := this.Parse(url)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New(fmt.Sprintf("command(%v) is illegal", command))
}

func (this *IniConfig) Parse(filename string) (Config, error) {
	filename = path.Clean(filename)
	file, err := os.Open(filename)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return this.ParseData(data)
}

func (this *IniConfig) ParseData(data []byte) (Config, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	var (
		line, section, key string
		slice              []string
	)
	for scanner.Scan() {
		line = strings.Trim(scanner.Text(), " \t\n")
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "@") {
			this.parseCommand(line[1:])
			continue
		}
		if strings.HasPrefix(line, "[") {
			section = strings.Trim(line, " []")
			continue
		}

		slice = strings.SplitN(line, "=", 2)
		if len(slice) != 2 {
			return nil, errors.New(fmt.Sprintf("parse line(%v) to slice illegal", line))
		}
		key = ""
		if len(section) > 0 {
			key = section + "."
		}
		key += this.envH(slice[0])
		this.dict[key] = this.envH(slice[1])
	}
	return this, nil
}

func (this *IniConfig) Set(key string, val interface{}) error {
	this.dict[key] = val
	return nil
}

func (this *IniConfig) Val(key string) interface{} {
	val, ok := this.dict[key]
	if !ok {
		val, ok = this.dict["loc."+key]
		if !ok {
			val = os.Getenv(key)
		}
	}
	return val
}

func (this *IniConfig) String(key string) string {
	val := this.Val(key)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

func (this *IniConfig) StrSlice(key string) []string {
	str := this.String(key)
	if len(str) == 0 {
		return nil
	}
	return strings.Split(str, ",")
}

func (this *IniConfig) Int(key string) (int, error) {
	val := this.Val(key)
	iv, ok := val.(int)
	if !ok {
		var err error
		iv, err = strconv.Atoi(this.String(key))
		if err != nil {
			return -1, err
		}
	}
	return iv, nil
}

func (this *IniConfig) Int64(key string) (int64, error) {
	val := this.Val(key)
	iv64, ok := val.(int64)
	if !ok {
		var err error
		iv64, err = strconv.ParseInt(this.String(key), 10, 64)
		if err != nil {
			return -1, err
		}
	}
	return iv64, nil
}

func (this *IniConfig) Bool(key string) (bool, error) {
	val := this.Val(key)
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		switch v {
		case "true", "True", "yes", "Yes", "1":
			return true, nil
		case "false", "False", "no", "No", "0":
			return false, nil
		default:
			return false, errors.New(fmt.Sprintf("value(%v), type(%v), key(%v) assert to bool fail", v, v, key))
		}
	case int, int32, int64, float32, float64:
		is := fmt.Sprintf("%v", v)
		if is == "1" {
			return true, nil
		} else if is == "0" {
			return false, nil
		} else {
			return false, errors.New(fmt.Sprintf("value(%v), type(%v), key(%v) assert to bool fail", v, v, key))
		}
	default:
		return false, errors.New(fmt.Sprintf("value(%v), type(%v), key(%v) assert to bool fail", v, v, key))
	}
}

func (this *IniConfig) Float(key string) (float64, error) {
	val := this.Val(key)
	fv, ok := val.(float64)
	if !ok {
		var err error
		fv, err = strconv.ParseFloat(this.String(key), 64)
		if err != nil {
			return 0, err
		}
	}
	return fv, nil
}

func (this *IniConfig) StringDef(key, defaultVal string) string {
	str := this.String(key)
	if str == "" {
		return defaultVal
	}
	return str
}

func (this *IniConfig) StrSliceDef(key string, defaultVal []string) []string {
	slice := this.StrSlice(key)
	if len(slice) == 0 {
		return defaultVal
	}
	return slice
}

func (this *IniConfig) IntDef(key string, defaultVal int) int {
	iv, err := this.Int(key)
	if err != nil {
		return defaultVal
	}
	return iv
}

func (this *IniConfig) Int64Def(key string, defaultVal int64) int64 {
	iv, err := this.Int64(key)
	if err != nil {
		return defaultVal
	}
	return iv
}

func (this *IniConfig) BoolDef(key string, defaultVal bool) bool {
	bv, err := this.Bool(key)
	if err != nil {
		return defaultVal
	}
	return bv
}

func (this *IniConfig) FloatDef(key string, defaultVal float64) float64 {
	fv, err := this.Float(key)
	if err != nil {
		return defaultVal
	}
	return fv
}

func (this *IniConfig) Section(key string) map[string]interface{} {
	panic("illegal methon")
}

func (this *IniConfig) Serialize() ([]byte, error) {
	panic("illegal methon")
}
