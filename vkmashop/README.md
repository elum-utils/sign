# `vkmashop` — VK Shop API request signature verification

`vkmashop` is a Go library for **verifying VK Shop API request signatures**.  
It implements **MD5-based validation** as specified in [VK Shop documentation](https://dev.vk.com/api/payments/notifications/overview).  
The package ensures requests are authentic before processing payment notifications.

---

## Features

- 🔒 Secure **MD5-based signature verification**  
- 🚀 **Zero extra allocations** except for the parsed `Params` struct  
- 📦 **Only 1 allocation** on success  
- 🛠 Optimized query parsing (no `net/url`)  
- 💨 Benchmark-proven efficiency (~488ns/op on Apple M4)  

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
	"github.com/elum-utils/sign/vkmashop"
)

func main() {
	rawQuery := "app_id=52333469&item=Subscribtion_Item_NoAd30&lang=ru_RU&notification_type=get_item_test&order_id=2256399&receiver_id=262959639&user_id=262959639&sig=871447748e3803be83acb30dec37b5e5"
	secrets := map[string]string{
		"52333469": "5STCdDl55VezBzYt0AUA",, // app_id → secret mapping
	}

	params, ok := vkmashop.Verify(rawQuery, secrets)
	if !ok {
		fmt.Println("Invalid VK Shop request ❌")
		return
	}

	fmt.Println("Signature verified ✅")
	fmt.Printf("Order ID: %d, Status: %s\n", params.OrderID, params.Status)
}
```

---

## API Reference

### `Verify`

```go
func Verify(rawQuery string, secrets map[string]string) (*Params, bool)
```

Validates the **MD5 signature** of VK Shop API requests.

#### Parameters

* `rawQuery` — raw query string with request parameters
* `secrets` — mapping of `app_id` to secret keys

#### Returns

* `*Params` — parsed request parameters if valid
* `bool` — `true` if the signature is valid, `false` otherwise

#### Verification process

1. Parse and validate required parameters (`app_id`, `sig`)
2. Select secret key for `app_id`
3. Build signature string (`key=value` + secret)
4. Compute **MD5 hash**
5. Compare with provided signature (hex) without allocations

#### Performance

* ❌ **0 allocations** for invalid requests
* ✅ **1 allocation** for successful validation (struct `Params`)

---

### `Params`

```go
type Params struct {
	Lang             string
	AppID            int
	UserID           int
	Date             int
	Item             string
	ItemDiscount     int
	ItemID           string
	ItemPhotoURL     string
	ItemPrice        int
	ItemTitle        string
	NotificationType NotificationType
	Status           Status
	CancelReason     CancelReason
	SubscriptionID   int
	OrderID          int
	ReceiverID       int
	Version          string
	Sig              string
}
```

Represents **payment notification parameters** sent by VK’s payment system.
Notifications are received via HTTP POST when a payment/subscription event occurs.

#### Key fields

* **AppID** — VK app identifier
* **UserID** — VK user who made the payment
* **Date** — Unix timestamp of the event
* **ItemID / ItemTitle** — product identifier and title
* **ItemPrice** — product price (in minor units, e.g. cents)
* **NotificationType** — event type (`get_item`, `order_status_change`, etc.)
* **Status** — payment status (`chargeable`, `canceled`, `refunded`, `active`)
* **CancelReason** — reason for cancellation (`user_decision`, `payment_fail`, etc.)
* **SubscriptionID** — recurring subscription identifier
* **OrderID** — transaction identifier in VK’s system
* **Sig** — MD5 signature (must be verified)

---

### `Status`

```go
type Status string

const (
	Chargeable Status = "chargeable"
	Canceled   Status = "canceled"
	Refunded   Status = "refunded"
	Active     Status = "active"
)
```

Represents the current **payment status**.

---

### `NotificationType`

```go
type NotificationType string

const (
	GetItem                      NotificationType = "get_item"
	GetItemTest                  NotificationType = "get_item_test"
	OrderStatusChange            NotificationType = "order_status_change"
	OrderStatusChangeTest        NotificationType = "order_status_change_test"
	GetSubscription              NotificationType = "get_subscription"
	SubscriptionStatusChange     NotificationType = "subscription_status_change"
	SubscriptionStatusChangeTest NotificationType = "subscription_status_change_test"
)
```

Represents the type of **payment/subscription notification**.

---

### `CancelReason`

```go
type CancelReason string

const (
	CancelUserDecision CancelReason = "user_decision"
	CancelAppDecision  CancelReason = "app_decision"
	CancelPaymentFail  CancelReason = "payment_fail"
	CancelUnknown      CancelReason = "unknown"
)
```

Represents the **reason for payment cancellation**.

---

## Benchmarks

```
goos: darwin
goarch: arm64
cpu: Apple M4
pkg: github.com/elum-utils/sign/vkmashop
```

| Test case      | ns/op | B/op | allocs/op |
| -------------- | ----- | ---- | --------- |
| Verify request | 488.1 | 224  | 1         |

---