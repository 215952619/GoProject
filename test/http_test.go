package test

import (
	"GoProject/util"
	"fmt"
	"net/url"
	"testing"
)

var urls = []string{"http://abc.com", "https://abc.com/a/b/c", "https://abc.com?fas=234&faf=basf", "https://abc.com/a/?b=1&c=2"}

func TestUrlParse(t *testing.T) {
	for _, path := range urls {
		urlObj, _ := url.Parse(path)
		fmt.Println(urlObj.String())
	}
}

func TestHttpClient(t *testing.T) {
	res, err := util.HttpClient("https://baidu.com", "get", nil, nil)
	if err != nil {
		t.Log(err)
	}

	t.Log(string(res))
}
