package tma

import (
	"testing"
	"time"
)

func TestParams_set(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected func(*Params) bool
	}{
		{
			name:  "user sets UserData",
			key:   "user",
			value: `{"id":123,"first_name":"John"}`,
			expected: func(p *Params) bool {
				return p.UserData == `{"id":123,"first_name":"John"}`
			},
		},
		{
			name:  "chat_instance sets ChatInstance",
			key:   "chat_instance",
			value: "abc123",
			expected: func(p *Params) bool {
				return p.ChatInstance == "abc123"
			},
		},
		{
			name:  "chat_type sets ChatType",
			key:   "chat_type",
			value: "private",
			expected: func(p *Params) bool {
				return p.ChatType == "private"
			},
		},
		{
			name:  "auth_date sets valid time",
			key:   "auth_date",
			value: "2301011200",
			expected: func(p *Params) bool {
				expectedTime, _ := time.Parse("0601021504", "2301011200")
				return p.AuthDate.Equal(expectedTime)
			},
		},
		{
			name:  "auth_date with invalid format doesn't set",
			key:   "auth_date",
			value: "invalid",
			expected: func(p *Params) bool {
				return p.AuthDate.IsZero()
			},
		},
		{
			name:  "unknown key doesn't modify Params",
			key:   "unknown",
			value: "value",
			expected: func(p *Params) bool {
				return p.UserData == "" && 
					p.ChatInstance == "" && 
					p.ChatType == "" && 
					p.AuthDate.IsZero()
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

func TestParams_User(t *testing.T) {
	tests := []struct {
		name     string
		userData string
		expected func(*User, error) bool
	}{
		{
			name:     "valid user data",
			userData: `{"id":123,"first_name":"John","last_name":"Doe","username":"johndoe"}`,
			expected: func(u *User, err error) bool {
				return u != nil && 
					u.ID == 123 && 
					u.FirstName == "John" && 
					u.LastName == "Doe" && 
					u.UserName == "johndoe" && 
					err == nil
			},
		},
		{
			name:     "invalid user data",
			userData: `invalid json`,
			expected: func(u *User, err error) bool {
				return u != nil && err != nil
			},
		},
		{
			name:     "empty user data",
			userData: "",
			expected: func(u *User, err error) bool {
				return u != nil && err != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{UserData: tt.userData}
			user, err := p.User()

			if !tt.expected(user, err) {
				t.Errorf("User() failed with input %s: user=%+v, err=%v", tt.userData, user, err)
			}
		})
	}
}

func BenchmarkParams_set(b *testing.B) {
	inputs := map[string]string{
		"user":           `{"id":123,"first_name":"John"}`,
		"chat_instance":  "abc123",
		"chat_type":      "private",
		"auth_date":      "2301011200",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Params
		for k, v := range inputs {
			p.set(k, v)
		}
	}
}

func BenchmarkParams_User(b *testing.B) {
	p := Params{UserData: `{"id":123,"first_name":"John","last_name":"Doe","username":"johndoe"}`}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.User()
	}
}