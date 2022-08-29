package security

import "encoding/base64"

func Base64DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Base64EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}
