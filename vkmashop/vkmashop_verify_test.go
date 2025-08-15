package vkmashop

import (
	"testing"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	secrets := map[string]string{
		"52333469": "5STCdDl55VezBzYt0AUA",
	}

	tests := []struct {
		name          string
		rawQuery      string
		clientSecrets map[string]string
		wantValid     bool
	}{
		{
			name:          "Missing secrets",
			rawQuery:      "app_id=52333469&item=Subscribtion_Item_NoAd30&lang=ru_RU",
			clientSecrets: nil,
			wantValid:     false,
		},
		{
			name:          "No signature param",
			rawQuery:      "app_id=52333469&item=Subscribtion_Item_NoAd30&lang=ru_RU",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Empty signature",
			rawQuery:      "app_id=52333469&item=Subscribtion_Item_NoAd30&lang=ru_RU&sig=",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Invalid signature",
			rawQuery:      "app_id=52333469&item=Subscribtion_Item_NoAd30&lang=ru_RU&sig=INVALIDSIG",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name: "Valid signature",
			rawQuery: "app_id=52333469" +
				"&item=Subscribtion_Item_NoAd30" +
				"&lang=ru_RU" +
				"&notification_type=get_item_test" +
				"&order_id=2256399" +
				"&receiver_id=262959639" +
				"&user_id=262959639" +
				"&sig=871447748e3803be83acb30dec37b5e5",
			clientSecrets: secrets,
			wantValid:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, ok := Verify(tt.rawQuery, tt.clientSecrets)
			if ok != tt.wantValid {
				t.Errorf("Verify() = %v, want %v", ok, tt.wantValid)
			}
			if ok && body == nil {
				t.Error("Expected non-nil *Body on valid signature")
			}
		})
	}
}

func BenchmarkVerify(b *testing.B) {
	secrets := map[string]string{
		"52333469": "5STCdDl55VezBzYt0AUA",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Verify("app_id=52333469&item=Subscribtion_Item_NoAd30&lang=ru_RU&notification_type=get_item_test&order_id=2256399&receiver_id=262959639&user_id=262959639&sig=871447748e3803be83acb30dec37b5e5", secrets)
	}
}
