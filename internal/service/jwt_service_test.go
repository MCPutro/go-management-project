package service

import (
	"github.com/MCPutro/go-management-project/internal/config"
	"reflect"
	"testing"
	"time"
)

func TestNewJwtService(t *testing.T) {
	type args struct {
		config *config.JwtConfig
	}
	tests := []struct {
		name string
		args args
		want JWTService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJwtService(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJwtService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jwtService_GenerateToken(t *testing.T) {
	type fields struct {
		secretKey []byte
		expiresIn time.Time
	}
	type args struct {
		userID int64
		email  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "coba1",
			fields: struct {
				secretKey []byte
				expiresIn time.Time
			}{secretKey: []byte("jwt_key"), expiresIn: time.Now().Add(1 * time.Hour)},
			args: struct {
				userID int64
				email  string
			}{userID: int64(1), email: "emnail@email.com"},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NTgwNDY0NDMsIm5iZiI6MTc1ODA0Mjg0MywiaWF0IjoxNzU4MDQyODQzfQ.kUWt13BKdF3izxc5riCc2icqih_fwoBqY9JMZdjMMt0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jwtService{
				secretKey: tt.fields.secretKey,
				expiresIn: tt.fields.expiresIn,
			}
			got, err := j.GenerateToken(tt.args.userID, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jwtService_ValidateToken(t *testing.T) {
	type fields struct {
		secretKey []byte
		expiresIn time.Time
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jwtService{
				secretKey: tt.fields.secretKey,
				expiresIn: tt.fields.expiresIn,
			}
			got, err := j.ValidateToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
