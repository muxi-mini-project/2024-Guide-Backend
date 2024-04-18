package services

import (
	"crypto/sha256"
	"encoding/base64"
)

func hashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil)) + "$" + salt
}
