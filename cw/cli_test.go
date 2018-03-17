package main

import (
    "testing"

    cli "github.com/rendon/testcli"
)

const (
    helpMessage = "cw -- Simple command line tool for chatwork API"
    versionRegex = `chatwork-cli/cw.*ver\.`
)

func TestSmokeRun(t *testing.T) {
    cli.Run("./cw")
    if !cli.Success() {
        t.Log(cli.Error())
        t.Fatalf("Command with no argument will be successful, but failed")
    }
}

func TestShowHelp(t *testing.T) {
    cli.Run("./cw", "-h")
    if !cli.StdoutContains(helpMessage) {
        t.Log(cli.Stdout())
        t.Fatalf("Command with '-h' will show help, but not supplied")
    }
}

func TestShowVersion(t *testing.T) {
    cli.Run("./cw", "-version")
    if !cli.StdoutMatches(versionRegex) {
        t.Log(cli.Stdout())
        t.Fatalf("Command with '-version' will show version, but not supplied")
    }
}

func TestArgumentError(t *testing.T) {
    cli.Run("./cw", "-invalidflag")
    if !cli.Failure() {
        t.Fatalf("Command with invalid flag will fail, but succeeded")
    }
}

