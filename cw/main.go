package main

import (
    "fmt"
    "flag"
    "net/http"
    "net/url"
    "io"
    "os"
    "strings"
    "sort"
)

const (
    APP_COMMAND = "cw"
    APP_NAME = "chatwork-cli/cw"
    APP_VERSION = "0.1"
)

var (
    optHelp bool
    optVerbose bool
    optProfile string
    optConfigFile string
    optVersion bool
)

func init() {
    flag.BoolVar(&optHelp, "h", false, "Show help message")
    flag.BoolVar(&optVerbose, "v", false, "Dump http headers")
    flag.StringVar(&optProfile, "p", "", "Specify `profile` name to use")
    flag.StringVar(&optConfigFile, "f", "", "Specify `configfile` to use")
    flag.BoolVar(&optVersion, "version", false, "Show version number")
}

func main() {

    flag.Parse()

    if optVersion {
        fmt.Println(getVersion())
        return
    }

    if optHelp || len(flag.Args()) < 2 {
        fmt.Printf(`%s -- Simple command line tool for chatwork API

Usage: %s [options] <verb> [paths...]

Available options:

`, APP_COMMAND, APP_COMMAND)
        flag.PrintDefaults()
        return
    }

    doRequest()
}

func parseArguments(args []string) (string, []string, url.Values) {
    method  := "GET"
    paths   := []string{"me"}
    params  := url.Values{}

    switch num := len(args); {
    case 3 <= num:
        for _, a := range args[2:] {
            if strings.Contains(a, "=") {
                p := strings.SplitN(a, "=", 2)
                params.Set(p[0], p[1])
            } else {
                paths = append(paths, a)
            }
        }
        fallthrough
    case 2 <= num:
        paths[0] = args[1]
        fallthrough
    case 1 <= num:
        method = args[0]
    }

    return method, paths, params
}

func getVersion() string {
    return fmt.Sprintf("%s ver.%s", APP_NAME, APP_VERSION)
}

func doRequest() {
    meth, paths, param := parseArguments(flag.Args())

    api, err := createApi()
    if err != nil {
        fmt.Println(err)
        return
    }
    api.Method = meth
    api.Paths = paths
    api.Param = param

    req, err := api.toRequest()
    if err != nil {
        fmt.Println(err)
        return
    }

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }

    if (optVerbose) {
        printReqHeader(req)
        printResHeader(res)
    }

    printResBody(res)
}

func createApi() (*CwApi, error) {
    cfg, err := ReadConfig(optConfigFile)
    if cfg != nil {
        // config exists
        api, err := NewCwApiFromConfig(cfg, optProfile)
        if err != nil {
            return nil, err
        }
        return api, nil
    } else {
        if os.IsExist(err) {
            // exists, but can not read
            return nil, err
        }
        // not exists. fallback to default
        return NewCwApi(), nil
    }
}

func warn(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, format, args...)
}

func printReqHeader(req *http.Request) {
    warn("> %s %s\n", req.Method, req.URL)
    printHeader(">", req.Header)
    warn(">\n")
}

func printResHeader(res *http.Response) {
    warn("< %s\n", res.Status)
    printHeader("<", res.Header)
    warn("<\n")
}

func printHeader(prefix string, h http.Header) {
    // get keys
    keys := make([]string, len(h))
    for k := range h {
        keys = append(keys, k)
    }
    // sort by name
    sort.Strings(keys)
    for _, name := range keys {
        for _, v := range h[name] {
            warn("%s %s: %s\n", prefix, name, v)
        }
    }
}

func printResBody(res *http.Response) {
    io.Copy(os.Stdout, res.Body)
    res.Body.Close()
}
