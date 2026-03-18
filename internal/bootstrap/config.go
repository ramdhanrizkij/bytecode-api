package bootstrap

import "github.com/ramdhanrizki/bytecode-api/configs"

func LoadConfig() (configs.Config, error) {
	return configs.Load()
}
