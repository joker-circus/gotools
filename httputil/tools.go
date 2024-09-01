package httputil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
)

func clientRequest(method, url string, body interface{}, header map[string]string, query map[string]interface{}) (*http.Request, error) {
	byteParams, err := BytesBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(byteParams))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = q.Encode()

	req.Close = true
	return req, nil
}

// Size get size of the header
func Size(h http.Header) (int64, error) {
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not present")
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// ContentType get Content-Type of the header
func ContentType(h http.Header) string {
	s := h.Get("Content-Type")
	// handle Content-Type like this: "text/html; charset=utf-8"
	return strings.Split(s, ";")[0]
}

// M3u8URLs get all urls from m3u8 url
func M3u8URLs(uri string) ([]string, error) {
	if len(uri) == 0 {
		return nil, errors.New("url is null")
	}

	html, err := Get(uri, nil, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	lines := strings.Split(string(html), "\n")
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "http") {
				urls = append(urls, line)
			} else {
				base, err := url.Parse(uri)
				if err != nil {
					continue
				}
				u, err := url.Parse(line)
				if err != nil {
					continue
				}
				urls = append(urls, base.ResolveReference(u).String())
			}
		}
	}
	return urls, nil
}

// 获取 body 的编码格式
func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

// 返回 UTF-8 的 html 文本。
// DecodeHTMLBody returns an decoding reader of the html Body for the specified `charset`
// If `charset` is empty, DecodeHTMLBody tries to guess the encoding from the content
func DecodeHTMLBody(body io.Reader, charsets ...string) (io.Reader, error) {
	var charsetName string
	if len(charsets) == 0 {
		charsetName = detectContentCharset(body)
	} else {
		charsetName = charsets[0]
	}

	e, err := htmlindex.Get(charsetName)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}
	return body, nil
}

func JoinPaths(uri string, paths ...string) string {
	urlObj, _ := url.Parse(uri)
	newPath := make([]string, 0, len(paths)+1)
	newPath = append(newPath, urlObj.Path)
	newPath = append(newPath, paths...)
	urlObj.Path = path.Join(newPath...)

	return urlObj.String()
}

// 展示 response body 内容
func ShowResponseBody(resp *http.Response) (string, error) {
	if resp == nil || resp.Body == nil {
		return "", nil
	}

	// 拿到body字节流数据
	b, _ := ioutil.ReadAll(resp.Body)

	// 用该方法继续将数据写入Body用于传递
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return string(b), nil
}

// 通过模拟 request 请求方式，展示 json body 内容
func ShowMockRequestBody(method, url string, body interface{}, header map[string]string, query map[string]interface{}) (string, error) {
	req, err := clientRequest(method, url, body, header, query)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return ShowRequestBody(req)
}

// 展示 request json body 内容
func ShowRequestBody(req *http.Request) (string, error) {
	return ShowRequestBodyByJson(req)
}

// 展示 request KV 形式 body 内容
func ShowRequestBodyByKV(req *http.Request) (string, error) {
	return showRequestBody(req, showKVForm)
}

// 展示 request json body 内容
func ShowRequestBodyByJson(req *http.Request) (string, error) {
	return showRequestBody(req, showJsonForm)
}

func showRequestBody(req *http.Request, showForm func(map[string][]string) string) (string, error) {
	if req == nil {
		return "", nil
	}

	switch ContentType(req.Header) {
	case "application/x-www-form-urlencoded":
		err := req.ParseForm()
		if err != nil {
			return "", err
		}
		// req.Form 比 req.PostForm 多了 url 参数
		return showForm(req.PostForm), nil

	case "multipart/form-data":
		// 10 M
		err := req.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			return "", err
		}
		if req.MultipartForm == nil {
			return "", nil
		}
		form := req.MultipartForm.Value
		for k, values := range req.MultipartForm.File {
			for _, v := range values {
				if v == nil {
					continue
				}
				form[k] = append(form[k], v.Filename)
			}
		}
		return showForm(form), nil
	default:
		if req.Body == nil {
			return "", nil
		}

		// 拿到body字节流数据
		b, _ := ioutil.ReadAll(req.Body)

		// 用该方法继续将数据写入Body用于传递
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		return string(b), nil
	}
}

func showKVForm(data map[string][]string) string {
	var res []string
	for k, values := range data {
		for _, v := range values {
			res = append(res, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return strings.Join(res, "\n")
}

func showJsonForm(data map[string][]string) string {
	res := make(map[string]interface{}, len(data))
	for k, values := range data {
		switch len(values) {
		case 0:
		case 1:
			res[k] = values[0]
		default:
			res[k] = values
		}
	}
	b, _ := json.Marshal(res)
	return string(b)
}
