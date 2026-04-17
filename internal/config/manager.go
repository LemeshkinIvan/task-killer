package cfg

import (
	"encoding/json"
	prov "task-killer/internal/provider"
)

type ConfigManager struct {
	provider prov.Provider
}

func NewConfigManager(provider prov.Provider) *ConfigManager {
	return &ConfigManager{
		provider: provider,
	}
}

func (c *ConfigManager) GetConfig() (*ConfigMutationDTO, error) {
	defer c.provider.Disconnect()

	/*bytesCopied for checking hash maybe in future*/
	file, err := c.provider.Get()
	if err != nil {
		return nil, err
	}

	cfg, err := c.UnmarshDTO(file)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *ConfigManager) UnmarshDTO(file []byte) (*ConfigMutationDTO, error) {
	var dto ConfigMutationDTO
	if err := json.Unmarshal(file, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
}
