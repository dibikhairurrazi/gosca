package cyclomatic

import (
	"go/ast"
	"strings"
)

type directives []string

func (ds directives) HasIgnore() bool {
	return ds.isPresent("ignore")
}

func (ds directives) isPresent(name string) bool {
	for _, d := range ds {
		if d == name {
			return true
		}
	}
	return false
}

func parseDirectives(doc *ast.CommentGroup) directives {
	if doc == nil {
		return directives{}
	}
	const prefix = "//go-sca:"
	var ds directives
	for _, comment := range doc.List {
		if strings.HasPrefix(comment.Text, prefix) {
			ds = append(ds, strings.TrimSpace(strings.TrimPrefix(comment.Text, prefix)))
		}
	}
	return ds
}
