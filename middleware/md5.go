package middleware

import (
	"crypto/md5"
	"encoding/hex"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempstr := h.Sum(nil)
	return hex.EncodeToString(tempstr)
}

// 加密
func MakePassWord(reqpwd, salt string) string {
	return Md5Encode(reqpwd + salt)
}

// 解密
func ValidPassWord(reqpwd, salt string, password string) bool {
	return Md5Encode(reqpwd+salt) == password
}
