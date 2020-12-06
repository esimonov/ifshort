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

/*
https://medium.com/justforfunc/understanding-go-programs-with-go-parser-c4e88a6edb87

https://astexplorer.net/

https://disaev.me/p/writing-useful-go-analysis-linter/
*/

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		fdecl := node.(*ast.FuncDecl)

		/*if fdecl.Name.Name != "notUsed_CondCallExpr_NotOK" {
			return
		}*/

		candidates := map[string]occurrenceInfo{}
		occurrences := getOccurrenceMap(fdecl, pass)

		for varName, occ := range occurrences {
			if occ.ifStmtPos > occ.declarationPos && occ.declarationPos != 0 ||
				isFoundByLhsMarker(occurrences, occ.lhsMarker) {
				candidates[varName] = occ
			}
		}

		for _, stmt := range fdecl.Body.List {
			switch v := stmt.(type) {
			case *ast.AssignStmt:
				for _, el := range v.Rhs {
					checkCandidate(candidates, el)
				}
			case *ast.DeferStmt:
				for _, a := range v.Call.Args {
					checkCandidate(candidates, a)
				}
			case *ast.ExprStmt:
				if callExpr, ok := v.X.(*ast.CallExpr); ok {
					checkCandidate(candidates, callExpr)
				}
			case *ast.IfStmt:
				switch cond := v.Cond.(type) {
				case *ast.BinaryExpr:
					checkIfCandidate(candidates, cond.X, v.If)
					checkIfCandidate(candidates, cond.Y, v.If)
				case *ast.CallExpr:
					checkIfCandidate(candidates, cond, v.If)
				}
				if init, ok := v.Init.(*ast.AssignStmt); ok {
					for _, e := range init.Rhs {
						checkIfCandidate(candidates, e, v.If)
					}
				}
			case *ast.GoStmt:
				for _, a := range v.Call.Args {
					checkCandidate(candidates, a)
				}
			case *ast.RangeStmt:
				checkCandidate(candidates, v.X)
			case *ast.ReturnStmt:
				for _, r := range v.Results {
					checkCandidate(candidates, r)
				}
			case *ast.SendStmt:
				checkCandidate(candidates, v.Value)
			case *ast.SwitchStmt:
				checkCandidate(candidates, v.Tag)
				for _, el := range v.Body.List {
					if clauses, ok := el.(*ast.CaseClause); ok {
						for _, c := range clauses.List {
							switch v := c.(type) {
							case *ast.BinaryExpr:
								checkCandidate(candidates, v.X)
								checkCandidate(candidates, v.Y)
							case *ast.Ident:
								checkCandidate(candidates, v)
							}
						}
						for _, c := range clauses.Body {
							if est, ok := c.(*ast.ExprStmt); ok {
								checkCandidate(candidates, est.X)
							}
						}
					}
				}
			}
		}

		for varName := range candidates {
			pass.Reportf(candidates[varName].declarationPos,
				"variable '%s' is only used in the if-statement (%s); consider using short syntax",
				varName, pass.Fset.Position(candidates[varName].ifStmtPos))
		}
	})
	return nil, nil
}

func isFoundByLhsMarker(candidates map[string]occurrenceInfo, lhsMarker int64) bool {
	var i int
	for _, v := range candidates {
		if v.lhsMarker == lhsMarker {
			i++
		}
	}
	return i >= 2
}

func checkIfCandidate(candidates map[string]occurrenceInfo, e ast.Expr, ifPos token.Pos) {
	switch v := e.(type) {
	case *ast.CallExpr:
		for _, arg := range v.Args {
			checkIfCandidate(candidates, arg, ifPos)
		}
		if fun, ok := v.Fun.(*ast.SelectorExpr); ok {
			checkIfCandidate(candidates, fun.X, ifPos)
		}
	case *ast.Ident:
		if _, ok := candidates[v.Name]; !ok {
			return
		}
		if ifPos != candidates[v.Name].ifStmtPos {
			lhsMarker := candidates[v.Name].lhsMarker
			delete(candidates, v.Name)
			for k, v := range candidates {
				if v.lhsMarker == lhsMarker {
					delete(candidates, k)
				}
			}
		}
	case *ast.UnaryExpr:
		checkIfCandidate(candidates, v.X, ifPos)
	}
}

func checkCandidate(candidates map[string]occurrenceInfo, e ast.Expr) {
	switch v := e.(type) {
	case *ast.CallExpr:
		for _, arg := range v.Args {
			checkCandidate(candidates, arg)
		}
		if fun, ok := v.Fun.(*ast.SelectorExpr); ok {
			checkCandidate(candidates, fun.X)
		}
	case *ast.Ident:
		if _, ok := candidates[v.Name]; !ok {
			return
		}
		if v.Pos() != candidates[v.Name].ifStmtPos && v.Pos() != candidates[v.Name].declarationPos {
			lhsMarker := candidates[v.Name].lhsMarker
			delete(candidates, v.Name)
			for k, v := range candidates {
				if v.lhsMarker == lhsMarker {
					delete(candidates, k)
				}
			}
		}
	case *ast.UnaryExpr:
		checkCandidate(candidates, v.X)
	}
}
