package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const templateName = "go_web_template"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	name, dest, err := promptProjectName()
	if err != nil {
		return err
	}

	fmt.Printf("\n正在生成项目 %s ...\n", dest)
	if err = materializeTemplate(dest); err != nil {
		_ = os.RemoveAll(dest)
		return err
	}
	if err = renameProject(dest, name); err != nil {
		_ = os.RemoveAll(dest)
		return err
	}
	if err = stripScaffoldOnlyDocs(dest); err != nil {
		_ = os.RemoveAll(dest)
		return err
	}

	fmt.Printf("\n完成。接下来：\n\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  cp .env.example .env\n")
	fmt.Printf("  make dev\n\n")
	return nil
}

func promptProjectName() (name, dest string, err error) {
	fmt.Println("项目名规则：字母开头，只能包含字母、数字、点(.)、下划线(_)、连字符(-)")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("项目名称: ")
		if !scanner.Scan() {
			if err = scanner.Err(); err != nil {
				return "", "", err
			}
			return "", "", fmt.Errorf("未输入项目名称")
		}
		name = strings.TrimSpace(scanner.Text())
		if err = validateProjectName(name); err != nil {
			fmt.Printf("  %v\n", err)
			continue
		}
		dest, err = filepath.Abs(name)
		if err != nil {
			return "", "", err
		}
		if _, statErr := os.Stat(dest); statErr == nil {
			fmt.Printf("  目录已存在: %s\n", dest)
			continue
		}
		return name, dest, nil
	}
}
