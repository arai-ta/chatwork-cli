package main

import (
    "fmt"
    // "net/http"
    "io/ioutil"
    "sort"
    "strings"

    //"github.com/k0kubun/pp"
    "gopkg.in/yaml.v2"
)

const (
    RamlFileUrl = "https://raw.githubusercontent.com/chatwork/api/master/RAML/api-ja.raml"
    RamlFileName = "api-ja.raml"
)

func main() {
    bytes, err := ioutil.ReadFile(RamlFileName)
    if err != nil {
        fmt.Println(err)
        return
    }
    raml, err := ParseRaml(bytes)

    for _, e := range raml {
        fmt.Printf("%s\t%s -- %s\n", e.Method, e.Path, e.Description)
    }

    if err != nil {
        fmt.Println("NG" + err.Error())
    } else {
        fmt.Println("OK")
    }
}

type EndPoint struct {
    Method  string
    Path    string
    Description string
}

func ParseRaml(data []byte) ([]EndPoint, error) {
    var raml map[interface{}]interface{}
    err := yaml.Unmarshal(data, &raml)

    ep := new([]EndPoint)

    parse(ep, "", raml)

    sort.Slice(*ep, func(i, j int) bool {
        return (*ep)[i].Path < (*ep)[j].Path
    })

    return *ep, err
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
    return
}

