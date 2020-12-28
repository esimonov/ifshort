package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var maxDeclChars, maxDeclLines int

const (
	maxDeclLinesUsage = `maximum length of variable declaration measured in number of lines, after which the linter won't suggest using short syntax. Has precedence over max-decl-chars.`
	maxDeclCharsUsage = `maximum length of variable declaration measured in number of characters, after which the linter won't suggest using short syntax.`
)

func init() {
	Analyzer.Flags.IntVar(&maxDeclLines, "max-decl-lines", 1, maxDeclLinesUsage)
	Analyzer.Flags.IntVar(&maxDeclChars, "max-decl-chars", 30, maxDeclCharsUsage)
}

// Analyzer is an analysis.Analyzer instance for ifshort linter.
var Analyzer = &analysis.Analyzer{
	Name:     "ifshort",
	Doc:      "Checks that your code uses short syntax for if-statements whenever possible.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		fdecl := node.(*ast.FuncDecl)

		/*if fdecl.Name.Name != "notUsed_For_OK" {
			return
		}*/

		candidates := getNamedOccurrenceMap(fdecl, pass)

		for _, stmt := range fdecl.Body.List {
			candidates.checkStatement(stmt)
		}

		for varName := range candidates {
			for marker, occ := range candidates[varName] {
				//  If two or more vars with the same lhs marker - skip them.
				if candidates.isFoundByLhsMarker(marker) {
					continue
				}

				pass.Reportf(occ.declarationPos,
					"variable '%s' is only used in the if-statement (%s); consider using short syntax",
					varName, pass.Fset.Position(occ.ifStmtPos))
			}
		}
	})
	return nil, nil
}

func (nom namedOccurrenceMap) checkStatement(stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.AssignStmt:
		for _, el := range v.Rhs {
			nom.check(el)
		}
	case *ast.DeferStmt:
		for _, a := range v.Call.Args {
			nom.check(a)
		}
	case *ast.ExprStmt:
		if callExpr, ok := v.X.(*ast.CallExpr); ok {
			nom.check(callExpr)
		}
	case *ast.IfStmt:
		switch cond := v.Cond.(type) {
		case *ast.BinaryExpr:
			nom.checkIf(cond.X, v.If)
			nom.checkIf(cond.Y, v.If)
		case *ast.CallExpr:
			nom.checkIf(cond, v.If)
		}
		if init, ok := v.Init.(*ast.AssignStmt); ok {
			for _, e := range init.Rhs {
				nom.checkIf(e, v.If)
			}
		}
	case *ast.IncDecStmt:
		nom.check(v.X)
	case *ast.ForStmt:
		for _, el := range v.Body.List {
			nom.checkStatement(el)
		}
		if bexpr, ok := v.Cond.(*ast.BinaryExpr); ok {
			nom.check(bexpr.X)
			nom.check(bexpr.Y)
		}
		nom.checkStatement(v.Post)
	case *ast.GoStmt:
		for _, a := range v.Call.Args {
			nom.check(a)
		}
	case *ast.RangeStmt:
		nom.check(v.X)
	case *ast.ReturnStmt:
		for _, r := range v.Results {
			nom.check(r)
		}
	case *ast.SendStmt:
		nom.check(v.Value)
	case *ast.SwitchStmt:
		nom.check(v.Tag)
		for _, el := range v.Body.List {
			if clauses, ok := el.(*ast.CaseClause); ok {
				for _, c := range clauses.List {
					switch v := c.(type) {
					case *ast.BinaryExpr:
						nom.check(v.X)
						nom.check(v.Y)
					case *ast.Ident:
						nom.check(v)
					}
				}
				for _, c := range clauses.Body {
					if est, ok := c.(*ast.ExprStmt); ok {
						nom.check(est.X)
					}
				}
			}
		}
	}
}

func (nom namedOccurrenceMap) checkIf(candidate ast.Expr, ifPos token.Pos) {
	switch v := candidate.(type) {
	case *ast.CallExpr:
		for _, arg := range v.Args {
			nom.checkIf(arg, ifPos)
		}
		if fun, ok := v.Fun.(*ast.SelectorExpr); ok {
			nom.checkIf(fun.X, ifPos)
		}
	case *ast.Ident:
		for lhsMarker1 := range nom[v.Name] {
			if !nom[v.Name].isEmponymousKey(ifPos) {
				delete(nom[v.Name], lhsMarker1)
				for k, v := range nom {
					for lhsMarker2 := range v {
						if lhsMarker1 == lhsMarker2 {
							delete(nom, k)
						}
					}
				}
			}
		}
	case *ast.UnaryExpr:
		nom.checkIf(v.X, ifPos)
	}
}

func (nom namedOccurrenceMap) check(candidate ast.Expr) {
	switch v := candidate.(type) {
	case *ast.CallExpr:
		for _, arg := range v.Args {
			nom.check(arg)
		}
		if fun, ok := v.Fun.(*ast.SelectorExpr); ok {
			nom.check(fun.X)
		}
	case *ast.CompositeLit:
		for _, el := range v.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if ident, ok := kv.Key.(*ast.Ident); ok {
				nom.check(ident)
			}
			if ident, ok := kv.Value.(*ast.Ident); ok {
				nom.check(ident)
			}
		}
	case *ast.Ident:
		lhsMarker1 := nom[v.Name].getLhsMarkerForPos(v.Pos())
		occ := nom[v.Name][lhsMarker1]
		if v.Pos() != occ.ifStmtPos && v.Pos() != occ.declarationPos {
			delete(nom[v.Name], lhsMarker1)
			for k := range nom {
				for lhsMarker2 := range nom[k] {
					if lhsMarker1 == lhsMarker2 {
						delete(nom[k], lhsMarker2)
					}
				}
			}
		}
	case *ast.IndexExpr:
		nom.check(v.X)
		index, ok := v.Index.(*ast.BinaryExpr)
		if !ok {
			return
		}
		nom.check(index.X)
	case *ast.SelectorExpr:
		nom.check(v.X)
	case *ast.UnaryExpr:
		nom.check(v.X)
	}
}
