# `github.com/elum-utils/sign/vkma`

A **high-performance** and **memory-efficient** signature verification package for **VK Mini Apps** launch parameters.  
Implements **HMAC-SHA256** validation with URL-safe Base64 encoding, following VK’s official specification.

> ⚡ Optimized to be faster and lighter than [`github.com/SevereCloud/vksdk/v3/vkapps`](https://github.com/SevereCloud/vksdk)!

---

## 📦 Installation

```bash
go get github.com/elum-utils/sign
```

---

## 🧩 Usage

```go
import "github.com/elum-utils/sign/vkma"

secrets := map[string]string{
    "6736218": "wvl68m4dR1UpLrVRli",
}

url := "https://example.com/?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=android&sign=htQFduJpLxz7ribXRZpDFUH-XEUhC9rBPTJkjUFEkRA"

params, ok := vkma.Verify(url, secrets)
if !ok {
    // invalid signature
    return
}

// access strongly typed fields:
_ = params.VkUserID
_ = params.VkAppID
_ = params.VkPlatform
```

---

## 🚦 Highlights

* ✅ Works with raw query strings, relative and absolute URLs
* 🔍 Graceful handling of edge cases (missing `sign`, malformed query, absent params)
* 🔐 HMAC-SHA256 signature validation with URL-safe Base64 encoding
* 💪 Zero allocations for invalid inputs, only **1 alloc/op** for valid signatures
* 🧵 Optimized for parallel verification

---

## 🔬 Performance

### 📉 Benchmark Results (`Apple M4`, `goarch: arm64`)

| Test                       | ns/op | B/op | allocs/op |
| -------------------------- | ----- | ---- | --------- |
| Missing secrets            | 5.457 | 0    | 0         |
| Invalid URL                | 38.82 | 0    | 0         |
| No signature               | 33.77 | 0    | 0         |
| Empty signature            | 50.97 | 0    | 0         |
| No signature param         | 45.86 | 0    | 0         |
| Malformed query            | 61.73 | 0    | 0         |
| Valid with special chars   | 642.3 | 160  | 1         |
| Invalid signature          | 567.1 | 160  | 1         |
| Valid signature            | 597.6 | 160  | 1         |
| Valid signature (parallel) | 127.4 | 160  | 1         |

**Key points:**

* 🚀 Up to **200x faster** than `vksdk/vkapps` in missing-secrets case
* 🔒 0 allocations for all failure scenarios
* ⚡ Stable performance under parallel workloads

---

## 🧭 Parameters

The `Params` struct provides **strongly typed access** to all launch parameters.

```go
type Params struct {
    VkUserID                  int
    VkAppID                   int
    VkIsAppUser               bool
    VkAreNotificationsEnabled bool
    VkIsFavorite              bool
    VkLanguage                string
    VkRef                     Referral
    VkAccessTokenSettings     string
    VkGroupID                 int
    VkViewerGroupRole         Role
    VkPlatform                Platform
    VkTs                      string
    VkClient                  Client
    Sign                      string
}
```

## 📚 How it works

1. Parse query parameters
2. Collect only `vk_*` params (excluding `sign`)
3. Sort lexicographically to build canonical string
4. Compute HMAC-SHA256 with the corresponding secret
5. Encode with **Base64 (URL-safe, no padding)**
6. Compare with provided `sign` (constant-time)

---