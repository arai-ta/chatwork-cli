package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io"
    "os"
    "strings"
)

func parseArguments(args []string) (string, []string, url.Values) {
    method  := "GET"
    paths   := []string{"me"}
    params  := url.Values{}

    num := len(args)
    switch {
    case 1 <= num:
        method = args[0]
        fallthrough
    case 2 <= num:
        paths = []string{args[1]}
        fallthrough
    case 3 <= num:
        for _, a := range args[2:] {
            if strings.Contains(a, "=") {
                p := strings.SplitN(a, "=", 2)
                params.Set(p[0], p[1])
            } else {
                paths = append(paths, a)
            }
        }
    }

    return method, paths, params
}

func main() {

    meth, paths, param := parseArguments(os.Args[1:])

    cfg, err := ReadConfig("")
    if err != nil {
        fmt.Println(err)
        return
    }

    api := NewCwApi()
    api.Version = cfg.Profiles["arai"].Version
    api.Method = meth
    api.Paths = paths
    api.Param = param

    req, err := api.toRequest()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Print(req.Method+" ")
    fmt.Println(req.URL)
    printHeader(req.Header)

    client := &http.Client{}

    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(res.Status)
    printHeader(res.Header)
    printBody(res)
}

func printHeader(h http.Header) {
    for name, values := range h {
        for _, v := range values {
            fmt.Println(name+": "+v)
        }
    }
}

func printBody(res *http.Response) {
    io.Copy(os.Stdout, res.Body)
    res.Body.Close()
}
