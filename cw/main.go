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
    "bufio"
)

var (
    optVerbose bool
    optConfigure string
    optProfile string
)

func init() {
    flag.BoolVar(&optVerbose, "v", false, "Dump http headers")
    flag.StringVar(&optConfigure, "configure", "", "Configure authentication")
    flag.StringVar(&optProfile, "profile", "", "Specify profile name to use")
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

    flag.Parse()

    if optConfigure != "" {
        doConfigure(optConfigure)
        return
    }

    doRequest()
}

func doConfigure(authType string) {
    switch authType {
    case "token":
        fmt.Print("Enter your API token: ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        token := scanner.Text()
        fmt.Println(token)
    default:
        fmt.Println("Error: invalid configure arg:" + authType)
        return
    }
}

func doRequest() {
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
