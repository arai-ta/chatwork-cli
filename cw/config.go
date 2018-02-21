package main

import (
    "os"
    "os/user"
    "path/filepath"
    "github.com/BurntSushi/toml"
)

type ApiConfig struct {
    DefaultProfile  string `toml:"default_profile"`
    Profiles        map[string]ApiConfigProfile
}

type ApiConfigProfile struct {
    Host    string
    Version string
    Token   string
}

func ReadConfig(filename string) (*ApiConfig, error) {
    if filename == "" {
        filename = getDefaultConfigPath()
    }
    _, err := os.Stat(filename)
    if err != nil {
        return nil, err
    }

    var cfg ApiConfig
    _, err = toml.DecodeFile(filename, &cfg)
    if err != nil {
        return nil, err
    }
    return &cfg, nil
}

func getDefaultConfigPath() string {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    return filepath.Join(user.HomeDir, ".chatwork.toml")
}

