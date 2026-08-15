package main

import (
	"fmt"
	"regexp"
	"strings"
)

var projectNameRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`)

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("项目名称不能为空")
	}
	if name == templateName {
		return fmt.Errorf("项目名称不能是 %q", templateName)
	}
	if strings.Contains(name, "/") || strings.Contains(name, `\`) {
		return fmt.Errorf("项目名称不能包含路径分隔符")
	}
	if !projectNameRE.MatchString(name) {
		return fmt.Errorf("项目名称不合法：必须以字母开头，且只能包含字母、数字、.、_、-")
	}
	return nil
}
