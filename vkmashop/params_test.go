package vkmashop

import "testing"

func TestBody_set(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected func(*Body) bool
	}{
		{
			name:  "app_id sets int",
			key:   "app_id",
			value: "123",
			expected: func(b *Body) bool {
				return b.AppID == 123
			},
		},
		{
			name:  "user_id sets int",
			key:   "user_id",
			value: "456",
			expected: func(b *Body) bool {
				return b.UserID == 456
			},
		},
		{
			name:  "item sets string",
			key:   "item",
			value: "example_item",
			expected: func(b *Body) bool {
				return b.Item == "example_item"
			},
		},
		{
			name:  "item_discount sets int",
			key:   "item_discount",
			value: "10",
			expected: func(b *Body) bool {
				return b.ItemDiscount == 10
			},
		},
		{
			name:  "order_id sets int",
			key:   "order_id",
			value: "777",
			expected: func(b *Body) bool {
				return b.OrderID == 777
			},
		},
		{
			name:  "sig sets string",
			key:   "sig",
			value: "testsig",
			expected: func(b *Body) bool {
				return b.Sig == "testsig"
			},
		},
		{
			name:  "invalid int input does not crash",
			key:   "app_id",
			value: "notanumber",
			expected: func(b *Body) bool {
				return b.AppID == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Body
			b.set(tt.key, tt.value)
			if !tt.expected(&b) {
				t.Errorf("set(%s, %s) failed: %+v", tt.key, tt.value, b)
			}
		})
	}
}

func BenchmarkBody_set(b *testing.B) {
	inputs := map[string]string{
		"app_id":         "101",
		"user_id":        "202",
		"item":           "gold_sub",
		"item_discount":  "5",
		"item_price":     "499",
		"notification_type": "chargeable",
		"order_id":       "80085",
		"sig":            "md5hashsig",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var body Body
		for k, v := range inputs {
			body.set(k, v)
		}
	}
}

