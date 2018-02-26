package main

import (
    "testing"
)

const ExampleConfig = "example.toml"

func Test_ReadConfig_canReadExampleFile(t *testing.T) {
    cfg, err := ReadConfig(ExampleConfig)
    if err != nil {
        t.Errorf("Failed to read example: %s", err)
    }

    if cfg.DefaultProfile != "default" {
        t.Errorf("DefaultProfile expected 'default', but: %s", cfg.DefaultProfile)
    }

    prof, ok := cfg.Profiles["default"]
    if !ok {
        t.Errorf("Profile[default] expected to be existed, but nil")
    }

    if prof.Token != "mysecrettoken" {
        t.Errorf("profiles.default.token is not expected value: %s", prof.Token)
    }

    if prof.Version != "" {
        t.Errorf("profiles.default.version is not expected value: %s", prof.Version)
    }

    if prof.Host != "" {
        t.Errorf("profiles.default.host is not expected value: %s", prof.Host)
    }

}
