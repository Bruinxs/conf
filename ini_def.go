package gconf

var (
	defaultIniConfig *IniConfig = NewIniConfig()
)

func Parse(filename string) error {
	_, err := defaultIniConfig.Parse(filename)
	return err
}

func String(key string) string {
	return defaultIniConfig.String(key)
}
