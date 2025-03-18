package config

import (
	"reflect"
	"testing"

	"github.com/uptrace/bun"
)

func TestLoadDatabaseConfig(t *testing.T) {
	tests := []struct {
		name string
		want *DBConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadDatabaseConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadDatabaseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadDatabase(t *testing.T) {
	tests := []struct {
		name string
		want *bun.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadDatabase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildDSN(t *testing.T) {
	type args struct {
		config DBConfig
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDSN(tt.args.config); got != tt.want {
				t.Errorf("buildDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
