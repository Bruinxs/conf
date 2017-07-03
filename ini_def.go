package conf

var (
	defIniConfig *IniConfig = NewIniConfig()
)

func Parse(filename string) error {
	_, err := defIniConfig.Parse(filename)
	return err
}

func String(key string) string {
	return defIniConfig.String(key)
}

func StringDef(key, def string) string {
	return defIniConfig.StringDef(key, def)
}

func Int(key string) int {
	return defIniConfig.IntDef(key, 0)
}

func IntDef(key string, def int) int {
	return defIniConfig.IntDef(key, def)
}
