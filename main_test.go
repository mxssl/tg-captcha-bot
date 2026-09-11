package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	err := readConfig()
	if err != nil {
		t.Errorf("Cannot read config file. Error: %v", err)
	}
}

func TestReadConfigDurations(t *testing.T) {
	original := config
	t.Cleanup(func() { config = original })

	tests := []struct {
		name, timeout, ban, invalidField string
	}{
		{"normal", "30", "10", ""},
		{"minimum", "1", "1", ""},
		{"maximum", "9223372036", "527040", ""},
		{"forever", "30", "forever", ""},
		{"timeout typo", "3o", "10", "welcome_timeout"},
		{"timeout units", "30s", "10", "welcome_timeout"},
		{"timeout missing", "", "10", "welcome_timeout"},
		{"timeout zero", "0", "10", "welcome_timeout"},
		{"timeout negative", "-1", "10", "welcome_timeout"},
		{"timeout fraction", "1.5", "10", "welcome_timeout"},
		{"timeout duration overflow", "9223372037", "10", "welcome_timeout"},
		{"timeout integer overflow", "9223372036854775808", "10", "welcome_timeout"},
		{"ban typo", "30", "1O", "ban_duration"},
		{"ban units", "30", "10m", "ban_duration"},
		{"ban missing", "30", "", "ban_duration"},
		{"ban zero", "30", "0", "ban_duration"},
		{"ban negative", "30", "-1", "ban_duration"},
		{"ban fraction", "30", "1.5", "ban_duration"},
		{"ban exceeds Telegram limit", "30", "527041", "ban_duration"},
		{"ban duration overflow", "30", "153722868", "ban_duration"},
		{"ban integer overflow", "30", "9223372036854775808", "ban_duration"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv(configPath, dir)
			var contents string
			if tt.timeout != "" {
				contents += fmt.Sprintf("welcome_timeout = %q\n", tt.timeout)
			}
			if tt.ban != "" {
				contents += fmt.Sprintf("ban_duration = %q\n", tt.ban)
			}
			if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(contents), 0600); err != nil {
				t.Fatal(err)
			}
			before := config
			err := readConfig()
			if tt.invalidField != "" {
				if err == nil || !strings.Contains(err.Error(), tt.invalidField) {
					t.Fatalf("expected error naming %s, got %v", tt.invalidField, err)
				}
				if config != before {
					t.Fatal("invalid config replaced the active config")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			wantTimeout, err := time.ParseDuration(tt.timeout + "s")
			if err != nil {
				t.Fatal(err)
			}
			if config.welcomeTimeout != wantTimeout {
				t.Fatalf("timeout = %v, want %v", config.welcomeTimeout, wantTimeout)
			}
			var wantBan time.Duration
			if tt.ban == "forever" {
				wantBan = 367 * 24 * time.Hour
			} else {
				wantBan, err = time.ParseDuration(tt.ban + "m")
				if err != nil {
					t.Fatal(err)
				}
			}
			earliest := time.Now().Add(wantBan).Unix()
			until := getBanDuration()
			latest := time.Now().Add(wantBan).Unix()
			if until < earliest || until > latest {
				t.Fatalf("ban timestamp = %d, want between %d and %d", until, earliest, latest)
			}
		})
	}
}

func TestCorrectToken(t *testing.T) {
	token := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	if err := os.Setenv("TEST_TOKEN", token); err != nil {
		t.Fatalf("Failed to set env variable: %v", err)
	}

	v, err := getToken("TEST_TOKEN")
	if err != nil {
		t.Errorf("Incorrect token. Error: %v", err)
	}

	if v != token {
		t.Errorf("Incorrect token. Expected: %v, Have: %v", token, v)
	}
}

func TestIncorrectToken(t *testing.T) {
	token := "a123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	if err := os.Setenv("TEST_TOKEN", token); err != nil {
		t.Fatalf("Failed to set env variable: %v", err)
	}

	v, _ := getToken("TEST_TOKEN")

	if v != "" {
		t.Errorf(`Case failed. Expected "", Have: %v`, v)
	}
}
