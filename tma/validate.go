package tma

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"sync"
)

var (
	hashSecret []byte
	bufPool    = sync.Pool{
		New: func() any { return make([]byte, 0, 256) },
	}
	// json       = jsoniter.ConfigCompatibleWithStandardLibrary
)

type TMAUser struct {
	ID                    int    `msgpack:"ni" json:"id"`
	FirstName             string `msgpack:"fn" json:"first_name"`
	LastName              string `msgpack:"ln" json:"last_name"`
	UserName              string `msgpack:"un" json:"username"`
	PhotoURL              string `msgpack:"pu" json:"photo_url"`
	Language              string `msgpack:"lc" json:"language_code"`
	ChatType              string `msgpack:"ct" json:"chat_type"`
	ChatInstance          string `msgpack:"ci" json:"chat_instance"`
	IsPremium             bool   `msgpack:"pr" json:"is_premium"`
	AllowsWriteToPM       bool   `msgpack:"aw" json:"allows_write_to_pm"`
	AddedToAttachmentMenu bool   `msgpack:"am" json:"added_to_attachment_menu"`
}

func Validate(params, secret string) (*TMAUser, bool) {

	if secret == "" {
		return nil, false
	}

	if hashSecret == nil {
		h := hmac.New(sha256.New, []byte("WebAppData"))
		h.Write([]byte(secret))
		hashSecret = h.Sum(nil)
	}

	var (
		hash  string
		pairs = make([][2]string, 0, 10)
	)

	// Ручной парсинг query string
	for start := 0; start < len(params); {
		end := strings.IndexByte(params[start:], '&')
		if end == -1 {
			end = len(params)
		} else {
			end += start
		}

		eq := strings.IndexByte(params[start:end], '=')
		if eq == -1 {
			start = end + 1
			continue
		}

		rawKey := params[start : start+eq]
		rawVal := params[start+eq+1 : end]

		key, err1 := url.QueryUnescape(rawKey)
		val, err2 := url.QueryUnescape(rawVal)
		if err1 != nil || err2 != nil {
			return &TMAUser{}, false
		}

		if key == "hash" {
			hash = val
		} else {
			pairs = append(pairs, [2]string{key, val})
		}

		start = end + 1
	}

	if hash == "" {
		return &TMAUser{}, false
	}

	// Сортировка пар по ключу
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})

	// Сборка строки для HMAC
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf[:0])

	for i, p := range pairs {
		if i > 0 {
			buf = append(buf, '\n')
		}
		buf = append(buf, p[0]...)
		buf = append(buf, '=')
		buf = append(buf, p[1]...)
	}

	mac := hmac.New(sha256.New, hashSecret)
	mac.Write(buf)
	computedHash := mac.Sum(nil)

	decodedHash, err := hex.DecodeString(hash)
	if err != nil || !hmac.Equal(computedHash, decodedHash) {
		return &TMAUser{}, false
	}

	

	return &TMAUser{}, true

}
