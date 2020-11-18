package rules

import (
	"go/ast"
	"go/token"
)

// AnalyzeASTFile calculates the complexities of the functions
// and methods in the abstract syntax tree (AST) of a parsed Go file and
// appends the results to the given Stats slice.
func AnalyzeASTFile(f *ast.File, fs *token.FileSet, s Stats) Stats {
	analyzer := &fileAnalyzer{
		file:    f,
		fileSet: fs,
		stats:   s,
	}
	return analyzer.analyze()
}

type fileAnalyzer struct {
	file    *ast.File
	fileSet *token.FileSet
	stats   Stats
}

func (a *fileAnalyzer) analyze() Stats {
	for _, decl := range a.file.Decls {
		a.analyzeDecl(decl)
	}
	return a.stats
}

func (a *fileAnalyzer) analyzeDecl(d ast.Decl) {
	switch decl := d.(type) {
	case *ast.FuncDecl:
		a.addStatIfNotIgnored(decl, funcName(decl), decl.Doc)
	case *ast.GenDecl:
		for _, spec := range decl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, value := range valueSpec.Values {
				funcLit, ok := value.(*ast.FuncLit)
				if !ok {
					continue
				}
				a.addStatIfNotIgnored(funcLit, valueSpec.Names[0].Name, decl.Doc)
			}
		}
	}
}

func (a *fileAnalyzer) addStatIfNotIgnored(node ast.Node, funcName string, doc *ast.CommentGroup) {
	if parseDirectives(doc).HasIgnore() {
		return
	}
	a.stats = append(a.stats, Stat{
		PkgName:    a.file.Name.Name,
		FuncName:   funcName,
		Cyclomatic: CyclomaticComplexity(node),
		Cognitive:  CognitiveComplexity(node),
		Pos:        a.fileSet.Position(node.Pos()),
	})
}
