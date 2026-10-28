package rule

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UseNewRule spots variables declared just to take their address, that can be replaced by a call to the new builtin.
type UseNewRule struct{}

// Apply applies the rule to given file.
func (*UseNewRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go126) {
		return nil // new accepts a value to point to only since Go 1.26
	}

	var failures []lint.Failure

	walker := lintUseNew{
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*UseNewRule) Name() string {
	return "use-new"
}

type lintUseNew struct {
	onFailure func(lint.Failure)
}

func (w lintUseNew) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.BlockStmt:
		w.checkPointerToTemporaryVar(n.List)
	case *ast.CaseClause:
		w.checkPointerToTemporaryVar(n.Body)
	case *ast.CommClause:
		w.checkPointerToTemporaryVar(n.Body)
	}

	return w
}

// checkPointerToTemporaryVar spots the
//
//	v := expr
//	p = &v
//
// idiom, where v is used nowhere else, that can be replaced by p = new(expr).
func (w lintUseNew) checkPointerToTemporaryVar(stmts []ast.Stmt) {
	for i := 0; i+1 < len(stmts); i++ {
		target, value := definedVar(stmts[i])
		if target == nil {
			continue
		}

		if !isSimpleValue(value) {
			continue
		}

		if !takesAddressOfVar(stmts[i+1], target.Name) {
			continue
		}

		if countUsesInStmts(stmts, target.Name) != 2 {
			continue
		}

		valueStr := astutils.GoFmt(value)

		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Node:       stmts[i],
			Confidence: 0.8,
			Failure:    fmt.Sprintf("replace the temporary variable %s with new(%s)", target.Name, valueStr),
		})

		i++ // the statement taking the address is rewritten by the failure above, reporting it in turn would contradict it
	}
}

// isSimpleValue returns true if the value reads well inside a call to new,
// for example 42, *size, status.count or float64(*size).
// A value hiding what it holds behind a call, e.g. ptr.Deref(count, 0) + 1 or kind.GroupVersionKind(),
// is better left in a variable, and so is a composite literal: &T{...} already allocates one.
func isSimpleValue(value ast.Expr) bool {
	switch value := value.(type) {
	case *ast.BasicLit, *ast.Ident:
		return true
	case *ast.ParenExpr:
		return isSimpleValue(value.X)
	case *ast.StarExpr:
		return isSimpleValue(value.X)
	case *ast.UnaryExpr:
		// a channel receive is left out: evaluated in place of the variable, it would block
		// where the call it moves into is evaluated, after the operands preceding it
		return value.Op != token.ARROW && isSimpleValue(value.X)
	case *ast.SelectorExpr:
		return isSimpleValue(value.X)
	case *ast.IndexExpr:
		return isSimpleValue(value.X) && isSimpleValue(value.Index)
	case *ast.CallExpr:
		// a conversion to a predeclared type, e.g. float64(*size).
		// The result of a function call is better left in a variable, where a name tells what it is.
		if len(value.Args) != 1 || value.Ellipsis.IsValid() {
			return false
		}

		convertedTo, ok := value.Fun.(*ast.Ident)

		return ok && builtInTypes[convertedTo.Name] && isSimpleValue(value.Args[0])
	default:
		return false
	}
}

// definedVar returns the variable and the value of a v := expr statement.
// It returns nil if the statement is not of that form.
func definedVar(stmt ast.Stmt) (*ast.Ident, ast.Expr) {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return nil, nil
	}

	target, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return nil, nil
	}

	return target, assign.Rhs[0]
}

// takesAddressOfVar returns true if the statement takes the address of the variable name and makes no other use of it,
// for example p = &name, return &name or f(&name).
func takesAddressOfVar(stmt ast.Stmt, name string) bool {
	if !isStraightLine(stmt) {
		return false
	}

	if countUses(stmt, name) != 1 {
		return false // the statement does not use the variable, or uses it more than once
	}

	// the address taken inside a function literal is not reported: the literal can run later,
	// when evaluating the value in place of the variable would not necessarily yield the same result
	inFuncLit := astutils.PickNodes(stmt, func(n ast.Node) bool {
		funcLit, ok := n.(*ast.FuncLit)
		return ok && countUses(funcLit, name) > 0
	})
	if len(inFuncLit) > 0 {
		return false
	}

	addressOf := astutils.PickNodes(stmt, func(n ast.Node) bool {
		unary, ok := n.(*ast.UnaryExpr)
		return ok && unary.Op == token.AND && astutils.IsIdent(unary.X, name)
	})

	return len(addressOf) == 1
}

// isStraightLine returns true if the statement runs exactly once, where it is written.
// The address taken in the body of a loop or of a conditional statement is not reported:
// the body runs many times or not at all, when evaluating the value in place of the variable
// would allocate one value per run, or would not evaluate it at all.
func isStraightLine(stmt ast.Stmt) bool {
	switch stmt.(type) {
	case *ast.AssignStmt, *ast.DeclStmt, *ast.DeferStmt, *ast.ExprStmt, *ast.GoStmt, *ast.ReturnStmt, *ast.SendStmt:
		return true
	default:
		return false
	}
}

// countUses returns the number of occurrences of the variable name in the AST subtree rooted at n.
// A field name is not an occurrence of a variable: neither status.count nor statusType{count: 1} uses the variable count.
func countUses(n ast.Node, name string) int {
	count := 0

	ast.Inspect(n, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == name {
				count++
			}
		case *ast.SelectorExpr:
			count += countUses(n.X, name) // n.Sel names a field or a method
			return false
		case *ast.CompositeLit:
			count += countUsesInCompositeLit(n, name)
			return false
		}

		return true
	})

	return count
}

// countUsesInCompositeLit returns the number of occurrences of the variable name in the given composite literal.
// Only a map literal has keys that are values: elsewhere a key names a field or indexes an array, and holds no variable.
func countUsesInCompositeLit(lit *ast.CompositeLit, name string) int {
	count := 0
	if lit.Type != nil {
		count += countUses(lit.Type, name)
	}

	_, keysAreValues := lit.Type.(*ast.MapType)
	for _, element := range lit.Elts {
		keyValue, ok := element.(*ast.KeyValueExpr)
		switch {
		case !ok:
			count += countUses(element, name)
		case keysAreValues:
			count += countUses(keyValue.Key, name) + countUses(keyValue.Value, name)
		default:
			count += countUses(keyValue.Value, name)
		}
	}

	return count
}

// countUsesInStmts returns the number of occurrences of the identifier name in the given statements.
func countUsesInStmts(stmts []ast.Stmt, name string) int {
	count := 0
	for _, stmt := range stmts {
		count += countUses(stmt, name)
	}

	return count
}
