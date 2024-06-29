package configs

type ServerConfig struct {
	Address string `json:"address"`

	DB DBConfig `json:"db"`
}

func GetServerConfig(path string) (*ServerConfig, error) {
	return getConfig[ServerConfig](path)
}
