package vkmashop

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
)

// kv represents a key-value pair for query parameters
type kv struct {
	key string // parameter name
	val string // parameter value
}

// Verify validates the VK Mini Apps request signature and parses parameters
//
// Parameters:
//   - rawQuery: raw query string (after "?" in URL)
//   - secrets: mapping of app_id to application secret key
//
// Returns:
//   - *Params: parsed request parameters
//   - bool: true if signature is valid, false otherwise
//
// Algorithm:
//   1. Parse query string into key-value pairs
//   2. Extract app_id and signature (sig)
//   3. Lookup secret key using app_id
//   4. Sort parameters alphabetically by key
//   5. Concatenate parameters for signature verification
//   6. Compute and verify MD5 hash
//
// Notes:
//   - Skips parameters without values (key without =value)
//   - Automatically decodes URL-encoded values
//   - Returns false on any parsing error
//   - Follows VK's signature verification protocol:
//     https://dev.vk.com/mini-apps/development/launch-params
func Verify(rawQuery string, secrets map[string]string) (*Params, bool) {
	body := &Params{}
	pairs := make([]kv, 0, 20) // pre-allocate for 20 parameters

	var appID, sig string // temporary storage for app_id and signature
	var start, end, eq int // parsing indices

	// Parse query string
	for start < len(rawQuery) {
		end = start
		// Find end of current parameter (& or end of string)
		for end < len(rawQuery) && rawQuery[end] != '&' {
			end++
		}

		eq = start
		// Find key-value separator (=)
		for eq < end && rawQuery[eq] != '=' {
			eq++
		}
		// Skip parameters without values
		if eq == end {
			start = end + 1
			continue
		}

		// Extract raw (possibly URL-encoded) key and value
		rawK := rawQuery[start:eq]
		rawV := rawQuery[eq+1 : end]

		// URL decode key and value
		k, err1 := unescapeMinimal(rawK)
		v, err2 := unescapeMinimal(rawV)
		if err1 != nil || err2 != nil {
			return nil, false
		}

		// Store parameter in struct
		body.set(k, v)

		// Save app_id and sig for verification
		if k == "app_id" {
			appID = v
		}
		if k == "sig" {
			sig = v
		} else {
			// All parameters except sig are included in signature check
			pairs = append(pairs, kv{key: k, val: v})
		}

		start = end + 1
	}

	// Verify required parameters exist
	if appID == "" || sig == "" {
		return nil, false
	}
	// Get application secret key
	secret, ok := secrets[appID]
	if !ok {
		return nil, false
	}

	// Sort parameters alphabetically (VK API requirement)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].key < pairs[j].key
	})

	// Buffer for parameter concatenation
	buf := make([]byte, 0, len(rawQuery)+32)

	// Build signature verification string
	for _, pair := range pairs {
		buf = append(buf, pair.key...)
		buf = append(buf, '=')
		buf = append(buf, pair.val...)
	}
	// Append secret key
	buf = append(buf, secret...)

	// Compute MD5 hash
	sum := md5.Sum(buf)
	expected := hex.EncodeToString(sum[:])

	// Compare with received signature
	return body, expected == sig
}

// unescapeMinimal decodes URL-encoded strings only when necessary
//
// Optimization: avoids memory allocations when no decoding is needed
//
// Returns:
//   - string: decoded string
//   - error: decoding error if any occurred
func unescapeMinimal(s string) (string, error) {
	// Fast check for escape sequences
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '%', '+':
			return url.QueryUnescape(s)
		}
	}
	return s, nil
}