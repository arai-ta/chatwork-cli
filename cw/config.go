package main

import (
    "fmt"
    "os"
    "bytes"
    "github.com/BurntSushi/toml"
)

const (
    DEFAULT_PROFILE_NAME = "default"
)

type ApiConfig struct {
    DefaultProfile  string `toml:"default_profile"`
    Profiles        map[string]ApiConfigProfile
}

type ApiConfigProfile struct {
    Host    string
    Version string
    Auth    string
    Token   string
}

func ReadConfig(filename string) (*ApiConfig, error) {
    if filename == "" {
        filename = "example.toml"
    }
    var cfg ApiConfig
    _, err := toml.DecodeFile(filename, &cfg)
    if err != nil {
        return nil, err
    }
    return &cfg, nil
}

func ReadDefaultConfigOrCreate() (*ApiConfig, error) {
    defaultPath := getDefaultConfigPath()
    _, err := os.Stat(defaultPath)
    if err == nil {
        // file exists
        cfg, err := ReadConfig(defaultPath)
        if err == nil {
            return cfg, nil
        } else {
            return nil, err
        }
    } else {
        // create it
        cfg := CreateNewConfig()
        buf := new(bytes.Buffer)
        enc := toml.NewEncoder(buf)
        enc.Indent = ""
        //if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
        if err := enc.Encode(cfg); err != nil {
            fmt.Println(err)
            return nil, err
        }
        fmt.Println(buf.String())
        return cfg, nil
    }
}

func CreateNewConfig() *ApiConfig {
    cfg := ApiConfig{}
    cfg.DefaultProfile = DEFAULT_PROFILE_NAME
    cfg.Profiles = make(map[string]ApiConfigProfile)
    cfg.Profiles[DEFAULT_PROFILE_NAME] = ApiConfigProfile{}
    return &cfg
}

func getDefaultConfigPath() string {
    return "~/.chatwork.toml"
}

