package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Endpoint string `json:"endpoint"`
}

// Dir returns the config directory: next to the executable in production,
// falls back to ./config in development (go run).
func Dir() string {
	exe, err := os.Executable()
	if err != nil {
		return "config"
	}
	dir := filepath.Join(filepath.Dir(exe), "config")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "config"
	}
	return dir
}

func Load(name string) (*Config, error) {
	path := filepath.Join(Dir(), filename(name))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return &cfg, nil
}

func Save(name string, cfg *Config) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, filename(name)), data, 0644)
}

func List() ([]string, error) {
	entries, err := os.ReadDir(Dir())
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

func Delete(name string) error {
	if name == "default" {
		return fmt.Errorf("cannot delete default config")
	}
	return os.Remove(filepath.Join(Dir(), filename(name)))
}

func filename(name string) string {
	if name == "" {
		return "default.json"
	}
	return name + ".json"
}
