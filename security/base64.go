package security

import "encoding/base64"

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

func Base64DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}
