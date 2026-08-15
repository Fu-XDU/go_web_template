package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
)

func renameProject(root, name string) error {
	replacements := [][2][]byte{
		{[]byte(publishModule), []byte(name)},
		{[]byte(templateName), []byte(name)},
	}
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			if shouldSkip(rel, d.Name()) {
				return fs.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if bytes.IndexByte(data, 0) >= 0 {
			return nil
		}
		updated := data
		changed := false
		for _, pair := range replacements {
			if bytes.Contains(updated, pair[0]) {
				updated = bytes.ReplaceAll(updated, pair[0], pair[1])
				changed = true
			}
		}
		if !changed {
			return nil
		}
		return os.WriteFile(path, updated, 0o644)
	})
}
