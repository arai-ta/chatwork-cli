package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io"
    "io/ioutil"
    "os"
    "strings"
)

const (
    DEFAULT_HOST    = "api.chatwork.com"
    DEFAULT_VERSION = "v2"
    DEFAULT_TOKEN_ENV   = "CW_API_TOKEN"
)

// APIへのリクエストに必要な情報を集めた構造体
type CwApi struct {
    // HTTPメソッド
    Method string

    // APIのホスト
    Host string

    // APIバージョン
    Version string

    // エンドポイントのパスまでの配列
    Paths []string

    // リクエストパラメタ
    Param url.Values

    // リクエストに認証情報をつけるオブジェクト
    Auth CwApiAuthorizer
}

// http.Requestをつくる
func (ca *CwApi) toRequest() (*http.Request, error) {
    url := "https://" + ca.Host + "/" + ca.Version + "/" + strings.Join(ca.Paths, "/")
    req, err := http.NewRequest(ca.Method, url, nil)

    if ca.Param != nil {
        query := ca.Param.Encode()
        if strings.ToUpper(ca.Method) == "GET" {
            req.URL.RawQuery = query
        } else {
            req.Body = ioutil.NopCloser(strings.NewReader(query))
        }
    }

    if err != nil {
        return req, err
    }
    ca.Auth.Authorize(req)
    return req, nil
}


// 何かの方法でリクエストに認証情報をつけるオブジェクトを示すinterface
type CwApiAuthorizer interface {
    Authorize(r *http.Request)
}

// 環境変数からAPIトークンを読み取るAuthorizerの実装
type TokenFromEnvAuthorizer struct {
    EnvName string
}

func (ta *TokenFromEnvAuthorizer) Authorize(r *http.Request) {
    token := os.Getenv(ta.EnvName)
    r.Header.Add("X-ChatWorkToken", token)
}

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

    api := CwApi{}
    api.Host = DEFAULT_HOST
    api.Version = DEFAULT_VERSION
    api.Method = meth
    api.Paths = paths
    api.Param = param

    api.Auth = &TokenFromEnvAuthorizer{DEFAULT_TOKEN_ENV}

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
