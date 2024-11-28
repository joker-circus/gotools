package datatypes

import (
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type DBZipData[T any] []T

func (za *DBZipData[T]) FromZipByte(bs []byte) error {
	if bs == nil {
		return nil
	}

	b, err := UnGzip(bs)
	if err != nil {
		return err
	}

	var znf DBZipData[T]
	err = json.Unmarshal(b, &znf)
	if err != nil {
		return err
	}
	*za = znf
	return nil
}

func (za DBZipData[T]) ToZipBytes() ([]byte, error) {
	if za == nil {
		return nil, nil
	}
	data, err := json.Marshal(za)
	if err != nil {
		return nil, err
	}

	return Gzip(data)
}

func (za DBZipData[T]) Value() (driver.Value, error) {
	return za.ToZipBytes()
}

func (za *DBZipData[T]) Scan(v interface{}) error {
	var val DBZipData[T]
	var fields T
	if v == nil {
		*za = val
		return nil
	}

	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("can not scan value %v to %T", v, fields)
	}

	if len(b) == 0 {
		*za = val
		return nil
	}

	return za.FromZipByte(b)
}

func Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	//result, err := io.ReadAll(r)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Base64Decode(src []byte) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(buf, src)
	return buf[:n], err
}

func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	return buf
}
