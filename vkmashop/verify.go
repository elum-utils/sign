package vkmashop

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

type kv struct {
	key string
	val string
}

func Verify(rawQuery string, secrets map[string]string) (*Body, bool) {
	body := &Body{}
	pairs := make([]kv, 0, 20)

	var appID, sig string

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

		k, err1 := url.QueryUnescape(rawQuery[start : start+eq])
		v, err2 := url.QueryUnescape(rawQuery[start+eq+1 : end])
		if err1 != nil || err2 != nil {
			return nil, false
		}

		body.set(k, v)

		if k == "app_id" {
			appID = v
		}
		if k == "sig" {
			sig = v
		} else {
			pairs = append(pairs, kv{key: k, val: v})
		}

		start = end + 1
	}

	if appID == "" || sig == "" {
		return nil, false
	}

	secret, ok := secrets[appID]
	if !ok {
		return nil, false
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].key < pairs[j].key
	})

	buf := make([]byte, 0, len(rawQuery)+len(secret))
	for _, pair := range pairs {
		buf = append(buf, pair.key...)
		buf = append(buf, '=')
		buf = append(buf, pair.val...)
	}
	buf = append(buf, secret...)

	sum := md5.Sum(buf)
	expected := hex.EncodeToString(sum[:])

	return body, expected == sig
}
