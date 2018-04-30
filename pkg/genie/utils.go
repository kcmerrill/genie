package genie

import (
	"path/filepath"
	"strings"
)

func cmd(hint string) string {
	ext := filepath.Ext(hint)
	switch ext {
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".sh":
		return "sh"
	case ".go":
		return "go run"
	default:
		return hint
	}
}

func dir(dir string) string {
	return strings.TrimRight(dir, "/") + "/"
}
