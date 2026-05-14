package agent

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds values from config.json.
type Config struct {
	LANIP     string `json:"lan_ip"`
	Port      int    `json:"port"`
	AuthToken string `json:"auth_token"`
}

// LoadConfig reads and parses the JSON config at path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}