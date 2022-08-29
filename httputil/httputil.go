package httputil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, body interface{}, header map[string]string) ([]byte, error) {
	return Do(http.MethodPost, url, body, header, make(map[string]interface{}))
}

func Get(url string, query map[string]interface{}, header map[string]string) ([]byte, error) {
	return Do(http.MethodGet, url, nil, header, query)
}

// 若 body 为 string、[]byte 类型直接返回，
// 其他类型返回 json.Marshal(body)
func BytesBody(body interface{}) ([]byte, error) {
	if v, ok := body.(string); ok {
		return []byte(v), nil
	}

	if v, ok := body.([]byte); ok {
		return v, nil
	}

	return json.Marshal(body)
}

func Do(method, url string, body interface{}, header map[string]string, query map[string]interface{}) ([]byte, error) {
	client := &http.Client{}

	byteParams, err := BytesBody(body)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequest(method, url, strings.NewReader(string(byteParams)))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	defer resp.Body.Close()

	repBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp.StatusCode is %v, resp body is %+v, req body is %+v", resp.StatusCode, string(repBody), string(byteParams))
	}

	return repBody, err
}
