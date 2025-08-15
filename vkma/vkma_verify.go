// Package vkma provides functionality for verifying VK Mini Apps launch parameters
// and validating their cryptographic signatures to ensure data integrity.
package vkma

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/elum-utils/sign/internal/utils"
)

// b64NoPad is a base64 URL encoder configured without padding characters.
// This matches VK's signature encoding format requirements.
var b64NoPad = base64.URLEncoding.WithPadding(base64.NoPadding)

// Verify validates the signature of VK Mini Apps launch parameters against
// provided application secrets using HMAC-SHA256.
//
// Parameters:
//   - rawQuery: The raw URL query string containing launch parameters
//   - secrets: A map of application IDs to their corresponding secret keys
//
// Returns:
//   - *Params: Parsed launch parameters if verification succeeds
//   - bool: true if signature is valid, false otherwise
//
// The verification process:
//   1. Parses and validates required parameters (vk_app_id and sign)
//   2. Selects the appropriate secret based on vk_app_id
//   3. Constructs the canonical parameter string
//   4. Computes HMAC-SHA256 signature
//   5. Compares with provided signature
func Verify(rawQuery string, secrets map[string]string) (*Params, bool) {
	// Early return if no secrets provided
	if len(secrets) == 0 {
		return nil, false
	}

	var appID, sign string

	// Get key-value pairs from sync.Pool to reduce allocations
	pairsPtr := utils.KVPool.Get().(*utils.KVSlice)
	pairs := (*pairsPtr)[:0] // Slice reset without reallocation
	defer utils.KVPool.Put(pairsPtr)

	// Get temporary buffer for URL unescaping from pool
	tmpBufPtr := utils.TmpBufPool.Get().(*[]byte)
	tmpBuf := *tmpBufPtr
	tmpBuf = tmpBuf[:0]
	defer utils.TmpBufPool.Put(tmpBufPtr)

	// Parse query string parameters
	for start := 0; start < len(rawQuery); {
		// Find parameter boundary
		end := strings.IndexByte(rawQuery[start:], '&')
		if end == -1 {
			end = len(rawQuery)
		} else {
			end += start
		}

		// Find key-value separator
		eq := start
		for eq < end && rawQuery[eq] != '=' {
			eq++
		}
		if eq == end {
			start = end + 1
			continue // Skip malformed parameters without values
		}

		// Unescape both key and value
		key, ok1 := utils.QueryUnescape(rawQuery[start:eq], &tmpBuf)
		val, ok2 := utils.QueryUnescape(rawQuery[eq+1:end], &tmpBuf)
		if !ok1 || !ok2 {
			return nil, false // Skip if unescaping fails
		}

		// Categorize parameters
		switch {
		case key == "sign":
			sign = val // Store signature separately
		case key == "vk_app_id":
			appID = val // Store app ID for secret lookup
			pairs = append(pairs, utils.KV{Key: key, Val: val})
		case strings.HasPrefix(key, "vk_"):
			// Include all vk_* parameters except vk_app_id already handled
			pairs = append(pairs, utils.KV{Key: key, Val: val})
		}

		start = end + 1
	}

	// Verify required parameters exist
	if appID == "" || sign == "" {
		return nil, false
	}

	// Lookup secret for this application
	secret, ok := secrets[appID]
	if !ok {
		return nil, false // Unknown app ID
	}

	// Sort parameters lexicographically by key for canonical string
	pairs.InsertionSort()

	// Get buffer for canonical string from pool
	bufPtr := utils.BufCanonicalPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]
	defer utils.BufCanonicalPool.Put(bufPtr)

	var params Params // Only allocation for result
	for i, p := range pairs {
		if i > 0 {
			buf = append(buf, '&') // Parameter separator
		}
		buf = utils.AppendEscape(buf, p.Key)
		buf = append(buf, '=')
		buf = utils.AppendEscape(buf, p.Val)
		
		// Store parameter while building canonical string
		params.set(p.Key, p.Val)
	}

	// Compute HMAC-SHA256 signature
	mac := utils.GetHMAC(secret)
	defer utils.PutHMAC(secret, mac)
	mac.Write(buf)

	// Get buffer for hash sum from pool
	sumPtr := utils.Sha256SumBufPool.Get().(*[]byte)
	sum := *sumPtr
	sum = sum[:sha256.Size] // Ensure correct size
	mac.Sum(sum[:0])        // Compute hash into buffer

	// Base64 encode the hash to compare with provided signature
	b64Ptr := utils.Base64BufPool.Get().(*[]byte)
	expectedSign := *b64Ptr
	expectedSign = expectedSign[:43] // Base64 URL encoded SHA-256 length
	b64NoPad.Encode(expectedSign, sum)
	
	// Return resources to pools
	utils.Base64BufPool.Put(b64Ptr)
	utils.Sha256SumBufPool.Put(sumPtr)

	// Constant-time comparison to prevent timing attacks
	return &params, string(expectedSign) == sign
}