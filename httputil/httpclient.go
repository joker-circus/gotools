package httputil

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	retryAlgorithmIncrement = iota
	retryAlgorithmRandom
	retryAlgorithmEqualWait
)

type (
	judgeFunc func(*HttpResp) error

	HttpClient struct {
		cli      http.Client
		req      *http.Request
		r        *retry
		resp     *http.Response
		header   map[string]string
		err      error
		httpResp *HttpResp
	}

	retry struct {
		times     int
		algorithm int
		sec       int
		judge     judgeFunc
	}

	HttpResp struct {
		StatusCode int
		Body       []byte
	}
)

func NewHTTPClient() *HttpClient {
	return &HttpClient{
		header: make(map[string]string),
	}
}

func (h *HttpClient) Request() *http.Request {
	return h.req
}

func (h *HttpClient) Response() *http.Response {
	return h.resp
}

func (h *HttpClient) WriteHeader(k, v string) *HttpClient {

	h.header[k] = v
	return h
}

func (h *HttpClient) IncRetry(times int, baseSeconds int, j judgeFunc) *HttpClient {

	h.r = &retry{
		times:     times,
		algorithm: retryAlgorithmIncrement,
		sec:       baseSeconds,
		judge:     j,
	}
	return h
}

func (h *HttpClient) RandomRetry(times int, maxSeconds int, j judgeFunc) *HttpClient {

	h.r = &retry{
		times:     times,
		algorithm: retryAlgorithmRandom,
		sec:       maxSeconds,
		judge:     j,
	}
	return h
}

func (h *HttpClient) EqualRetry(times int, waitSeconds int, j judgeFunc) *HttpClient {

	h.r = &retry{
		times:     times,
		algorithm: retryAlgorithmEqualWait,
		sec:       waitSeconds,
		judge:     j,
	}
	return h
}

func (h *HttpClient) doWithRetry() {

	i := 0
	resp := &HttpResp{}

	defer func() {
		if h.err == nil {
			h.httpResp = resp
		}
	}()

	if h.r.algorithm == retryAlgorithmEqualWait {
		for {
			i++

			h.resp, h.err = h.cli.Do(h.req)
			h.getResult(resp)

			if h.err == nil && h.r.judge(resp) == nil {
				return
			}

			if i >= h.r.times {
				break
			}

			time.Sleep(time.Duration(h.r.sec) * time.Second)
		}

		if h.err == nil {
			h.err = h.r.judge(resp)
		}

		h.err = fmt.Errorf("%s, retry %d time(s)", h.err.Error(), h.r.times)
		return
	}

	if h.r.algorithm == retryAlgorithmIncrement {
		waitSeconds := h.r.sec
		for {
			i++
			h.resp, h.err = h.cli.Do(h.req)
			h.getResult(resp)

			if h.err == nil && h.r.judge(resp) == nil {
				return
			}

			if i >= h.r.times {
				break
			}

			time.Sleep(time.Duration(waitSeconds) * time.Second)
			waitSeconds = waitSeconds << 1
		}

		if h.err == nil {
			h.err = h.r.judge(resp)
		}

		h.err = fmt.Errorf("%s, retry %d time(s)", h.err.Error(), h.r.times)
		return
	}

	if h.r.algorithm == retryAlgorithmRandom {
		for {
			i++

			h.resp, h.err = h.cli.Do(h.req)
			h.getResult(resp)

			if h.err == nil && h.r.judge(resp) == nil {
				return
			}

			rand.Seed(time.Now().Unix())
			waitSeconds := rand.Int() % h.r.sec
			if waitSeconds == 0 {
				waitSeconds = 1 // minimal is 1
			}

			if i >= h.r.times {
				break
			}

			time.Sleep(time.Duration(waitSeconds) * time.Second)
		}

		if h.err == nil {
			h.err = h.r.judge(resp)
		}

		h.err = fmt.Errorf("%s, retry %d time(s)", h.err.Error(), h.r.times)
		return
	}

}

func (h *HttpClient) doWriteHeader() {

	if h.req == nil {
		return
	}

	for k, v := range h.header {
		h.req.Header.Set(k, v)
	}
}

func (h *HttpClient) Timeout(second int) *HttpClient {
	h.cli.Timeout = time.Duration(second) * time.Second

	return h
}

func (h *HttpClient) Get(url string) *HttpClient {

	return h.Do("GET", url, nil)
}

func (h *HttpClient) Post(url string, body string) *HttpClient {

	return h.Do("POST", url, strings.NewReader(body))

}

func (h *HttpClient) PostBytes(url string, body []byte) *HttpClient {

	return h.Do("POST", url, bytes.NewReader(body))

}

func (h *HttpClient) Do(method, url string, body io.Reader) *HttpClient {

	h.req, h.err = http.NewRequest(strings.ToUpper(method), url, body)
	h.doWriteHeader()

	if h.r != nil {
		h.doWithRetry()
	} else {
		h.resp, h.err = h.cli.Do(h.req)
	}

	return h
}

func (h *HttpClient) Result(resp *HttpResp) *HttpClient {
	if h.httpResp != nil {
		resp.StatusCode = h.httpResp.StatusCode
		resp.Body = h.httpResp.Body
	}

	return h.getResult(resp)
}

func (h *HttpClient) getResult(resp *HttpResp) *HttpClient {

	if h.err != nil {
		return h
	}

	if h.httpResp != nil {
		resp.StatusCode = h.httpResp.StatusCode
		resp.Body = h.httpResp.Body
		return h
	}

	if h.resp == nil {
		h.err = errors.New("nil response, please do a request first. ")
		return h
	}

	defer h.resp.Body.Close()
	body, err := ioutil.ReadAll(h.resp.Body)

	if err != nil {
		h.err = errors.New("nil response, please do a request first. ")
		return h
	}

	resp.StatusCode = h.resp.StatusCode
	resp.Body = body

	return h
}

func (h *HttpClient) Error() error {
	return h.err
}

func (h *HttpClient) SetBasicAuth(username, password string) *HttpClient {
	auth := username + ":" + password
	return h.WriteHeader("Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
}

func (h *HttpClient) SetOAuth2(token string) *HttpClient {
	return h.WriteHeader("Authorization",
		"Bearer "+token)
}
