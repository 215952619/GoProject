package util

import (
	"GoProject/global"
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"time"
)

type RequestMethod string
type UpdateHeadFunc func(header *http.Header)

func HttpClient(url string, method string, params []byte, handler UpdateHeadFunc) ([]byte, error) {
	parse, err := url2.Parse(url)
	if err != nil {
		return nil, errors.New("invalid url path")
	}

	client := &http.Client{
		Timeout: 100 * time.Second,
	}
	if global.Mode == "dev" {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		if parse.Hostname() == "github.com" {
			transport.Proxy = func(_ *http.Request) (*url2.URL, error) {
				return url2.Parse("https://www.envbits.com:443")
			}
		}

		client.Transport = transport
	}

	var req *http.Request
	if method == http.MethodPost {
		req, err = http.NewRequest(http.MethodPost, parse.String(), bytes.NewBuffer(params))
	} else {
		req, err = http.NewRequest(http.MethodGet, parse.String(), nil)
	}
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":    err,
			"method": method,
			"url":    parse.String(),
			"params": params,
		}).Error("make request error")
		return nil, err
	}

	if handler != nil {
		handler(&req.Header)
	}

	resp, err := client.Do(req)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":    err.Error(),
			"method": method,
			"url":    parse.String(),
			"params": params,
			"resp":   resp,
		}).Error("receive response error")
		return nil, err
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":      err,
			"method":   method,
			"url":      parse.String(),
			"response": response,
		}).Error("read response error")
		return nil, err
	}
	return response, nil
}

func GetRequest(url string, handler UpdateHeadFunc) ([]byte, error) {
	return HttpClient(url, http.MethodGet, nil, handler)
}

func PostRequest(url string, params []byte, handler UpdateHeadFunc) ([]byte, error) {
	return HttpClient(url, http.MethodPost, params, handler)
}
