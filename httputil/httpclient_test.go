package httputil

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestHttpClient_Timeout(t *testing.T) {

	timeout := 3
	t1 := time.Now()

	NewHTTPClient().Timeout(timeout).Get("http://google.com")

	duration := time.Since(t1)

	if int(duration.Seconds()) > timeout {
		t.Error("client running exceed timeout. ")
	}
}

func TestHttpClient_WriteHeader(t *testing.T) {

	cli := NewHTTPClient().
		Timeout(1).
		WriteHeader("test-header-1", "test1").
		WriteHeader("test-header-2", "test2").
		Get("http://google.com")

	v, ok := cli.Request().Header["Test-Header-1"]

	if !ok {
		t.Logf("header: %+v", cli.Request().Header)
		t.Error("client request header not been set. ")
		return
	}

	if len(v) < 1 || v[0] != "test1" {
		t.Error("client request header not match as set. ")
	}
}

func TestHttpClient_Result(t *testing.T) {

	resp := HttpResp{}
	err := NewHTTPClient().
		Timeout(5).
		WriteHeader("test-header-1", "test1").
		Get("https://postman-echo.com/get").
		getResult(&resp).Error()

	if err != nil {
		t.Errorf("request echo service fail: %s", err.Error())
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("request echo service return non-ok code: %d", resp.StatusCode)
		return
	}

	data := struct {
		Args    map[string]string `json:"args"`
		Headers map[string]string `json:"headers"`
		Url     string            `json:"url"`
	}{}

	err = json.Unmarshal(resp.Body, &data)
	if err != nil {
		t.Errorf("parse body data fail: %s", err.Error())
		return
	}

	if v, ok := data.Headers["test-header-1"]; !ok || v != "test1" {
		t.Logf("headers:  %+v", data.Headers)
		t.Errorf("header of response check fail: ok: %t, value: %s", ok, v)
		return
	}

}

func TestHttpClient_Retry(t *testing.T) {

	resp := HttpResp{}
	err := NewHTTPClient().
		Timeout(3).
		WriteHeader("test-header-1", "test1").
		RandomRetry(3, 1, func(resp *HttpResp) error {
			if resp.StatusCode != 200 {
				return errors.New("http code error")
			}
			return nil
	}).
		Get("http://apigw-sre-dev.ucloudadmin.com/south-gate/v1/groups").
		getResult(&resp).Error()

	if err != nil {
		t.Errorf("request echo service fail: %s", err.Error())
		return
	}
}
