package tma

import (
	"testing"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	secrets := "1111111111:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	tests := []struct {
		name          string
		params        string
		clientSecrets string
		wantValid     bool
	}{
		{
			name:          "Valid params",
			params:        `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=ef19060b40a2277fa4debd9c6ad9b37b1e7ac1b6f467e53c66ca6d8df2c3c168`,
			clientSecrets: secrets,
			wantValid:     true,
		},
		{
			name:          "Missing secrets",
			params:        `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=ef19060b40a2277fa4debd9c6ad9b37b1e7ac1b6f467e53c66ca6d8df2c3c168`,
			clientSecrets: "",
			wantValid:     false,
		},
		{
			name:          "Invalid hash",
			params:        `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=invalid_hash`,
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Missing hash",
			params:        `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745`,
			clientSecrets: secrets,
			wantValid:     false,
		},
		{
			name:          "Malformed params",
			params:        `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=ef19060b40a2277fa4debd9c6ad9b37b1e7ac1b6f467e53c66ca6d8df2c3c168&%gh&%ij`,
			clientSecrets: secrets,
			wantValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, isValid := Validate(tt.params, tt.clientSecrets)
			if isValid != tt.wantValid {
				t.Errorf("Validate() validity = %v, want %v", isValid, tt.wantValid)
			}

			if isValid && p == nil {
				t.Error("Expected non-nil *Params on valid signature")
			}
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	secrets := "1111111111:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	validParams := `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=ef19060b40a2277fa4debd9c6ad9b37b1e7ac1b6f467e53c66ca6d8df2c3c168`
	invalidParams := `user=%7B%22id%22%3A1093776793%2C%22first_name%22%3A%22%D0%90%D1%80%D1%82%D1%83%D1%80%22%2C%22last_name%22%3A%22%D0%A4%D1%80%D0%B0%D0%BD%D0%BA%22%2C%22username%22%3A%22gmelum%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=3411281046910109270&chat_type=private&auth_date=1710181745&hash=invalid_hash`

	benchmarks := []struct {
		name     string
		params   string
		secrets  string
		parallel bool
	}{
		{
			name:    "Valid params",
			params:  validParams,
			secrets: secrets,
		},
		{
			name:    "Invalid params",
			params:  invalidParams,
			secrets: secrets,
		},
		{
			name:    "Missing secrets",
			params:  validParams,
			secrets: "",
		},
		{
			name:     "Valid params (parallel)",
			params:   validParams,
			secrets:  secrets,
			parallel: true,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.parallel {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, _ = Validate(bm.params, bm.secrets)
					}
				})
			} else {
				for i := 0; i < b.N; i++ {
					_, _ = Validate(bm.params, bm.secrets)
				}
			}
		})
	}
}
