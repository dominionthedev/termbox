// Package envutil provides helpers for reading and writing termbox.env.
package envutil

import (
	"fmt"
	"os"
	"strings"
)

// UpdateVar rewrites (or appends) a KEY="VALUE" export line in the given env file.
// It performs an atomic write: temp file → rename.
func UpdateVar(envPath, key, value string) error {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("reading %q: %w — run 'termbox setup --wizard' first", envPath, err)
	}

	lines := strings.Split(string(data), "\n")
	found := false
	prefix1 := "export " + key + "="
	prefix2 := key + "="
	for i, line := range lines {
		if strings.HasPrefix(line, prefix1) || strings.HasPrefix(line, prefix2) {
			lines[i] = fmt.Sprintf("export %s=%q", key, value)
			found = true
		}
	}
	if !found {
		lines = append(lines, fmt.Sprintf("export %s=%q", key, value))
	}

	tmp := envPath + ".tmp"
	if err := os.WriteFile(tmp, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("writing temp: %w", err)
	}
	if err := os.Rename(tmp, envPath); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("replacing env file: %w", err)
	}
	return nil
}

// ReadVar returns the value of KEY from the given env file, or "" if not found.
func ReadVar(envPath, key string) string {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return ""
	}
	prefix := "export " + key + "="
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, prefix) {
			v := strings.TrimPrefix(line, prefix)
			return strings.Trim(v, `"'`)
		}
	}
	return ""
}
