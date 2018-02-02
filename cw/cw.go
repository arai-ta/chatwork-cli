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

    req.Header.Add("X-ChatWork-Token", "hoge")

    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }

    printHeader(res)
    printBody(res)
}

func printHeader(res *http.Response) {
    fmt.Println("STATUS: ", res.Status)
    fmt.Println("")
    for name, values := range res.Header {
        for _, v := range values {
            fmt.Println(name, ": ", v)
        }
    }
}

func printBody(res *http.Response) {
    io.Copy(os.Stdout, res.Body)
    res.Body.Close()
}
