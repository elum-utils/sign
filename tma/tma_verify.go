// Package tma provides functionality for verifying and processing Telegram Mini Apps
// initialization data. It handles the secure validation of incoming web app data
// using HMAC-SHA256 signatures.
package tma

import (
	"crypto/hmac"
	"crypto/sha256"
	"strings"

	"github.com/elum-utils/sign/internal/utils"
)

// hashSecret stores the precomputed HMAC key derived from the secret.
// It's lazily initialized on first use to avoid unnecessary computation.
var hashSecret []byte

// Verify validates the raw query string from Telegram Mini Apps initialization
// against the provided secret key using HMAC-SHA256 signature verification.
//
// Parameters:
//   - rawQuery: The URL-encoded query string received from Telegram
//   - secret: The application secret key provided by Telegram Bot API
//
// Returns:
//   - *Params: Parsed parameters if verification succeeds
//   - bool: Verification result (true if valid)
//
// The function performs several security checks:
//   1. Validates input parameters aren't empty
//   2. Parses and sorts query parameters
//   3. Computes HMAC-SHA256 signature
//   4. Compares with provided signature
//   5. Returns parsed parameters only if verification succeeds
func Verify(rawQuery, secret string) (*Params, bool) {
	// Early return for empty inputs
	if secret == "" || rawQuery == "" {
		return nil, false
	}

	// Lazy initialization of HMAC key
	if hashSecret == nil {
		h := hmac.New(sha256.New, []byte("WebAppData"))
		h.Write([]byte(secret))
		hashSecret = h.Sum(nil)
	}

	var hash string

	// Get key-value pairs from pool to avoid allocations
	pairsPtr := utils.KVPool.Get().(*utils.KVSlice)
	pairs := (*pairsPtr)[:0] // Slice reset without reallocation
	defer utils.KVPool.Put(pairsPtr)

	// Get temporary buffer from pool for unescaping
	tmpBufPtr := utils.TmpBufPool.Get().(*[]byte)
	tmpBuf := *tmpBufPtr
	tmpBuf = tmpBuf[:0]
	defer utils.TmpBufPool.Put(tmpBufPtr)

	// Parse query string
	for start := 0; start < len(rawQuery); {
		// Find next parameter boundary
		end := strings.IndexByte(rawQuery[start:], '&')
		if end == -1 {
			end = len(rawQuery)
		} else {
			end += start
		}

		// Split key-value pair
		eq := strings.IndexByte(rawQuery[start:end], '=')
		if eq == -1 {
			return nil, false // Malformed parameter
		}
		eq += start

		// Unescape both key and value
		key, ok1 := utils.QueryUnescape(rawQuery[start:eq], &tmpBuf)
		val, ok2 := utils.QueryUnescape(rawQuery[eq+1:end], &tmpBuf)
		if !ok1 || !ok2 {
			return nil, false // Unescape failed
		}

		// Separate hash parameter from others
		if key == "hash" {
			hash = val
		} else {
			pairs = append(pairs, utils.KV{Key: key, Val: val})
		}

		start = end + 1
	}

	// Hash parameter is mandatory
	if hash == "" {
		return nil, false
	}

	// Sort parameters lexicographically by key
	pairs.InsertionSort()

	// Build canonical string for signing
	bufPtr := utils.BufCanonicalPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]
	defer utils.BufCanonicalPool.Put(bufPtr)

	var params Params
	for i, p := range pairs {
		if i > 0 {
			buf = append(buf, '\n') // Parameters separator
		}
		buf = append(buf, p.Key...)
		buf = append(buf, '=')
		buf = append(buf, p.Val...)
		
		// Store parameter while building canonical string
		params.set(p.Key, p.Val)
	}

	// Compute HMAC-SHA256 signature
	mac := utils.GetHMACBytes(hashSecret)
	defer utils.PutHMACBytes(hashSecret, mac)
	mac.Write(buf)

	// Get buffer for computed hash from pool
	sumPtr := utils.Sha256SumBufPool.Get().(*[]byte)
	sum := (*sumPtr)[:sha256.Size]
	defer utils.Sha256SumBufPool.Put(sumPtr)
	computedHash := mac.Sum(sum[:0]) // Reuses the sum buffer

	// Decode provided hex hash
	decodedPtr := utils.Sha256SumBufPool.Get().(*[]byte)
	decodedHash := (*decodedPtr)[:sha256.Size]
	defer utils.Sha256SumBufPool.Put(decodedPtr)

	if _, err := utils.DecodeHexStringInto(hash, decodedHash); err != nil {
		return nil, false // Invalid hex encoding
	}

	// Constant-time comparison to prevent timing attacks
	if !hmac.Equal(computedHash, decodedHash) {
		return nil, false
	}

	// Return successfully parsed parameters
	return &params, true
}