package config

import (
	"reflect"
	"testing"
)

func TestLoadAuthConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *AuthConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadAuthConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAuthConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAuthConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
