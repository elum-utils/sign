package vkmashop

import (
	"strconv"
)

// Params represents the payment notification parameters sent by VK's payment system.
// These parameters are received via HTTP POST when a payment event occurs in your VK Mini App.
//
// VK Payments API Documentation: 
// - Overview: https://dev.vk.com/api/payments/getting-started
// - Notifications: https://dev.vk.com/api/payments/notifications/overview
//
// Security Note: All notifications are signed with an MD5 hash using your app's secret key
// to prevent forgery. Always verify the signature before processing notifications.
type Params struct {
	// Language code for localization (e.g., "ru", "en")
	Lang             string
	
	// Unique identifier of the VK application that initiated the payment
	AppID            int
	
	// VK user ID who performed the payment action
	UserID           int
	
	// Unix timestamp when the payment event occurred
	Date             int
	
	// [Deprecated] Product name - use ItemTitle instead for new implementations
	Item             string
	
	// Discount percentage applied to the item (0-100)
	ItemDiscount     int
	
	// Merchant's unique product identifier (SKU)
	ItemID           string
	
	// URL of the product image displayed during checkout
	ItemPhotoURL     string
	
	// Product price in minor currency units (e.g., cents, kopecks)
	ItemPrice        int
	
	// Localized product title displayed to user
	ItemTitle        string
	
	// Type of payment notification:
	// - "get_item": Product information request
	// - "order_status_change": Payment status update
	// - "get_subscription": Subscription information request
	// - "subscription_status_change": Subscription status update
	NotificationType string
	
	// Current payment status:
	// - "chargeable": Ready to be processed
	// - "canceled": Payment failed
	// - "refunded": Payment was returned
	Status           string
	
	// Reason for cancellation when status="canceled":
	// - "user_decision": User canceled
	// - "app_decision": App canceled via API
	// - "payment_fail": Payment failed
	// - "unknown": Other reason
	CancelReason     string
	
	// Unique identifier for recurring subscriptions
	SubscriptionID   int
	
	// Unique transaction ID for this payment in VK's system
	OrderID          int
	
	// Merchant account ID receiving the funds
	ReceiverID       int
	
	// VK API version used (e.g., "5.131")
	Version          string
	
	// MD5 signature for request validation (concatenated params + secret key)
	Sig              string
}

// set assigns a value to the appropriate struct field based on the parameter name.
// This is used internally to parse incoming notification parameters.
//
// Parameters:
//   - key: The parameter name from VK (e.g., "app_id", "user_id")
//   - value: The raw string value from the request
//
// Note: 
// - Numeric values are automatically converted from string
// - Unknown parameters are silently ignored
// - Conversion errors result in zero values
func (b *Params) set(key, value string) {
	switch key {
	case "lang":
		b.Lang = value
	case "app_id":
		b.AppID = atoi(value)
	case "user_id":
		b.UserID = atoi(value)
	case "date":
		b.Date = atoi(value)
	case "item":
		b.Item = value
	case "item_discount":
		b.ItemDiscount = atoi(value)
	case "item_id":
		b.ItemID = value
	case "item_photo_url":
		b.ItemPhotoURL = value
	case "item_price":
		b.ItemPrice = atoi(value)
	case "item_title":
		b.ItemTitle = value
	case "notification_type":
		b.NotificationType = value
	case "status":
		b.Status = value
	case "cancel_reason":
		b.CancelReason = value
	case "subscription_id":
		b.SubscriptionID = atoi(value)
	case "order_id":
		b.OrderID = atoi(value)
	case "receiver_id":
		b.ReceiverID = atoi(value)
	case "version":
		b.Version = value
	case "sig":
		b.Sig = value
	}
}

// atoi safely converts a string to integer, returning 0 if conversion fails.
// This matches VK API behavior where missing numeric parameters default to zero.
//
// Note: This intentionally ignores conversion errors as they're handled by
// the payment system's validation logic. For critical numeric fields,
// additional validation should be performed after parsing.
//
// Security Consideration: Always validate numeric ranges after conversion,
// especially for prices and discounts.
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}