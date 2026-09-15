package main

import (
	"bytes"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		buildInfoFn   func() (*debug.BuildInfo, bool)
		expectedValue string
	}{
		{
			name:    "ldflags release version takes precedence",
			version: "v1.2.3",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "v0.1.0"},
				}, true
			},
			expectedValue: "1.2.3",
		},
		{
			name:    "ldflags version without v prefix is unchanged",
			version: "1.2.3",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "v0.1.0"},
				}, true
			},
			expectedValue: "1.2.3",
		},
		{
			name:    "build info main version when ldflags is empty",
			version: "",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "v0.2.0"},
				}, true
			},
			expectedValue: "0.2.0",
		},
		{
			name:    "vcs revision fallback when main version is devel",
			version: "",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "(devel)"},
					Settings: []debug.BuildSetting{
						{Key: "vcs.revision", Value: "49923c4d090e281c1fc62723f4182fdbc4c3d0ba"},
					},
				}, true
			},
			expectedValue: "49923c4",
		},
		{
			name:    "vcs revision short length",
			version: "",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Settings: []debug.BuildSetting{
						{Key: "vcs.revision", Value: "abc"},
					},
				}, true
			},
			expectedValue: "abc",
		},
		{
			name:    "fallback to devel when no build info",
			version: "",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return nil, false
			},
			expectedValue: "devel",
		},
		{
			name:    "fallback to devel when build info has no relevant settings",
			version: "",
			buildInfoFn: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{}, true
			},
			expectedValue: "devel",
		},
		{
			name:          "fallback to devel when build info function is nil",
			version:       "",
			buildInfoFn:   nil,
			expectedValue: "devel",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := resolveVersion(tc.version, tc.buildInfoFn)
			if actual != tc.expectedValue {
				t.Errorf("expected %q, got %q", tc.expectedValue, actual)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	v := getVersion()
	if v == "" {
		t.Errorf("expected non-empty version")
	}
}

func TestPrintHelp(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	printHelp()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected help output to contain 'Usage:', got %q", output)
	}
	if !strings.Contains(output, "--version") {
		t.Errorf("expected help output to contain '--version', got %q", output)
	}
}
