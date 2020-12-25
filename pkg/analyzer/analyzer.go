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

		/*if fdecl.Name.Name != "notUsed_MultipleAssignments_WhenFlagSettingsAreNotSatisfied_OK" {
			return
		}*/

		candidates := map[string]map[int64]occurrenceInfo{}
		occurrences := getOccurrenceMap(fdecl, pass)

		for varName, occ := range occurrences {
			for lhs, o := range occ {
				if o.declarationPos != 0 || isFoundByLhsMarker(occurrences, lhs) {
					if _, ok := candidates[varName]; !ok {
						candidates[varName] = map[int64]occurrenceInfo{
							lhs: o,
						}
					} else {
						candidates[varName][lhs] = o
					}
				}
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
			for _, o := range candidates[varName] {
				pass.Reportf(o.declarationPos,
					"variable '%s' is only used in the if-statement (%s); consider using short syntax",
					varName, pass.Fset.Position(o.ifStmtPos))
			}
		}
	})
	return nil, nil
}

func isFoundByLhsMarker(candidates map[string]map[int64]occurrenceInfo, lhsMarker int64) bool {
	var i int
	for _, v := range candidates {
		for lhs := range v {
			if lhs == lhsMarker {
				i++
			}
		}
	}
	return i >= 2
}

func checkIfCandidate(candidates map[string]map[int64]occurrenceInfo, e ast.Expr, ifPos token.Pos) {
	switch v := e.(type) {
	case *ast.CallExpr:
		for _, arg := range v.Args {
			checkIfCandidate(candidates, arg, ifPos)
		}
		if fun, ok := v.Fun.(*ast.SelectorExpr); ok {
			checkIfCandidate(candidates, fun.X, ifPos)
		}
	case *ast.Ident:
		for _, occ := range candidates[v.Name] {
			if !isEmponymousKey(ifPos, candidates[v.Name]) {
				lhsMarker := occ.lhsMarker
				delete(candidates[v.Name], lhsMarker)
				for k, v := range candidates {
					for _, o := range v {
						if o.lhsMarker == lhsMarker {
							delete(candidates, k)
						}
					}
				}
			}
		}
	case *ast.UnaryExpr:
		checkIfCandidate(candidates, v.X, ifPos)
	}
}

func isEmponymousKey(pos token.Pos, occs map[int64]occurrenceInfo) bool {
	for _, o := range occs {
		if o.ifStmtPos == pos {
			return true
		}
	}
	return false
}

func checkCandidate(candidates map[string]map[int64]occurrenceInfo, e ast.Expr) {
	switch v := e.(type) {
	case *ast.CallExpr:
		for _, arg := range v.Args {
			checkCandidate(candidates, arg)
		}
		if fun, ok := v.Fun.(*ast.SelectorExpr); ok {
			checkCandidate(candidates, fun.X)
		}
	case *ast.CompositeLit:
		for _, el := range v.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if ident, ok := kv.Key.(*ast.Ident); ok {
				checkCandidate(candidates, ident)
			}
			if ident, ok := kv.Value.(*ast.Ident); ok {
				checkCandidate(candidates, ident)
			}
		}
	case *ast.Ident:
		lhsMarker := getLhsMarker(candidates[v.Name], v.Pos())
		occ := candidates[v.Name][lhsMarker]
		if v.Pos() != occ.ifStmtPos && v.Pos() != occ.declarationPos {
			delete(candidates[v.Name], lhsMarker)
			for k := range candidates {
				for _, occ2 := range candidates[k] {
					if occ2.lhsMarker == lhsMarker {
						delete(candidates[k], lhsMarker)
					}
				}
			}
		}

	case *ast.UnaryExpr:
		checkCandidate(candidates, v.X)
	}
}

// find lhs marker of the greatest token.Pos that is smaller than provided.
func getLhsMarker(occs map[int64]occurrenceInfo, pos token.Pos) int64 {
	var m int64
	var foundPos token.Pos

	for lhsMarker, occ := range occs {
		if occ.declarationPos < pos && occ.declarationPos >= foundPos {
			m = lhsMarker
			foundPos = occ.declarationPos
		}
	}

	return m
}
