// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"strings"
	"testing"
)

func TestGuestAdditionsConfigPrepare(t *testing.T) {
	tests := []struct {
		name          string
		config        GuestAdditionsConfig
		communicator  string
		wantErrSubstr []string
		wantMode      string
		wantPath      string
		wantSHA256    string
	}{
		{
			name:         "valid disable with communicator none",
			config:       GuestAdditionsConfig{GuestAdditionsMode: "disable"},
			communicator: "none",
		},
		{
			name:         "valid attach",
			config:       GuestAdditionsConfig{GuestAdditionsMode: "attach"},
			communicator: "ssh",
		},
		{
			name:         "valid upload, default communicator",
			config:       GuestAdditionsConfig{GuestAdditionsMode: "upload"},
			communicator: "ssh",
		},
		{
			name:         "default mode upload",
			config:       GuestAdditionsConfig{},
			communicator: "ssh",
			wantMode:     "upload",
			wantPath:     "utm-guest-tools.iso",
		},
		{
			name:          "invalid mode",
			config:        GuestAdditionsConfig{GuestAdditionsMode: "foobar"},
			communicator:  "ssh",
			wantErrSubstr: []string{"guest_additions_mode is invalid"},
		},
		{
			name:          "upload with communicator none is error",
			config:        GuestAdditionsConfig{GuestAdditionsMode: "upload"},
			communicator:  "none",
			wantErrSubstr: []string{"communicator must not be 'none'"},
		},
		{
			name:         "SHA256 upper -> lower",
			config:       GuestAdditionsConfig{GuestAdditionsSHA256: "ABCDEF0123456789"},
			communicator: "ssh",
			wantSHA256:   "abcdef0123456789",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := test.config
			errs := cfg.Prepare(test.communicator)
			var errStrs []string
			for _, err := range errs {
				errStrs = append(errStrs, err.Error())
			}

			for _, want := range test.wantErrSubstr {
				found := false
				for _, got := range errStrs {
					if strings.Contains(got, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got errors: %v", want, errStrs)
				}
			}

			if len(test.wantErrSubstr) == 0 && len(errs) != 0 {
				t.Fatalf("should not have error: %v", errStrs)
			}

			if test.wantMode != "" && cfg.GuestAdditionsMode != test.wantMode {
				t.Errorf("wanted mode default %q, got %q", test.wantMode, cfg.GuestAdditionsMode)
			}
			if test.wantPath != "" && cfg.GuestAdditionsPath != test.wantPath {
				t.Errorf("wanted path default %q, got %q", test.wantPath, cfg.GuestAdditionsPath)
			}

			if test.wantSHA256 != "" && cfg.GuestAdditionsSHA256 != test.wantSHA256 {
				t.Errorf("wanted sha256 %q, got %q", test.wantSHA256, cfg.GuestAdditionsSHA256)
			}
		})
	}
}
