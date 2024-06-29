package configs

import (
	"encoding/json"
	"os"
)

func getConfig[T any](path string) (*T, error) {
	plan, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(T)
	err = json.Unmarshal(plan, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
