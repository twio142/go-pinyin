package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mozillazg/go-pinyin"
)

// Helper function to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestProcessLine(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		line          string
		isInitialMode bool
		isXiaoheMode  bool
		isOnlyMode    bool
		expected      string
	}{
		{
			name:          "Normal Pinyin",
			line:          "你好",
			isInitialMode: false,
			isXiaoheMode:  false,
			isOnlyMode:    false,
			expected:      "你好	ni hao",
		},
		{
			name:          "Initials Mode",
			line:          "你好",
			isInitialMode: true,
			isXiaoheMode:  false,
			isOnlyMode:    false,
			expected:      "你好	n h",
		},
		{
			name:          "Xiaohe Shuangpin Mode",
			line:          "你好",
			isInitialMode: false,
			isXiaoheMode:  true,
			isOnlyMode:    false,
			expected:      "你好	ni hc",
		},
		{
			name:          "Xiaohe Shuangpin & Initials Mode",
			line:          "你好世界",
			isInitialMode: true,
			isXiaoheMode:  true,
			isOnlyMode:    false,
			expected:      "你好世界	n h u j",
		},
		{
			name:          "Only Mode",
			line:          "你好",
			isInitialMode: false,
			isXiaoheMode:  false,
			isOnlyMode:    true,
			expected:      "ni hao",
		},
		{
			name:          "Full-width to Half-width",
			line:          "你好，world！",
			isInitialMode: false,
			isXiaoheMode:  false,
			isOnlyMode:    false,
			expected:      "你好，world！	ni hao , world!",
		},
		{
			name:          "Mixed Text",
			line:          "Hello你好World",
			isInitialMode: false,
			isXiaoheMode:  false,
			isOnlyMode:    false,
			expected:      "Hello你好World	Hello ni hao World",
		},
		{
			name:          "Empty Input",
			line:          "",
			isInitialMode: false,
			isXiaoheMode:  false,
			isOnlyMode:    false,
			expected:      "",
		},
		{
			name:          "No Chinese Characters",
			line:          "Hello World",
			isInitialMode: false,
			isXiaoheMode:  false,
			isOnlyMode:    false,
			expected:      "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new pinyin Args object for each test
			args := pinyin.NewArgs()

			// Capture the output of processLine
			output := captureOutput(func() {
				processLine(tt.line, args, tt.isInitialMode, tt.isXiaoheMode, tt.isOnlyMode)
			})

			// Trim whitespace and compare
			if strings.TrimSpace(output) != strings.TrimSpace(tt.expected) {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expected, output)
			}
		})
	}
}
