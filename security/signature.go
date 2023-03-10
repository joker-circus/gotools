package security

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

// API 签名方法：
// 	key - 私钥或 token
//	nonce - 随机数
// 	timestamp - 时间戳
type APISignatureFunc func(key, nonce, timestamp string) string

// 获取签名随机数、及当前时间戳
func GetNonceAndTimeStamp() (nonce string, timestamp string) {
	ts := time.Now().Unix()
	rand.Seed(ts)
	nonce = strconv.FormatInt(rand.Int63(), 10)
	timestamp = strconv.FormatInt(ts, 10)
	return
}

// 便捷签名方法
// 	key - 私钥或 token
//	signatureFunc - 生成签名的方法
func MakeSignature(key string, signatureFunc APISignatureFunc) string {
	nonce, timestamp := GetNonceAndTimeStamp()
	return signatureFunc(key, nonce, timestamp)
}

// MakeSignature sign with the given appKey, random string and timestamp string.
func Sha256Signature(appkey, randStr, timestampStr string) string {
	hash := sha256.New()
	hash.Write([]byte(appkey))
	hash.Write([]byte(randStr))
	hash.Write([]byte(timestampStr))
	return hex.EncodeToString(hash.Sum(nil))
}

// more see: https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Getting_Started_Guide.html
func WeChatPlatformSign(token, timestamp, nonce string) string {
	data := []string{token, timestamp, nonce}
	sort.Strings(data)

	encrypt := sha1.New()
	encrypt.Write([]byte(strings.Join(data, "")))
	return hex.EncodeToString(encrypt.Sum(nil))
}
