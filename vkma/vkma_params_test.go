package vkma

import (
	"testing"
)

func TestParams_set(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected func(*Params) bool
	}{
		{
			name:  "vk_user_id sets int",
			key:   "vk_user_id",
			value: "42",
			expected: func(p *Params) bool {
				return p.VkUserID == 42
			},
		},
		{
			name:  "vk_app_id sets int",
			key:   "vk_app_id",
			value: "123",
			expected: func(p *Params) bool {
				return p.VkAppID == 123
			},
		},
		{
			name:  "vk_is_app_user sets true",
			key:   "vk_is_app_user",
			value: "1",
			expected: func(p *Params) bool {
				return p.VkIsAppUser
			},
		},
		{
			name:  "vk_are_notifications_enabled sets true",
			key:   "vk_are_notifications_enabled",
			value: "1",
			expected: func(p *Params) bool {
				return p.VkAreNotificationsEnabled
			},
		},
		{
			name:  "vk_is_favorite sets true",
			key:   "vk_is_favorite",
			value: "1",
			expected: func(p *Params) bool {
				return p.VkIsFavorite
			},
		},
		{
			name:  "vk_language sets string",
			key:   "vk_language",
			value: "en",
			expected: func(p *Params) bool {
				return p.VkLanguage == "en"
			},
		},
		{
			name:  "vk_ref sets Referral",
			key:   "vk_ref",
			value: "catalog_shopping",
			expected: func(p *Params) bool {
				return p.VkRef == CatalogShopping
			},
		},
		{
			name:  "vk_access_token_settings sets string",
			key:   "vk_access_token_settings",
			value: "friends,photos",
			expected: func(p *Params) bool {
				return p.VkAccessTokenSettings == "friends,photos"
			},
		},
		{
			name:  "vk_group_id sets int",
			key:   "vk_group_id",
			value: "777",
			expected: func(p *Params) bool {
				return p.VkGroupID == 777
			},
		},
		{
			name:  "vk_viewer_group_role sets Role",
			key:   "vk_viewer_group_role",
			value: "admin",
			expected: func(p *Params) bool {
				return p.VkViewerGroupRole == RoleAdmin
			},
		},
		{
			name:  "vk_platform sets Platform",
			key:   "vk_platform",
			value: "mobile_android",
			expected: func(p *Params) bool {
				return p.VkPlatform == MobileAndroid
			},
		},
		{
			name:  "vk_ts sets string",
			key:   "vk_ts",
			value: "1690000000",
			expected: func(p *Params) bool {
				return p.VkTs == "1690000000"
			},
		},
		{
			name:  "vk_client sets Client",
			key:   "vk_client",
			value: "ok",
			expected: func(p *Params) bool {
				return p.VkClient == ClientOk
			},
		},
		{
			name:  "sign sets string",
			key:   "sign",
			value: "abc123",
			expected: func(p *Params) bool {
				return p.Sign == "abc123"
			},
		},
		{
			name:  "invalid int does not panic",
			key:   "vk_user_id",
			value: "notanint",
			expected: func(p *Params) bool {
				return p.VkUserID == 0 // default value
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Params
			p.set(tt.key, tt.value)

			if !tt.expected(&p) {
				t.Errorf("set(%s, %s) failed: %+v", tt.key, tt.value, p)
			}
		})
	}
}

func BenchmarkParams_set(b *testing.B) {
	inputs := map[string]string{
		"vk_user_id":                  "123456",
		"vk_app_id":                   "6736218",
		"vk_is_app_user":              "1",
		"vk_are_notifications_enabled":"1",
		"vk_is_favorite":              "1",
		"vk_language":                 "ru",
		"vk_ref":                      "catalog_events",
		"vk_access_token_settings":    "friends",
		"vk_group_id":                 "9999",
		"vk_viewer_group_role":        "moder",
		"vk_platform":                 "mobile_android",
		"vk_ts":                       "1691000000",
		"vk_client":                   "ok",
		"sign":                        "sigvalue123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Params
		for k, v := range inputs {
			p.set(k, v)
		}
	}
}

func BenchmarkParams_set_vk_user_id(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p Params
		p.set("vk_user_id", "123456")
	}
}
