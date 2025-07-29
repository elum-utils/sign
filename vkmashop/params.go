package vkmashop

import (
	"strconv"
)

type Body struct {
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
	NotificationType string
	Status           string
	CancelReason     string
	SubscriptionID   int
	OrderID          int
	ReceiverID       int
	Version          string
	Sig              string
}

func (b *Body) set(key, value string) {
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

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
