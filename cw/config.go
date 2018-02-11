package main

import (
    //"fmt"
    "github.com/BurntSushi/toml"
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

/*
func main() {
    cfg, err := ReadConfig("")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(cfg.DefaultProfile)
    fmt.Println(cfg.Profiles["arai"].Host)
    fmt.Println(cfg.Profiles["arai"].Token)
}
*/
