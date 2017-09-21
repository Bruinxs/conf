package conf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	//ErrLoadLoop appoint happen error when load remote config data may be in loop
	ErrLoadLoop = errors.New("load may be in loop error")
)

var (
	varRegex = regexp.MustCompile("\\$\\{.+?\\}")
)

//$ command
func assignVariable(cfg Config, raw string) string {
	val := strings.Trim(raw, " ")
	val = varRegex.ReplaceAllStringFunc(raw, func(match string) string {
		keys := strings.Split(strings.Trim(match, "${}"), ",")
		for _, key := range keys {
			if cfg.Exists(key) {
				return cfg.String(key)
			}
		}
		return match
	})
	return val
}

//@ command
func runOrder(cfg Config, raw string) ([]byte, error) {
	raw = strings.Trim(raw, "@ ")
	vals := strings.Split(raw, ":")
	switch vals[0] {
	case "load":
		uri := assignVariable(cfg, vals[1])
		return load(uri)
	}
	return nil, fmt.Errorf("illegal order '%v'", vals[0])
}

//load command
func load(uri string) ([]byte, error) {
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		resp, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(resp.Body)
	}
	return ioutil.ReadFile(uri)
}
