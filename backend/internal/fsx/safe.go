package fsx

import (
	"fmt"
	"path/filepath"
	"strings"
)

func SafeJoinUploadDir(base string, rel string) (string, error) {
	rel = filepath.ToSlash(rel)
	if strings.Contains(rel, "..") {
		return "", fmt.Errorf("invalid path")
	}
	full := filepath.Join(base, filepath.FromSlash(rel))
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	sep := string(filepath.Separator)
	if absFull != absBase && !strings.HasPrefix(absFull+sep, absBase+sep) {
		return "", fmt.Errorf("path escapes upload directory")
	}
	return absFull, nil
}
