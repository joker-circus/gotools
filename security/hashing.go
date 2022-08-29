/*散列加密算法*/
package security

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

func Sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

// Hmac: Hash-based Message Authentication Code
// 通过一个标准算法，在计算哈希的过程中，把key混入计算过程中

func HmacSha256(key, data []byte) string {
	return hmacSum(sha256.New, key, data)
}

func HmacSha1(key, data []byte) string {
	return hmacSum(sha1.New, key, data)
}

func HmacMd5(key, data []byte) string {
	return hmacSum(md5.New, key, data)
}

func hmacSum(h func() hash.Hash, key, data []byte) string {
	hashed := hmac.New(h, key)
	hashed.Write(data)
	return Base64EncodeToString(hashed.Sum(nil))
}
