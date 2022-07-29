package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//type UpdateHeadFunc func(header *http.Header)
//
//func DefaultClient(url *url2.URL) *http.Client {
//	client := &http.Client{
//		Timeout: 100 * time.Second,
//	}
//	if global.Mode == "dev" {
//		transport := &http.Transport{
//			TLSClientConfig: &tls.Config{
//				InsecureSkipVerify: true,
//			},
//		}
//
//		if url.Hostname() == "github.com" {
//			transport.Proxy = func(_ *http.Request) (*url2.URL, error) {
//				return url2.Parse("https://www.envbits.com:443")
//			}
//		}
//		client.Transport = transport
//	}
//	return client
//}
//
//func HttpClient(url string, method string, params []byte, handler UpdateHeadFunc) ([]byte, error) {
//	parse, err := url2.Parse(url)
//	if err != nil {
//		return nil, errors.New("invalid url path")
//	}
//
//	var req *http.Request
//	if method == http.MethodPost {
//		req, err = http.NewRequest(http.MethodPost, parse.String(), bytes.NewBuffer(params))
//	} else {
//		req, err = http.NewRequest(http.MethodGet, parse.String(), nil)
//	}
//	if err != nil {
//		global.Logger.WithFields(logrus.Fields{
//			"err":    err,
//			"method": method,
//			"url":    parse.String(),
//			"params": params,
//		}).Error("make request error")
//		return nil, err
//	}
//
//	if handler != nil {
//		handler(&req.Header)
//	}
//
//	resp, err := DefaultClient(parse).Do(req)
//	if err != nil {
//		global.Logger.WithFields(logrus.Fields{
//			"err":    err.Error(),
//			"method": method,
//			"url":    parse.String(),
//			"params": params,
//			"resp":   resp,
//		}).Error("receive response error")
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//	response, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		global.Logger.WithFields(logrus.Fields{
//			"err":      err,
//			"method":   method,
//			"url":      parse.String(),
//			"response": response,
//		}).Error("read response error")
//		return nil, err
//	}
//	return response, nil
//}
//
//func GetRequest(url string, params map[string]string, header map[string]string) ([]byte, error) {
//	httpclient := &HttpSend{
//		Link:   url,
//		Params: params,
//		Header: header,
//	}
//	return httpclient.Get()
//}
//
//func PostJsonRequest(url string, dataType string, params map[string]string, header map[string]string, data map[string]string) ([]byte, error) {
//	httpclient := &HttpSend{
//		Link:     url,
//		SendType: dataType,
//		Params:   params,
//		Header:   header,
//		Body:     data,
//	}
//	return httpclient.Post()
//}
//
//var (
//	GetMethod    = "GET"
//	PostMethod   = "POST"
//	SendTypeFrom = "from"
//	SendTypeJson = "json"
//)
//
//type HttpSend struct {
//	Link     string
//	SendType string
//	Params   map[string]string
//	Header   map[string]string
//	Body     map[string]string
//	FullPath string
//	sync.RWMutex
//}
//
//func (h *HttpSend) SetBody(body map[string]string) {
//	h.Lock()
//	defer h.Unlock()
//	h.Body = body
//}
//func (h *HttpSend) SetHeader(header map[string]string) {
//	h.Lock()
//	defer h.Unlock()
//	h.Header = header
//}
//func (h *HttpSend) SetSendType(sendType string) {
//	h.Lock()
//	defer h.Unlock()
//	h.SendType = sendType
//}
//func (h *HttpSend) Get() ([]byte, error) {
//	return h.send(GetMethod)
//}
//func (h *HttpSend) Post() ([]byte, error) {
//	return h.send(PostMethod)
//}
//func (h *HttpSend) Client() *http.Client {
//	if h.FullPath == "" {
//		return nil
//	}
//	client := &http.Client{
//		Timeout: 100 * time.Second,
//	}
//	if global.Mode == "dev" {
//		transport := &http.Transport{
//			TLSClientConfig: &tls.Config{
//				InsecureSkipVerify: true,
//			},
//		}
//		if strings.Contains(h.FullPath, "github.com") {
//			transport.Proxy = func(_ *http.Request) (*url.URL, error) {
//				return url.Parse("https://www.envbits.com:443")
//			}
//		}
//		client.Transport = transport
//	}
//	return client
//}
//func (h *HttpSend) ParseUrl() {
//	u, _ := url.Parse(h.Link)
//	q := u.Query()
//	for k, v := range h.Params {
//		q.Set(k, v)
//	}
//	u.RawQuery = q.Encode()
//	h.FullPath = u.String()
//}
//func (h *HttpSend) ParseData() (string, error) {
//	if len(h.Body) > 0 {
//		if strings.ToLower(h.SendType) == SendTypeJson {
//			sendBody, jsonErr := json.Marshal(h.Body)
//			if jsonErr != nil {
//				return "", jsonErr
//			}
//			return string(sendBody), nil
//		} else {
//			sendBody := http.Request{}
//			sendBody.ParseForm()
//			for k, v := range h.Body {
//				sendBody.Form.Add(k, v)
//			}
//			return sendBody.Form.Encode(), nil
//		}
//	}
//	return "", nil
//}
//func (h *HttpSend) ParseHeader(req *http.Request) {
//	//设置默认header
//	if len(h.Header) == 0 {
//		if strings.ToLower(h.SendType) == SendTypeJson {
//			h.Header = map[string]string{
//				"Content-Type": "application/json; charset=utf-8",
//			}
//		} else {
//			h.Header = map[string]string{
//				"Content-Type": "application/x-www-form-urlencoded",
//			}
//		}
//	}
//	for k, v := range h.Header {
//		if strings.ToLower(k) == "host" {
//			req.Host = v
//		} else {
//			req.Header.Add(k, v)
//		}
//	}
//}
//func (h *HttpSend) send(method string) ([]byte, error) {
//	h.ParseUrl()
//	data, err := h.ParseData()
//	if err != nil {
//		return nil, err
//	}
//	req, err := http.NewRequest(method, h.FullPath, strings.NewReader(data))
//	if err != nil {
//		return nil, err
//	}
//	defer req.Body.Close()
//	h.ParseHeader(req)
//
//	resp, err := h.Client().Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, errors.New(fmt.Sprintf("error http code :%d", resp.StatusCode))
//	}
//	return ioutil.ReadAll(resp.Body)
//}

type UploadFile struct {
	// 表单名称
	Name string
	// 文件全路径
	Filepath string
}

// 请求客户端
var httpClient = &http.Client{}

func GetRequest(reqUrl string, reqParams map[string]string, headers map[string]string) ([]byte, error) {
	urlParams := url.Values{}
	Url, err := url.Parse(reqUrl)
	if err != nil {
		return nil, nil
	}
	for key, val := range reqParams {
		urlParams.Set(key, val)
	}

	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = urlParams.Encode()
	// 得到完整的url，http://xx?query
	urlPath := Url.String()

	httpRequest, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, nil
	}
	// 添加请求头
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return getResponse(*resp, respData)
}

func PostFormRequest(reqUrl string, reqParams map[string]string, headers map[string]string) ([]byte, error) {
	return post(reqUrl, reqParams, "application/x-www-form-urlencoded", nil, headers)
}

func PostJsonRequest(reqUrl string, reqParams map[string]string, headers map[string]string) ([]byte, error) {
	return post(reqUrl, reqParams, "application/json", nil, headers)
}

func PostFileRequest(reqUrl string, reqParams map[string]string, files []UploadFile, headers map[string]string) ([]byte, error) {
	return post(reqUrl, reqParams, "multipart/form-data", files, headers)
}

func post(reqUrl string, reqParams map[string]string, contentType string, files []UploadFile, headers map[string]string) ([]byte, error) {
	requestBody, realContentType, err := getReader(reqParams, contentType, files)
	if err != nil {
		return nil, nil
	}
	httpRequest, err := http.NewRequest("POST", reqUrl, requestBody)
	// 添加请求头
	httpRequest.Header.Add("Content-Type", realContentType)
	if headers != nil {
		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}
	}
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return getResponse(*resp, respData)
}

func getReader(reqParams map[string]string, contentType string, files []UploadFile) (io.Reader, string, error) {
	if strings.Index(contentType, "json") > -1 {
		bytesData, err := json.Marshal(reqParams)
		return bytes.NewReader(bytesData), contentType, err
	} else if files != nil {
		body := &bytes.Buffer{}
		// 文件写入 body
		writer := multipart.NewWriter(body)
		for _, uploadFile := range files {
			file, err := os.Open(uploadFile.Filepath)
			if err != nil {
				return nil, "", err
			}
			part, err := writer.CreateFormFile(uploadFile.Name, filepath.Base(uploadFile.Filepath))
			if err != nil {
				return nil, "", err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, "", err
			}
			file.Close()
		}
		// 其他参数列表写入 body
		for k, v := range reqParams {
			if err := writer.WriteField(k, v); err != nil {
				return nil, "", err
			}
		}
		if err := writer.Close(); err != nil {
			return nil, "", err
		}
		// 上传文件需要自己专用的contentType
		return body, writer.FormDataContentType(), nil
	} else {
		urlValues := url.Values{}
		for key, val := range reqParams {
			urlValues.Set(key, val)
		}
		reqBody := urlValues.Encode()
		return strings.NewReader(reqBody), contentType, nil
	}
}

func getResponse(res http.Response, resData []byte) ([]byte, error) {
	contentType := res.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		var result = map[string]string{}
		for _, item := range strings.Split(string(resData), "&") {
			val := strings.Split(item, "=")
			value, err := url.QueryUnescape(val[1])
			if err != nil {
				return nil, err
			}
			if result[val[0]] != "" {
				result[val[0]] = fmt.Sprintf("%s,%s", result[val[0]], value)
			} else {
				result[val[0]] = value
			}
		}
		out, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		return out, err
	}
	return resData, nil
}
