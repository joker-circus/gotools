package httputil

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
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

	"github.com/joker-circus/gotools/internal"
)

type ResponseValidator func(resp *http.Response) error

func ValidatorStatusCode(resp *http.Response, targetStatusCode int) error {
	if resp.StatusCode == targetStatusCode {
		return nil
	}
	return fmt.Errorf("the current status code is %d, but the expected status code is %d", resp.StatusCode, targetStatusCode)
}

func StatusOK(resp *http.Response) error {
	return ValidatorStatusCode(resp, http.StatusOK)
}

func Put(url string, body interface{}, header map[string]string, validators ...ResponseValidator) ([]byte, error) {
	return Do(http.MethodPut, url, body, header, make(map[string]interface{}), validators...)
}

func PostForm(url string, data url.Values, header map[string]string, validators ...ResponseValidator) ([]byte, error) {
	header["Content-Type"] = "application/x-www-form-urlencoded"
	body := internal.S2b(data.Encode())
	return Post(url, body, header, validators...)
}

func Post(url string, body interface{}, header map[string]string, validators ...ResponseValidator) ([]byte, error) {
	return Do(http.MethodPost, url, body, header, make(map[string]interface{}), validators...)
}

func Do(method, url string, body interface{}, header map[string]string, query map[string]interface{}, validators ...ResponseValidator) ([]byte, error) {
	resp, err := Request(method, url, body, header, make(map[string]interface{}), validators...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func Get(url string, query map[string]interface{}, header map[string]string, validators ...ResponseValidator) ([]byte, error) {
	resp, err := Request(http.MethodGet, url, nil, header, query, validators...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}
	defer reader.Close() // nolint

	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

// 若 body 为 string、[]byte 类型直接返回，
// 其他类型返回 json.Marshal(body)
func BytesBody(body interface{}) ([]byte, error) {
	if v, ok := body.(string); ok {
		return internal.S2b(v), nil
	}

	if v, ok := body.([]byte); ok {
		return v, nil
	}

	return json.Marshal(body)
}

// 如果 validators 为 null，默认验证 http code 是否为 200。
func Request(method, url string, body interface{}, header map[string]string, query map[string]interface{}, validators ...ResponseValidator) (*http.Response, error) {
	client := &http.Client{}

	byteParams, err := BytesBody(body)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequest(method, url, bytes.NewReader(byteParams))
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

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(validators) == 0 {
		validators = append(validators, StatusOK)
	}
	for _, validator := range validators {
		err = validator(resp)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return resp, nil
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
func ContentType(h http.Header) (string, error) {
	s := h.Get("Content-Type")
	// handle Content-Type like this: "text/html; charset=utf-8"
	return strings.Split(s, ";")[0], nil
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
