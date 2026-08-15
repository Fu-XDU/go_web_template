//go:build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var skipNames = map[string]struct{}{
	".git":         {},
	".github":      {},
	".idea":        {},
	".vscode":      {},
	".env":         {},
	"node_modules": {},
	"dist":         {},
	"assets":       {},
	"bin":          {},
	"data":         {},
	"cmd":          {},
	"scripts":      {},
}

func main() {
	root, err := filepath.Abs("../..")
	if err != nil {
		fatal(err)
	}
	out := "template.zip"
	if err = pack(root, out); err != nil {
		fatal(err)
	}
	info, err := os.Stat(out)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("wrote %s (%d bytes) from %s\n", out, info.Size(), root)
}

func pack(root, dest string) error {
	tmp := dest + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(f)
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if _, skip := skipNames[d.Name()]; skip {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)
		header.Method = zip.Deflate
		w, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, in)
		_ = in.Close()
		return err
	})
	if closeErr := zw.Close(); err == nil {
		err = closeErr
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dest)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "pack template: %v\n", err)
	os.Exit(1)
}
