package models

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"math/rand"
	"time"
)

var DB *sql.DB

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func encodePassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func CreateRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
