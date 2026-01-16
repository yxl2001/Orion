package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Concurrency int           `json:"concurrency"`
	RateLimit   int           `json:"rate_limit"`
	Timeout     time.Duration `json:"timeout"`
	Retry       int           `json:"retry"`
	TargetsFile string        `json:"targets_file"`
	OutputPath  string        `json:"output_path"`
}

func Default() Config {
	return Config{
		Concurrency: 20,
		RateLimit:   200,
		Timeout:     10 * time.Second,
		Retry:       1,
		TargetsFile: "targets.txt",
		OutputPath:  "output.jsonl",
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
