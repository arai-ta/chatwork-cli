package main

import (
    "fmt"
    "testing"
)

func TestReadConfig(t *testing.T) {
    cfg, err := ReadConfig("")
    if err != nil {
        t.Errorf("something go wrong: %s", err)
    }
    if cfg.DefaultProfile == "" {
        t.Errorf("DefaultProfile is nil")
    }
}

func _TestCreateConfig(t *testing.T) {
    cfg := CreateNewConfig()
    fmt.Println(cfg)
}

func TestCreate(t *testing.T) {
    cfg, err := ReadDefaultConfigOrCreate()
    fmt.Println(cfg)
    fmt.Println(err)
}

