package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
)

func renameProject(root, name string) error {
	old := []byte(templateName)
	neu := []byte(name)
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
		if bytes.IndexByte(data, 0) >= 0 || !bytes.Contains(data, old) {
			return nil
		}
		updated := bytes.ReplaceAll(data, old, neu)
		return os.WriteFile(path, updated, 0o644)
	})
}
