package handler

import (
	_"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

type shortcode struct {
	processedUrl string
	originalurl string
	expiresDate time.Time
	isCustom bool
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ShorteningUrl(originalurl, alias string, expire time.Time) *shortcode {
	sc := &shortcode{}
	if alias != "" {	
		sc.isCustom = true
    	sc.originalurl = originalurl
		sc.processedUrl = alias
		sc.expiresDate = expire
    	return sc
	}

	hasher := sha256.New()
	hasher.Write([]byte(originalurl))
	hashed := hex.EncodeToString(hasher.Sum(nil))

	var seed int64
	for _, c := range hashed[:8] {
		seed = seed * 31 + int64(c)
	}

	r := rand.New(rand.NewSource(seed * time.Now().UnixNano()))
	b := make([]rune, 8)
	for i := range b {
		b[i] =letterRunes[r.Intn(len(letterRunes))]
	}

	sc.originalurl = originalurl
	sc.processedUrl = string(b)
	sc.isCustom = false
	sc.expiresDate = expire

	return sc
}







