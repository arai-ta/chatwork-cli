package main

import (
    "fmt"
    "flag"
    "net/http"
    "net/url"
    "io"
    "os"
    "strings"
)

var (
    optVerbose bool
)

func initFlags() {
    flag.BoolVar(&optVerbose, "v", false, "Dump http headers")
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

func main() {

    initFlags()
    flag.Parse()

    meth, paths, param := parseArguments(flag.Args())

    cfg, err := ReadConfig("")
    if err != nil {
        fmt.Println(err)
        return
    }

    api := NewCwApiFromConfig(cfg)
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
    for name, values := range h {
        for _, v := range values {
            warn("%s %s: %s\n", prefix, name, v)
        }
    }
}

func printResBody(res *http.Response) {
    io.Copy(os.Stdout, res.Body)
    res.Body.Close()
}
