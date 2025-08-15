package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"sync"
)

var (
	// keyedHMAC stores a sync.Pool of HMAC objects per secret key.
	keyedHMAC = sync.Map{}
)

// putHMAC returns an HMAC hash.Hash back into the pool for the given secret.
func PutHMAC(secret string, h hash.Hash) {
	h.Reset()
	if poolVal, ok := keyedHMAC.Load(secret); ok {
		poolVal.(*sync.Pool).Put(h)
	}
}

// getHMAC retrieves an HMAC hash.Hash from a pool for the given secret.
// If a pool for the secret does not yet exist, it creates one.
func GetHMAC(secret string) hash.Hash {
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

func GetHMACBytes(secretKey []byte) hash.Hash {
    poolVal, ok := keyedHMAC.Load(string(secretKey))
    if !ok {
        // копируем ключ, чтобы он жил в памяти
        keyCopy := append([]byte(nil), secretKey...)
        newPool := &sync.Pool{
            New: func() any {
                return hmac.New(sha256.New, keyCopy)
            },
        }
        keyedHMAC.Store(string(secretKey), newPool)
        return newPool.Get().(hash.Hash)
    }
    return poolVal.(*sync.Pool).Get().(hash.Hash)
}

func PutHMACBytes(secretKey []byte, h hash.Hash) {
    h.Reset()
    if poolVal, ok := keyedHMAC.Load(string(secretKey)); ok {
        poolVal.(*sync.Pool).Put(h)
    }
}
