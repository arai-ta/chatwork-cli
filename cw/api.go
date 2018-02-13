package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "os"
    "strings"
    "regexp"
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

func NewCwApi() *CwApi {
    api := CwApi{}
    api.Host    = DEFAULT_HOST
    api.Version = DEFAULT_VERSION
    api.Auth    = &TokenFromEnvAuthorizer{DEFAULT_TOKEN_ENV}
    return &api
}

func NewCwApiFromConfig(cfg *ApiConfig) *CwApi {
    name := cfg.DefaultProfile
    prof, ok := cfg.Profiles[name]
    if ok {
        return NewCwApiWithProfile(&prof)
    } else {
        return NewCwApi()
    }
}

func NewCwApiWithProfile(prof *ApiConfigProfile) *CwApi {
    api := NewCwApi()
    if prof.Host != "" {
        api.Host = prof.Host
    }
    if prof.Version != "" {
        api.Version = prof.Version
    }
    switch prof.Auth {
    case "token":
        api.Auth = &TokenFromValueAuthorizer{prof.Token}
    }
    return api
}

// http.Requestをつくる
func (a *CwApi) toRequest() (*http.Request, error) {

    meth := strings.ToUpper(a.Method)
    ok, err := regexp.MatchString(`[A-Z]+`, meth)
    if !ok || err != nil {
        fmt.Println(err)
        return nil, fmt.Errorf("Error: invalid method or error with: %s", a.Method)
    }

    url := "https://" + a.Host + "/" + a.Version + "/" + strings.Join(a.Paths, "/")
    req, err := http.NewRequest(meth, url, nil)

    // fmt.Printf("param len = %d\n", len(a.Param))

    req.Header.Set("User-Agent", getVersion())

    if a.Param != nil && 0 < len(a.Param) {
        query := a.Param.Encode()
        if meth == "GET" {
            req.URL.RawQuery = query
        } else {
            req.Body = ioutil.NopCloser(strings.NewReader(query))
            req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
            req.Header.Set("Content-Length", fmt.Sprintf("%d", len(query)))
            //fmt.Printf("q = %s, len = %d\n", query, len(query))
        }
    }

    if err != nil {
        return req, err
    }
    a.Auth.Authorize(req)
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

// APIトークンをそのまま設定するAuthorizerの実装
type TokenFromValueAuthorizer struct {
    Token string
}

func (ta *TokenFromValueAuthorizer) Authorize(r *http.Request) {
    r.Header.Add("X-ChatWorkToken", ta.Token)
}

