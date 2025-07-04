package vkma

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"net/url"
	"sort"
	"strings"
	"sync"
)

var (
	// b64NoPad encodes HMAC using base64 without padding
	b64NoPad = base64.URLEncoding.WithPadding(base64.NoPadding)

	// keyedHMAC caches sync.Pool of HMAC hashers by secret key
	keyedHMAC = sync.Map{} // map[string]*sync.Pool
)

// getHMAC returns a hash.Hash from a pool for the given secret key.
// If no pool exists for the key, it is created on first use.
func getHMAC(secret string) hash.Hash {
	poolVal, ok := keyedHMAC.Load(secret)
	if !ok {
		newPool := &sync.Pool{
			New: func() any {
				return hmac.New(sha256.New, []byte(secret))
			},
		}
		keyedHMAC.Store(secret, newPool)
		return newPool.Get().(hash.Hash)
	}
	return poolVal.(*sync.Pool).Get().(hash.Hash)
}

// putHMAC returns the hash.Hash back to the pool after resetting.
func putHMAC(secret string, h hash.Hash) {
	h.Reset()
	if poolVal, ok := keyedHMAC.Load(secret); ok {
		poolVal.(*sync.Pool).Put(h)
	}
}

// Verify validates the VK Mini Apps query signature using HMAC.
// It returns a filled Params struct and a boolean indicating signature validity.
func Verify(rawQuery string, secrets map[string]string) (*Params, bool) {
	if len(secrets) == 0 {
		return nil, false
	}

	var (
		appID string
		sign  string
		seen  = make(map[string]string, 20)
	)

	// Parse query string manually to reduce overhead
	for start := 0; start < len(rawQuery); {
		end := strings.IndexByte(rawQuery[start:], '&')
		if end == -1 {
			end = len(rawQuery)
		} else {
			end += start
		}

		eq := strings.IndexByte(rawQuery[start:end], '=')
		if eq == -1 {
			start = end + 1
			continue
		}

		key, err1 := url.QueryUnescape(rawQuery[start : start+eq])
		val, err2 := url.QueryUnescape(rawQuery[start+eq+1 : end])
		if err1 != nil || err2 != nil {
			return nil, false
		}

		switch {
		case key == "sign":
			sign = val
		case key == "vk_app_id":
			appID = val
			seen[key] = val
		case strings.HasPrefix(key, "vk_"):
			seen[key] = val
		}

		start = end + 1
	}

	if appID == "" || sign == "" {
		return nil, false
	}

	secret, ok := secrets[appID]
	if !ok {
		return nil, false
	}

	// Sort keys for canonical HMAC string
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := make([]byte, 0, 512)
	var p Params
	for i, k := range keys {
		if i > 0 {
			buf = append(buf, '&')
		}
		buf = appendEscape(buf, k)
		buf = append(buf, '=')
		buf = appendEscape(buf, seen[k])
		p.set(k, seen[k])
	}

	// Use cached HMAC hasher
	mac := getHMAC(secret)
	defer putHMAC(secret, mac)

	mac.Write(buf)
	expectedSign := b64NoPad.EncodeToString(mac.Sum(nil))

	return &p, expectedSign == sign
}

// appendEscape performs percent-encoding of s into dst ([]byte), similar to url.QueryEscape.
func appendEscape(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			dst = append(dst, c)
		} else {
			dst = append(dst, '%')
			dst = append(dst, "0123456789ABCDEF"[c>>4])
			dst = append(dst, "0123456789ABCDEF"[c&15])
		}
	}
	return dst
}
