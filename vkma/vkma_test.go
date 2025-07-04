package vkma

import (
	"testing"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	secrets := map[string]string{
		"6736218": "wvl68m4dR1UpLrVRli",
	}

	tests := []struct {
		name          string
		URL           string
		clientSecrets map[string]string
		wantValid     bool
	}{
		{
			name:          "Missing secrets",
			URL:           "http://[%10::1]",
			clientSecrets: nil,
			wantValid:     false,
		},
		{
			name:          "Invalid URL",
			URL:           "http://[%10::1]",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "No signature",
			URL:           "https://example.com",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Empty signature",
			URL:           "https://example.com?sign=abc",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "No signature param",
			URL:           "https://example.com?q=abc",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Malformed query",
			URL:           "https://example.com?sign=abc&%gh&%ij",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Valid with special chars in query",
			URL:           "q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=andr%26oid&sign=gAgvKPEe3wJiC9ZdT16XuZ65_KSH5WkGSeDp_CQofws",
			clientSecrets: secrets,
			wantValid:     true,
		},
		{
			name:          "Valid with question mark",
			URL:           "?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=andr%26oid&sign=gAgvKPEe3wJiC9ZdT16XuZ65_KSH5WkGSeDp_CQofws",
			clientSecrets: secrets,
			wantValid:     true,
		},
		{
			name:          "Valid full URL",
			URL:           "https://example.com/?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=andr%26oid&sign=gAgvKPEe3wJiC9ZdT16XuZ65_KSH5WkGSeDp_CQofws",
			clientSecrets: secrets,
			wantValid:     true,
		},
		{
			name:          "Invalid signature",
			URL:           "https://example.com/?vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=notify&vk_platform=android&sign=exTIBPYTrAKDTHLLm2AwJkmcVcvFCzQUNyoa6wAjvW6k",
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Valid signature",
			URL:           "https://example.com/?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=android&sign=htQFduJpLxz7ribXRZpDFUH-XEUhC9rBPTJkjUFEkRA",
			clientSecrets: secrets,
			wantValid:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, isValid := Verify(tt.URL, tt.clientSecrets)
			if isValid != tt.wantValid {
				t.Errorf("Verify() validity = %v, want %v", isValid, tt.wantValid)
			}

			if isValid && p == nil {
				t.Error("Expected non-nil *Params on valid signature")
			}
		})
	}
}

func BenchmarkVerify(b *testing.B) {
	secrets := map[string]string{
		"6736218": "wvl68m4dR1UpLrVRli",
	}

	benchmarks := []struct {
		name     string
		URL      string
		secrets  map[string]string
		parallel bool
	}{
		{
			name:    "Missing secrets",
			URL:     "http://[%10::1]",
			secrets: nil,
		},
		{
			name:    "Invalid URL",
			URL:     "http://[%10::1]",
			secrets: secrets,
		},
		{
			name:    "No signature",
			URL:     "https://example.com",
			secrets: secrets,
		},
		{
			name:    "Empty signature",
			URL:     "https://example.com?sign=abc",
			secrets: secrets,
		},
		{
			name:    "No signature param",
			URL:     "https://example.com?q=abc",
			secrets: secrets,
		},
		{
			name:    "Malformed query",
			URL:     "https://example.com?sign=abc&%gh&%ij",
			secrets: secrets,
		},
		{
			name:    "Valid with special chars",
			URL:     "https://example.com/?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=andr%26oid&sign=gAgvKPEe3wJiC9ZdT16XuZ65_KSH5WkGSeDp_CQofws",
			secrets: secrets,
		},
		{
			name:    "Invalid signature",
			URL:     "https://example.com/?vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=notify&vk_platform=android&sign=exTIBPYTrAKDTHLLm2AwJkmcVcvFCzQUNyoa6wAjvW6k",
			secrets: secrets,
		},
		{
			name:    "Valid signature",
			URL:     "https://example.com/?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=android&sign=htQFduJpLxz7ribXRZpDFUH-XEUhC9rBPTJkjUFEkRA",
			secrets: secrets,
		},
		{
			name:     "Valid signature (parallel)",
			URL:      "https://example.com/?q=1&vk_user_id=494075&vk_app_id=6736218&vk_is_app_user=1&vk_are_notifications_enabled=1&vk_language=ru&vk_access_token_settings=&vk_platform=android&sign=htQFduJpLxz7ribXRZpDFUH-XEUhC9rBPTJkjUFEkRA",
			secrets:  secrets,
			parallel: true,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.parallel {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, _ = Verify(bm.URL, bm.secrets)
					}
				})
			} else {
				for i := 0; i < b.N; i++ {
					_, _ = Verify(bm.URL, bm.secrets)
				}
			}
		})
	}
}
