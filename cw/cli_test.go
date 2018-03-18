package main

import (
    "testing"

    cli "github.com/rendon/testcli"
)

const (
    targetBinary    = "./cw"
    helpMessage     = "cw -- Simple command line tool for chatwork API"
    versionRegex    = `chatwork-cli/cw.*ver\.`
)

func TestSmokeRun(t *testing.T) {
    cli.Run(targetBinary)
    if !cli.Success() {
        t.Log(cli.Error())
        t.Fatalf("Command with no argument will be successful, but failed")
    }
}

func TestShowHelp(t *testing.T) {
    cli.Run(targetBinary, "-h")
    if !cli.StdoutContains(helpMessage) {
        t.Log(cli.Stdout())
        t.Fatalf("Command with '-h' will show help, but not supplied")
    }
}

func TestShowVersion(t *testing.T) {
    cli.Run(targetBinary, "-version")
    if !cli.StdoutMatches(versionRegex) {
        t.Log(cli.Stdout())
        t.Fatalf("Command with '-version' will show version, but not supplied")
    }
}

func TestArgumentError(t *testing.T) {
    cli.Run(targetBinary, "-invalidflag")
    if !cli.Failure() {
        t.Fatalf("Command with invalid flag will fail, but succeeded")
    }
}

func TestRunWithEnvironment(t *testing.T) {
    cmd := cli.Command(targetBinary, "get", "me")
    cmd.SetEnv([]string{"CW_API_TOKEN=token"})
    cmd.Run()

    if !cmd.Success() {
        t.Log(cmd.Error())
        t.Log(cmd.Stdout())
        t.Log(cmd.Stderr())
        t.Fatalf("Command using CW_API_TOKEN environment will send request, but failed")
    }
}


