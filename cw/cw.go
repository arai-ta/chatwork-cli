package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
)

func main() {
    fmt.Println("Let's go!")

    client := &http.Client{}

    req, err := http.NewRequest("GET", "https://api.chatwork.com/v2/me", nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    req.Header.Add("X-ChatWorkToken", getApiToken())

    fmt.Print(req.Method+" ")
    fmt.Println(req.URL)
    printHeader(req.Header)

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
