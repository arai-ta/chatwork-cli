package main

import (
    "testing"

    cli "github.com/rendon/testcli"
)

func TestSmokeRun(t *testing.T) {
    cli.Run("./cw")
    expectSuccess(t)
    expectShowHelp(t)
}

func TestShowHelp(t *testing.T) {
    cli.Run("./cw", "-h")
    expectSuccess(t)
    expectShowHelp(t)
}

func TestShowVersion(t *testing.T) {
    cli.Run("./cw", "-version")
    expectSuccess(t)
    expectShowVersion(t)
}

func TestArgumentError(t *testing.T) {
    cli.Run("./cw", "-invalidflag")
    expectFailure(t)
}


func expectSuccess(t *testing.T) {
    if !cli.Success() {
        t.Fatalf("Expected to succeed, but failed: %s", cli.Error())
    }
}

func expectFailure(t *testing.T) {
    if !cli.Failure() {
        t.Fatalf("Expected to fail, but succeeded")
    }
}

func expectShowHelp(t *testing.T) {
    if !cli.StdoutContains("Simple command line tool") {
        t.Fatalf("Expected to show help, but not matched: %s", cli.Stdout())
    }
}

func expectShowVersion(t *testing.T) {
    regex := `chatwork-cli/cw.*ver\.`
    if !cli.StdoutMatches(regex) {
        t.Fatalf("Expected: %s, actual: %s", regex, cli.Stdout())
    }
}

