package utils

import "testing"

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name    string
		cpf     string
		wantErr bool
	}{
		{"Valid CPF with formatting", "123.456.789-09", false},
		{"Valid CPF without formatting", "12345678909", false},
		{"Invalid CPF - wrong check digit", "123.456.789-00", true},
		{"Invalid CPF - all zeros", "000.000.000-00", true},
		{"Invalid CPF - all same digit", "111.111.111-11", true},
		{"Invalid CPF - too short", "123.456.789", true},
		{"Invalid CPF - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCPF(tt.cpf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCPF(%q) error = %v, wantErr %v", tt.cpf, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCNPJ(t *testing.T) {
	tests := []struct {
		name    string
		cnpj    string
		wantErr bool
	}{
		{"Valid CNPJ with formatting", "11.222.333/0001-81", false},
		{"Valid CNPJ without formatting", "11222333000181", false},
		{"Invalid CNPJ - wrong check digit", "11.222.333/0001-00", true},
		{"Invalid CNPJ - all zeros", "00.000.000/0000-00", true},
		{"Invalid CNPJ - too short", "11.222.333", true},
		{"Invalid CNPJ - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCNPJ(tt.cnpj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCNPJ(%q) error = %v, wantErr %v", tt.cnpj, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCRECI(t *testing.T) {
	tests := []struct {
		name    string
		creci   string
		wantErr bool
	}{
		{"Valid CRECI with F", "12345-F/SP", false},
		{"Valid CRECI with J", "12345-J/SP", false},
		{"Valid CRECI without letter", "12345/SP", false},
		{"Invalid CRECI - wrong format", "CRECI-12345", true},
		{"Invalid CRECI - no state", "12345-F", true},
		{"Invalid CRECI - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCRECI(tt.creci)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCRECI(%q) error = %v, wantErr %v", tt.creci, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid email", "user@example.com", false},
		{"Valid email with subdomain", "user@mail.example.com", false},
		{"Valid email with plus", "user+tag@example.com", false},
		{"Invalid email - no @", "userexample.com", true},
		{"Invalid email - no domain", "user@", true},
		{"Invalid email - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneBR(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"Valid mobile with formatting", "(11) 98765-4321", false},
		{"Valid mobile without formatting", "11987654321", false},
		{"Valid landline with formatting", "(11) 3456-7890", false},
		{"Valid landline without formatting", "1134567890", false},
		{"Invalid phone - too short", "119876543", true},
		{"Invalid phone - wrong DDD", "00987654321", true},
		{"Invalid phone - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneBR(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneBR(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeCPF(t *testing.T) {
	tests := []struct {
		name string
		cpf  string
		want string
	}{
		{"CPF with formatting", "123.456.789-09", "123.456.789-09"},
		{"CPF without formatting", "12345678909", "123.456.789-09"},
		{"CPF with spaces", "123 456 789 09", "123.456.789-09"},
		{"Empty CPF", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeCPF(tt.cpf)
			if got != tt.want {
				t.Errorf("NormalizeCPF(%q) = %q, want %q", tt.cpf, got, tt.want)
			}
		})
	}
}

func TestNormalizeCNPJ(t *testing.T) {
	tests := []struct {
		name string
		cnpj string
		want string
	}{
		{"CNPJ with formatting", "11.222.333/0001-81", "11.222.333/0001-81"},
		{"CNPJ without formatting", "11222333000181", "11.222.333/0001-81"},
		{"Empty CNPJ", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeCNPJ(tt.cnpj)
			if got != tt.want {
				t.Errorf("NormalizeCNPJ(%q) = %q, want %q", tt.cnpj, got, tt.want)
			}
		})
	}
}

func TestNormalizePhoneBR(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{"Phone with formatting", "(11) 98765-4321", "(11) 98765-4321"},
		{"Phone without formatting", "11987654321", "(11) 98765-4321"},
		{"Phone with spaces", "11 9 8765 4321", "(11) 98765-4321"},
		{"Empty phone", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizePhoneBR(tt.phone)
			if got != tt.want {
				t.Errorf("NormalizePhoneBR(%q) = %q, want %q", tt.phone, got, tt.want)
			}
		})
	}
}
