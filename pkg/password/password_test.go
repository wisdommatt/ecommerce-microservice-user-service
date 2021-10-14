package password

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		ctx         context.Context
		passwordStr string
		cost        int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty password",
			args: args{
				passwordStr: "",
				cost:        bcrypt.DefaultCost,
			},
			wantErr: true,
		},
		{
			name: "non-empty password",
			args: args{
				passwordStr: "hello",
				cost:        bcrypt.DefaultCost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(context.Background(), tt.args.passwordStr, tt.args.cost)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" && tt.wantErr == false {
				t.Errorf("HashPassword() should not return an empty string '%s'", got)
			}
		})
	}
}
