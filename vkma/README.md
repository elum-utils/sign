# `github.com/elum-utils/sign`

A **high-performance** and **feature-rich** signature verification module for `VK Mini Apps` parameters.

> 🚀 Faster and more memory-efficient than [`github.com/SevereCloud/vksdk/v3/vkapps`](https://github.com/SevereCloud/vksdk)!

---

## 📦 Installation

```bash
go get github.com/elum-utils/sign
```

## 🧩 Usage

```go
import "github.com/elum-utils/sign/vkma"

secrets := map[string]string{
    "6736218": "wvl68m4dR1UpLrVRli",
}

url := "https://example.com/?vk_user_id=494075&vk_app_id=6736218&sign=..."

params, ok := vkma.ParamsVerify(url, secrets)
if !ok {
    // invalid signature
    return
}

// access params.VkUserID, params.VkAppID, etc.
```

## 🚦 Highlights

- ✅ Supports raw query strings, relative and absolute URLs
- 🔍 Handles edge cases: missing `sign`, malformed query, absent parameters
- 🔐 HMAC‑SHA256 signature validation with URL‑safe base64
- 💪 Zero external dependencies

## 🔬 Performance

### 📉 Benchmark Comparison (`Apple M4`, `goarch: arm64`)

| Test                       | `elum-utils/sign` | `vksdk/vkapps` | Improvement       |
|----------------------------|-------------------|----------------|-------------------|
| Missing secrets            | **1.554 ns/op**   | 323.7 ns/op    | 🟢 ~208x faster   |
| Invalid URL                | **19.63 ns/op**   | 314.6 ns/op    | 🟢 ~16x faster    |
| No signature               | **18.96 ns/op**   | 85.96 ns/op    | 🟢 ~4.5x faster   |
| Empty signature            | **34.66 ns/op**   | 511.2 ns/op    | 🟢 ~14.7x faster  |
| No signature param         | **34.57 ns/op**   | 126.8 ns/op    | 🟢 ~3.7x faster   |
| Malformed query            | **49.64 ns/op**   | 279.3 ns/op    | 🟢 ~5.6x faster   |
| Valid with special chars   | **637.8 ns/op**   | 2362 ns/op     | 🟢 ~3.7x faster   |
| Invalid signature          | **583.5 ns/op**   | 1993 ns/op     | 🟢 ~3.4x faster   |
| Valid signature            | **645.7 ns/op**   | 2401 ns/op     | 🟢 ~3.7x faster   |
| Valid signature (parallel) | **141.9 ns/op**   | 1107 ns/op     | 🟢 ~7.8x faster   |

**Key improvements:**
- All operations show significantly better performance
- 0 memory allocations for most test cases
- Only 1 alloc/op for signature operations
- Parallel processing shows ~5.4x speedup

---

