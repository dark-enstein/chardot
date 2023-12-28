package streams

import (
	"github.com/dark-enstein/chardot/cfg"
	"gopkg.in/yaml.v3"
)

func YamlDecode(data []byte) (*cfg.Config, error) {
	var cfg cfg.Config
	err := yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
