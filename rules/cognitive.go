package rules

import (
	"go/ast"
	"go/token"
)

// CognitiveComplexity calculates the cognitive complexity of a function.
func CognitiveComplexity(fn ast.Node) int {
	if fn, ok := fn.(*ast.FuncDecl); ok {
		v := cognitiveVisitor{
			name: fn.Name,
		}
		ast.Walk(&v, fn)
		return v.complexity
	}

	return 0
}

type cognitiveVisitor struct {
	name            *ast.Ident
	complexity      int
	nesting         int
	elseNodes       map[ast.Node]bool
	calculatedExprs map[ast.Expr]bool
}

func (v *cognitiveVisitor) incNesting() {
	v.nesting++
}

func (v *cognitiveVisitor) decNesting() {
	v.nesting--
}

func (v *cognitiveVisitor) incComplexity() {
	v.complexity++
}

func (v *cognitiveVisitor) nestIncComplexity() {
	v.complexity += (v.nesting + 1)
}

func (v *cognitiveVisitor) markAsElseNode(n ast.Node) {
	if v.elseNodes == nil {
		v.elseNodes = make(map[ast.Node]bool)
	}

	v.elseNodes[n] = true
}

func (v *cognitiveVisitor) markedAsElseNode(n ast.Node) bool {
	if v.elseNodes == nil {
		return false
	}

	return v.elseNodes[n]
}

func (v *cognitiveVisitor) markCalculated(e ast.Expr) {
	if v.calculatedExprs == nil {
		v.calculatedExprs = make(map[ast.Expr]bool)
	}

	v.calculatedExprs[e] = true
}

func (v *cognitiveVisitor) isCalculated(e ast.Expr) bool {
	if v.calculatedExprs == nil {
		return false
	}

	return v.calculatedExprs[e]
}

// Visit implements the ast.Visitor interface.
func (v *cognitiveVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.IfStmt:
		return v.visitIfStmt(n)
	case *ast.SwitchStmt:
		return v.visitSwitchStmt(n)
	case *ast.SelectStmt:
		return v.visitSelectStmt(n)
	case *ast.ForStmt:
		return v.visitForStmt(n)
	case *ast.RangeStmt:
		return v.visitRangeStmt(n)
	case *ast.FuncLit:
		return v.visitFuncLit(n)
	case *ast.BranchStmt:
		return v.visitBranchStmt(n)
	case *ast.BinaryExpr:
		return v.visitBinaryExpr(n)
	case *ast.CallExpr:
		return v.visitCallExpr(n)
	}
	return v
}

func (v *cognitiveVisitor) visitIfStmt(n *ast.IfStmt) ast.Visitor {
	v.incIfComplexity(n)

	if n.Init != nil {
		ast.Walk(v, n.Init)
	}

	ast.Walk(v, n.Cond)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()

	if _, ok := n.Else.(*ast.BlockStmt); ok {
		v.incComplexity()

		v.incNesting()
		ast.Walk(v, n.Else)
		v.decNesting()
	} else if _, ok := n.Else.(*ast.IfStmt); ok {
		v.markAsElseNode(n.Else)
		ast.Walk(v, n.Else)
	}
	return nil
}

func (v *cognitiveVisitor) visitSwitchStmt(n *ast.SwitchStmt) ast.Visitor {
	v.nestIncComplexity()

	if n.Init != nil {
		ast.Walk(v, n.Init)
	}

	if n.Tag != nil {
		ast.Walk(v, n.Tag)
	}

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *cognitiveVisitor) visitSelectStmt(n *ast.SelectStmt) ast.Visitor {
	v.nestIncComplexity()

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *cognitiveVisitor) visitForStmt(n *ast.ForStmt) ast.Visitor {
	v.nestIncComplexity()

	if n.Init != nil {
		ast.Walk(v, n.Init)
	}

	if n.Cond != nil {
		ast.Walk(v, n.Cond)
	}

	if n.Post != nil {
		ast.Walk(v, n.Post)
	}

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *cognitiveVisitor) visitRangeStmt(n *ast.RangeStmt) ast.Visitor {
	v.nestIncComplexity()

	if n.Key != nil {
		ast.Walk(v, n.Key)
	}

	if n.Value != nil {
		ast.Walk(v, n.Value)
	}

	ast.Walk(v, n.X)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *cognitiveVisitor) visitFuncLit(n *ast.FuncLit) ast.Visitor {
	ast.Walk(v, n.Type)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *cognitiveVisitor) visitBranchStmt(n *ast.BranchStmt) ast.Visitor {
	if n.Label != nil {
		v.incComplexity()
	}
	return v
}

func (v *cognitiveVisitor) visitBinaryExpr(n *ast.BinaryExpr) ast.Visitor {
	if (n.Op == token.LAND || n.Op == token.LOR) && !v.isCalculated(n) {
		ops := v.collectBinaryOps(n)

		var lastOp token.Token
		for _, op := range ops {
			if lastOp != op {
				v.incComplexity()
				lastOp = op
			}
		}
	}
	return v
}

func (v *cognitiveVisitor) visitCallExpr(n *ast.CallExpr) ast.Visitor {
	if name, ok := n.Fun.(*ast.Ident); ok {
		if name.Obj == v.name.Obj && name.Name == v.name.Name {
			v.incComplexity()
		}
	}
	return v
}

func (v *cognitiveVisitor) collectBinaryOps(exp ast.Expr) []token.Token {
	v.markCalculated(exp)
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return mergeBinaryOps(v.collectBinaryOps(exp.X), exp.Op, v.collectBinaryOps(exp.Y))
	case *ast.ParenExpr:
		return v.collectBinaryOps(exp.X)
	default:
		return []token.Token{}
	}
}

func (v *cognitiveVisitor) incIfComplexity(n *ast.IfStmt) {
	if v.markedAsElseNode(n) {
		v.incComplexity()
	} else {
		v.nestIncComplexity()
	}
}

func mergeBinaryOps(x []token.Token, op token.Token, y []token.Token) []token.Token {
	var out []token.Token
	if len(x) != 0 {
		out = append(out, x...)
	}
	out = append(out, op)
	if len(y) != 0 {
		out = append(out, y...)
	}
	return out
}
