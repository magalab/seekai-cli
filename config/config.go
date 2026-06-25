package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type File struct {
	Default  Profile            `toml:"default"`
	Profiles map[string]Profile `toml:"profiles"`
}

type Profile struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.seekai/config.toml"
	}
	return filepath.Join(home, ".seekai", "config.toml")
}

func LoadDefault() (*File, error) {
	path := DefaultPath()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var file File
	if err := toml.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if file.Profiles == nil {
		file.Profiles = map[string]Profile{}
	}
	return &file, nil
}

func SaveDefault(file File) error {
	data, err := toml.Marshal(file)
	if err != nil {
		return err
	}
	path := DefaultPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
