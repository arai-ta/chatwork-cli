package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "sort"
    "strings"

    "gopkg.in/yaml.v2"
)

const (
    OfficialRamlFileUrl = "https://raw.githubusercontent.com/chatwork/api/master/RAML/api-ja.raml"
)

type EndPoint struct {
    Method  string
    Path    string
    Description string
}

func ShowEndPoints(raml string) error {
    if raml == "" {
        raml = OfficialRamlFileUrl
    }

    bytes, err := GetRaml(raml)
    if err != nil {
        return err
    }

    endpoints, err := ParseRaml(bytes)
    if err != nil {
        return err
    }

    for _, e := range endpoints {
        fmt.Printf("%s\t%s -- %s\n", e.Method, e.Path, e.Description)
    }

    return nil
}

func GetRaml(location string) ([]byte, error) {
    var data []byte
    var err error

    u, err := url.Parse(location)
    if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
        // It's a URL
        resp, err := http.Get(location)
        if err == nil {
            data, err = ioutil.ReadAll(resp.Body)
            defer resp.Body.Close()
        }
    } else {
        // It may be a file
        data, err = ioutil.ReadFile(location)
    }
    if err != nil {
        return nil, err
    }

    return data, nil
}

func ParseRaml(data []byte) ([]EndPoint, error) {
    var raml map[interface{}]interface{}
    var ep []EndPoint

    err := yaml.Unmarshal(data, &raml)
    if err != nil {
        return ep, err
    }

    parse(&ep, "", raml)
    sort.Slice(ep, func(i, j int) bool {
        return (ep)[i].Path < (ep)[j].Path
    })

    return ep, nil
}

func parse(ep *[]EndPoint, current string, node map[interface{}]interface{}) {
    for k, v := range node {
        ks := k.(string)
        switch ks {
        case "GET", "POST", "PUT", "DELETE":
            m := v.(map[interface{}]interface{})
            d := strings.Trim(m["description"].(string), " \n")
            e := EndPoint{ks, current, d}
            // fmt.Println(e)
            *ep = append(*ep, e)
        default:
            if ks[0] == '/' {
                m := v.(map[interface{}]interface{})
                parse(ep, current + ks, m)
            }
        }
    }
}

