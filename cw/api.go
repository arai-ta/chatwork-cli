package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	DefaultHost         = "api.chatwork.com"
	DefaultVersion      = "v2"
	DefaultTokenEnvName = "CW_API_TOKEN"
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

	// multipartで送るか
	UseMultiPart bool
}

func NewCwApi() *CwApi {
	api := CwApi{}
	api.Host = DefaultHost
	api.Version = DefaultVersion
	api.Auth = &TokenFromEnvAuthorizer{DefaultTokenEnvName}
	return &api
}

func NewCwApiByConfig(cfg *ApiConfig, profile string) (*CwApi, error) {
	if cfg == nil {
		return NewCwApi(), nil
	}

	if profile == "" {
		profile = cfg.DefaultProfile
	}

	prof, ok := cfg.Profiles[profile]
	if ok {
		return NewCwApiWithProfile(&prof), nil
	} else {
		return nil, fmt.Errorf("profile not found: %s", profile)
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
	if prof.Token != "" {
		api.Auth = &TokenFromValueAuthorizer{prof.Token}
	}
	return api
}

// http.Requestをつくる
func (a *CwApi) toRequest() (*http.Request, error) {
	meth := strings.ToUpper(a.Method)
	ok, _ := regexp.MatchString(`^[A-Z]+$`, meth)
	if !ok {
		return nil, fmt.Errorf("invalid method: %s", a.Method)
	}

	url := "https://" + a.Host + "/" + a.Version + "/" + strings.Join(a.Paths, "/")
	req, err := http.NewRequest(meth, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", getVersion())

	if a.Param != nil && 0 < len(a.Param) {
		query := a.Param.Encode()
		if meth == "GET" {
			req.URL.RawQuery = query
		} else if a.UseMultiPart {
			// multipart
			body := &bytes.Buffer{}
			mw := multipart.NewWriter(body)
			for key, val := range a.Param {
				v := val[0] // 1つめを常に見る
				// attach file
				if v[0] == '@' {
					// open file
					path := v[1:]
					fh, err := os.Open(string(path))
					if err != nil {
						return nil, err
					}
					defer fh.Close()

					// attach
					fw, err := mw.CreateFormFile(key, filepath.Base(path))
					_, err = io.Copy(fw, fh)
					if err != nil {
						return nil, err
					}
				} else {
					fw, err := mw.CreateFormField(key)
					_, err = fw.Write([]byte(v))
					if err != nil {
						return nil, err
					}
				}
			}
			err = mw.Close()
			if err != nil {
				return nil, err
			}
			req.Body = ioutil.NopCloser(body)
			req.Header.Set("Content-Type", mw.FormDataContentType())
		} else {
			// form-urlencoded
			req.Body = ioutil.NopCloser(strings.NewReader(query))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Content-Length", fmt.Sprintf("%d", len(query)))
		}
	}

	err = a.Auth.Authorize(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// 何かの方法でリクエストに認証情報をつけるオブジェクトを示すinterface
type CwApiAuthorizer interface {
	Authorize(r *http.Request) error
}

// 環境変数からAPIトークンを読み取るAuthorizerの実装
type TokenFromEnvAuthorizer struct {
	EnvName string
}

func (ta *TokenFromEnvAuthorizer) Authorize(r *http.Request) error {
	token, ok := os.LookupEnv(ta.EnvName)
	if !ok {
		return fmt.Errorf("environment variable not set: " + ta.EnvName)
	}
	r.Header.Add("X-ChatWorkToken", token)
	return nil
}

// APIトークンをそのまま設定するAuthorizerの実装
type TokenFromValueAuthorizer struct {
	Token string
}

func (ta *TokenFromValueAuthorizer) Authorize(r *http.Request) error {
	r.Header.Add("X-ChatWorkToken", ta.Token)
	return nil
}
