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

| Test                               | `elum-utils/sign` | `vksdk/vkapps`       | Improvement      |
|------------------------------------|-------------------|----------------------|------------------|
| Missing secrets                    | **1.39 ns/op**    | 323.7 ns/op          | 🟢 ~233x faster   |
| Invalid URL                        | **113.4 ns/op**   | 314.6 ns/op          | 🟢 ~2.8x faster   |
| No signature                       | **87.12 ns/op**   | 85.96 ns/op          | ➖ comparable     |
| Empty signature                    | **105.8 ns/op**   | 511.2 ns/op          | 🟢 ~4.8x faster   |
| No signature param                 | **107.7 ns/op**   | 126.8 ns/op          | 🟢 ~1.2x faster   |
| Malformed query                    | **150.1 ns/op**   | 279.3 ns/op          | 🟢 ~1.9x faster   |
| Valid with special chars           | **1100 ns/op**    | 2362 ns/op           | 🟢 ~2.15x faster  |
| Invalid signature                  | **1004 ns/op**    | 1993 ns/op           | 🟢 ~2x faster     |
| Valid signature                    | **1034 ns/op**    | 2401 ns/op           | 🟢 ~2.3x faster   |
| Valid signature (parallel)         | **468.9 ns/op**   | 1107 ns/op           | 🟢 ~2.36x faster  |

### 🧠 Memory Efficiency
- **50-80% fewer bytes allocated** in most cases  
- **70-80% fewer allocations** for signature validation  

---

