package configs

type ClientConfig struct {
	ServerAddress string `json:"server_address"`
}

func GetClientConfig(path string) (*ClientConfig, error) {
	return getConfig[ClientConfig](path)
}
