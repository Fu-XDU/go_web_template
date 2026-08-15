package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:generate go run pack.go

//go:embed template.zip
var templateZip []byte

var skipAnywhere = map[string]struct{}{
	".git":         {},
	".github":      {},
	".idea":        {},
	".vscode":      {},
	".env":         {},
	"node_modules": {},
	"dist":         {},
}

// Only skip these at the repo root so web/src/assets is kept.
var skipRoot = map[string]struct{}{
	"assets":  {},
	"bin":     {},
	"data":    {},
	"cmd":     {},
	"scripts": {},
}

func materializeTemplate(dest string) error {
	if len(templateZip) > 0 {
		return extractZip(templateZip, dest)
	}
	src, err := findTemplateRoot()
	if err != nil {
		return err
	}
	return copyTemplate(src, dest)
}

func extractZip(raw []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return fmt.Errorf("read embedded template: %w", err)
	}
	if err = os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	for _, f := range r.File {
		name := filepath.FromSlash(f.Name)
		if name == "" || strings.Contains(name, "..") {
			continue
		}
		target := filepath.Join(dest, name)
		if !isInside(dest, target) {
			continue
		}
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode().Perm())
		if err != nil {
			_ = rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		_ = rc.Close()
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func findTemplateRoot() (string, error) {
	if p := os.Getenv("GO_WEB_TEMPLATE_ROOT"); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", err
		}
		if !isTemplateRoot(abs) {
			return "", fmt.Errorf("GO_WEB_TEMPLATE_ROOT is not a template root: %s", abs)
		}
		return abs, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if root, ok := lookUpTemplateRoot(cwd); ok {
		return root, nil
	}

	if _, file, _, ok := runtime.Caller(0); ok && filepath.IsAbs(file) {
		root := filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
		if isTemplateRoot(root) {
			return root, nil
		}
	}

	return "", fmt.Errorf("embedded template is missing; run go generate ./cmd/create-go-web")
}

func lookUpTemplateRoot(start string) (string, bool) {
	dir := start
	for {
		if isTemplateRoot(dir) {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func isTemplateRoot(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "server.go")); err != nil {
		return false
	}
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), templateName)
}

func copyTemplate(src, dest string) error {
	src = filepath.Clean(src)
	dest = filepath.Clean(dest)
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if shouldSkip(rel, d.Name()) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		absDest := filepath.Join(dest, rel)
		if isInside(dest, path) {
			return fs.SkipDir
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() && !d.IsDir() {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(absDest, 0o755)
		}
		return copyFile(path, absDest, info.Mode())
	})
}

func shouldSkip(rel, name string) bool {
	if _, ok := skipAnywhere[name]; ok {
		return true
	}
	if _, ok := skipRoot[name]; ok && filepath.Dir(rel) == "." {
		return true
	}
	return false
}

func isInside(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func copyFile(src, dest string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

var scaffoldMarkers = [][2]string{
	{"<!-- scaffold-only:start -->", "<!-- scaffold-only:end -->"},
	{"# scaffold-only:start", "# scaffold-only:end"},
}

func stripScaffoldOnlyDocs(dest string) error {
	for _, name := range []string{"README.md", "Makefile"} {
		if err := stripScaffoldMarkers(filepath.Join(dest, name)); err != nil {
			return err
		}
	}
	return nil
}

func stripScaffoldMarkers(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	text := string(data)
	for _, pair := range scaffoldMarkers {
		start, end := pair[0], pair[1]
		for {
			i := strings.Index(text, start)
			j := strings.Index(text, end)
			if i < 0 || j < 0 || j < i {
				break
			}
			j += len(end)
			text = strings.TrimSpace(text[:i]) + "\n\n" + strings.TrimSpace(text[j:]) + "\n"
		}
	}
	if text == string(data) {
		return nil
	}
	if err = os.WriteFile(path, []byte(text), 0o644); err != nil {
		return fmt.Errorf("strip scaffold docs in %s: %w", filepath.Base(path), err)
	}
	return nil
}
