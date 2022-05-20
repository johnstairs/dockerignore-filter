package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	includedAbsPath, err := filepath.Abs("foo.txt")
	require.Nil(t, err)
	excludedAbsPath, err := filepath.Abs("x.txt")
	require.Nil(t, err)
	testCases := []struct {
		input    string
		expected string
	}{
		{"x.txt", ""},
		{"foo.txt", "foo.txt"},
		{"./foo.txt", "./foo.txt"},
		{includedAbsPath, includedAbsPath},
		{excludedAbsPath, ""},
		{"bar.txt", "bar.txt"},
		{"foo.txt\nx.txt", "foo.txt"},
		{"foo.txt\nbar.txt", "foo.txt\nbar.txt"},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			ignoreFile := strings.NewReader(`
			# Ignore everything...
			**
			
			# Except...
			!foo.txt
			!/bar.txt
			`)
			output := &bytes.Buffer{}
			process(ignoreFile, strings.NewReader(tC.input), output)
			assert.Equal(t, tC.expected, strings.TrimSpace(output.String()))
		})
	}
}
