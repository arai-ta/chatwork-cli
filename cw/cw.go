package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "strings"
)

func getMethodAndPaths() (string, []string) {
    args := os.Args[1:]
    switch len(args) {
    case 0:
        return http.MethodGet, []string{"me"}
    default:
        return args[0], args[1:]
    }
}

func main() {

    meth, paths := getMethodAndPaths()

    path := strings.Join(paths, "/")

    req, err := http.NewRequest(meth, "https://api.chatwork.com/v2/" + path, nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    req.Header.Add("X-ChatWorkToken", getApiToken())

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

func getApiToken() string {
    return os.Getenv("CW_API_TOKEN")
}
