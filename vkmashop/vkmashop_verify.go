// Package vkmashop provides functionality for verifying VK Shop API request signatures.
// It implements MD5-based signature validation as specified in VK Shop documentation.
package vkmashop

import (
	"crypto/md5"
	"strings"

	"github.com/elum-utils/sign/internal/utils"
)

// Verify validates the signature of VK Shop API requests against provided application secrets.
//
// Parameters:
//   - rawQuery: The raw URL query string containing request parameters
//   - secrets: A map of application IDs to their corresponding secret keys
//
// Returns:
//   - *Params: Parsed request parameters if verification succeeds
//   - bool: true if signature is valid, false otherwise
//
// The verification process:
//   1. Parses and validates required parameters (app_id and sig)
//   2. Selects the appropriate secret based on app_id
//   3. Constructs the signature string according to VK Shop specification
//   4. Computes MD5 hash
//   5. Compares with provided signature without string allocations
func Verify(rawQuery string, secrets map[string]string) (*Params, bool) {
	// Early return if no secrets provided
	if len(secrets) == 0 {
		return nil, false
	}

	var appID, sig string

	// Get key-value pairs from sync.Pool to reduce allocations
	pairsPtr := utils.KVPool.Get().(*utils.KVSlice)
	pairs := (*pairsPtr)[:0] // Slice reset without reallocation
	defer utils.KVPool.Put(pairsPtr)

	// Get temporary buffer for URL unescaping from pool
	tmpBufPtr := utils.TmpBufPool.Get().(*[]byte)
	tmpBuf := *tmpBufPtr
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
		switch key {
		case "app_id":
			appID = val // Store app ID for secret lookup
			pairs = append(pairs, utils.KV{Key: key, Val: val})
		case "sig":
			sig = val // Store signature separately
		default:
			// Include all other parameters in verification
			pairs = append(pairs, utils.KV{Key: key, Val: val})
		}

		start = end + 1
	}

	// Verify required parameters exist
	if appID == "" || sig == "" {
		return nil, false
	}

	// Lookup secret for this application
	secret, ok := secrets[appID]
	if !ok {
		return nil, false // Unknown app ID
	}

	// Sort parameters lexicographically by key
	pairs.InsertionSort()

	// Get buffer for signature string from pool
	bufPtr := utils.BufCanonicalPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]
	defer utils.BufCanonicalPool.Put(bufPtr)

	body := &Params{} // Only allocation for result
	for _, p := range pairs {
		// Build signature string format: key=value
		buf = append(buf, p.Key...)
		buf = append(buf, '=')
		buf = append(buf, p.Val...)
		
		// Store parameter while building signature string
		body.set(p.Key, p.Val)
	}
	// Append secret key as specified in VK Shop docs
	buf = append(buf, secret...)

	// Compute MD5 hash of the signature string
	sum := md5.Sum(buf)

	// Validate signature format and compare
	if len(sig) != 32 { // MD5 hex string should be 32 chars
		return nil, false
	}

	// Compare computed hash with provided signature
	// without converting to string to avoid allocations
	for i := 0; i < 16; i++ {
		// Decode hex digits directly
		hi := utils.FromHex(sig[i*2])
		lo := utils.FromHex(sig[i*2+1])
		
		// Check for invalid hex digits (255 indicates error)
		if hi == 255 || lo == 255 {
			return nil, false
		}
		
		// Compare each byte of the hash
		if sum[i] != (hi<<4|lo) {
			return nil, false
		}
	}

	return body, true
}