package main

import (
	"log"
	"net/url"
	"testing"
)

func Test_Api(t *testing.T) {

	//	api := NewCwApi()
	//	api.Method = "POST"
	//	api.Version = "v2"
	//	api.Host = "api.chatwork.com"
	//	api.Paths = []string{"aaa", "bbbb"}

	param, _ := url.ParseQuery("foo=bar&hoge=fufa&file=@path")

	// log.Printf("%#v", api)
	log.Printf("%#v", param)

	var f string

	for key, val := range param {
		if val[0][0] == '@' {
			param.Del(key)
			f = val[0][1:]
		}
		log.Printf("%s -> %s\n", key, val)
		log.Printf("%T -> %T\n", key, val)
	}

	//	api.Param = param
	//	req, _ := api.toRequest()
	//	log.Printf("%#v", req.URL.String())
	log.Printf("%#v", f)

	t.Errorf("api test error end")

}
