package cfg

type ConfigMutationDTO struct {
	Blacklist   []string `json:"blacklist"`
	TimeRequest string   `json:"time_cfg_request"`
	TimeSleep   string   `json:"time_sleep"`
	LogPath     string   `json:"log_path"`
}
