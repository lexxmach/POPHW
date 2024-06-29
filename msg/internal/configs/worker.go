package configs

type WorkerConfig struct {
	Address string `json:"address"`

	DB DBConfig `json:"db"`
}

func GetWorkerConfig(path string) (*WorkerConfig, error) {
	return getConfig[WorkerConfig](path)
}
