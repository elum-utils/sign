# `tma` — Telegram Mini App data validation with zero allocations

`tma` is a high-performance Go library for verifying and parsing Telegram Mini App initialization data.  
It provides **secure HMAC-SHA256 validation**, optimized parsing, and **zero allocations** in all error cases.

---

## Features

- 🚀 **Zero-allocation** in all validation failure paths (0 allocs/op)
- 📦 **Only 1 allocation** on successful validation (for the parsed `Params` struct)
- 🔒 Constant-time HMAC comparison (protection against timing attacks)
- 🛠 Optimized query parsing without `net/url`
- 💨 Benchmark-proven efficiency: ~148ns/op (parallel, valid params) on Apple M4

---

## Installation

```bash
go get github.com/elum-utils/sign
```

---

## Usage Example

```go
package main

import (
	"fmt"
	"github.com/elum-utils/sign/tma"
)

func main() {
	rawQuery := "user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=ef19060b40a2277fa4debd9c6ad9b37b1e7ac1b6f467e53c66ca6d8df2c3c168"
	secret := "1111111111:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	params, ok := tma.Verify(rawQuery, secret)
	if !ok {
		fmt.Println("Invalid Telegram Mini App data")
		return
	}

	fmt.Println("Verification passed ✅")
	user, err := params.User()
	if err != nil {
		panic(err)
	}

	fmt.Printf("User ID: %d\n", user.ID)
}
```

---

## API Reference

### `Verify`

```go
func Verify(rawQuery, secret string) (*Params, bool)
```

Validates Telegram Mini App initialization data.

#### Parameters

* `rawQuery` — the URL-encoded query string received from Telegram
* `secret` — the app secret key from Bot API (without `bot` prefix)

#### Returns

* `*Params` — structured parameters if validation succeeds
* `bool` — validation result

#### Behavior

1. Rejects empty input
2. Parses and sorts query parameters
3. Builds a canonical string for signing
4. Computes HMAC-SHA256 signature
5. Compares with provided hash in **constant time**
6. Returns parsed `Params` only if valid

#### Performance

* ❌ **0 allocations** for invalid data
* ✅ **1 allocation** for successful validation (struct `Params`)

---

### `Params`

```go
type Params struct {
	UserData     string    `json:"user" msgpack:"user"`
	ChatInstance string    `json:"chat_instance" msgpack:"chat_instance"`
	ChatType     string    `json:"chat_type" msgpack:"chat_type"`
	AuthDate     time.Time `json:"auth_date" msgpack:"auth_date"`
	Hash         string    `json:"hash" msgpack:"hash"`
}
```

Represents initialization parameters received from Telegram Mini Apps.
Includes user info, chat context, and authentication metadata.

#### Fields

* **UserData** — raw JSON with user info (use `User()` to decode)
* **ChatInstance** — unique chat session identifier
* **ChatType** — type of chat (`private`, `group`, `channel`, etc.)
* **AuthDate** — authentication timestamp (parsed from Unix time)
* **Hash** — verification hash (used for integrity check)

---

### `User`

```go
type User struct {
	ID                    int    `json:"id"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	UserName              string `json:"username"`
	PhotoURL              string `json:"photo_url"`
	Language              string `json:"language_code"`
	ChatType              string `json:"chat_type"`
	ChatInstance          string `json:"chat_instance"`
	IsPremium             bool   `json:"is_premium"`
	AllowsWriteToPM       bool   `json:"allows_write_to_pm"`
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"`
}
```

Represents a Telegram user as provided by Mini Apps initialization.
Includes profile and chat-related details.

---

### `Params.User()`

```go
func (p *Params) User() (*User, error)
```

Parses the `UserData` field into a structured `User`.

#### Returns

* `*User` — structured user object
* `error` — JSON unmarshalling error if any

#### Example

```go
user, err := params.User()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Hello, %s!\n", user.FirstName)
```

---

## Benchmarks

```
goos: darwin
goarch: arm64
cpu: Apple M4
pkg: github.com/elum-utils/sign/tma
```

| Test case                | ns/op | B/op | allocs/op |
| ------------------------ | ----- | ---- | --------- |
| Missing\_secrets         | 6.50  | 0    | 0         |
| Invalid\_params          | 27.69 | 0    | 0         |
| No\_signature            | 28.08 | 0    | 0         |
| Empty\_signature         | 42.89 | 0    | 0         |
| No\_signature\_param     | 42.31 | 0    | 0         |
| Malformed\_query         | 59.85 | 0    | 0         |
| Valid\_params            | 655.4 | 96   | 1         |
| Invalid\_params#01       | 605.3 | 96   | 1         |
| Valid\_params (parallel) | 148.4 | 96   | 1         |

---