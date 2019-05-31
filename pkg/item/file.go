package item

import (
	"encoding/json"
	"io/ioutil"
)

// FromConfig parses a config from a given file
func FromConfig(filepath string) (Config, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
