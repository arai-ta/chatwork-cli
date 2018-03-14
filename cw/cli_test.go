package main

import (
    "testing"

    cli "github.com/rendon/testcli"
)

func TestSmokeRun(t *testing.T) {
    cli.Run("./cw")
    if !cli.Success() {
        t.Fatalf("Expected to succeed, but failed: %s", cli.Error())
    }

    if !cli.StdoutContains("Simple command line tool") {
        t.Fatalf("Expected to show help, but not matched: %s", cli.Stdout())
    }
}

