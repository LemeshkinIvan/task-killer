package cfg

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigDTO struct {
	Blacklist   []string `json:"blacklist"`
	TimeRequest string   `json:"time_cfg_request"`
	TimeSleep   string   `json:"time_sleep"`
	LogPath     string   `json:"log_path"`
}

func GetConfig(path string) (*ConfigDTO, error) {
	if path == "" {
		return nil, fmt.Errorf("cfg path is empty")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ConfigDTO
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
