package httputil

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"

	"github.com/joker-circus/gotools/internal"
	"github.com/pkg/errors"
)

type ResponseValidator func(resp *http.Response) error

func EmptyValidator(resp *http.Response) error {
	return nil
}

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
	resp, err := Request(method, url, body, header, query, validators...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint
	return ReadRespBody(resp)
}

func Get(url string, query map[string]interface{}, header map[string]string, validators ...ResponseValidator) ([]byte, error) {
	resp, err := Request(http.MethodGet, url, nil, header, query, validators...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint
	return ReadRespBody(resp)
}

func DownFileByGet(url string, query map[string]interface{}, header map[string]string, validators ...ResponseValidator) error {
	return DownFile(http.MethodPost, url, nil, header, query, validators...)
}

func DownFileByPost(url string, body interface{}, header map[string]string, validators ...ResponseValidator) error {
	return DownFile(http.MethodPost, url, body, header, nil, validators...)
}

func DownFile(method, url string, body interface{}, header map[string]string, query map[string]interface{}, validators ...ResponseValidator) error {
	resp, err := Request(method, url, body, header, query, validators...)
	if err != nil {
		return errors.WithMessage(err, "Request")
	}

	defer resp.Body.Close() // nolint

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return errors.WithMessage(err, "ParseMediaType")
	}
	fileName, ok := params["filename"]
	if !ok {
		return errors.New("filename parameter not exist")
	}

	//创建文件
	file, err := os.Create(fileName)
	if err != nil {
		return errors.WithMessage(err, "create file")
	}

	// defer延迟调用 关闭文件，释放资源
	defer file.Close()

	//添加缓冲 bufio 是通过缓冲来提高效率。
	bufWriter := bufio.NewWriter(file)
	_, err = io.Copy(bufWriter, resp.Body)
	if err != nil {
		return errors.WithMessage(err, "io.Copy")
	}
	//将缓存的数据写入到文件中
	_ = bufWriter.Flush()
	return nil
}

func ReadRespBody(resp *http.Response) ([]byte, error) {
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
	req, err := clientRequest(method, url, body, header, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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
