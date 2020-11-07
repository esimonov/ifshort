package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func getOccurrenceMap(fdecl *ast.FuncDecl, pass *analysis.Pass) map[string]occurrenceInfo {
	occs := map[string]occurrenceInfo{}

	for _, stmt := range fdecl.Body.List {
		switch v := stmt.(type) {
		case *ast.AssignStmt:
			trimmed := trimAssignmentSides(v.Lhs, v.Rhs)
			if len(trimmed) != 1 {
				continue
			}

			pair := trimmed[0]

			if pair.Lh.Obj != nil && pair.Lh.Obj.Pos() == pair.Lh.Pos() {
				if oi, ok := occs[pair.Lh.Name]; ok {
					oi.declarationPos = pair.Lh.Pos()
					occs[pair.Lh.Name] = oi
				} else if pair.Lh.Name != "nil" && areExtraConditionsSatisfied(pair, pass) {
					occs[pair.Lh.Name] = occurrenceInfo{declarationPos: pair.Lh.Pos()}
				}
			}
		case *ast.IfStmt:
			addOccurrenceFromCondition(v.Cond, occs)
			addOccurrenceFromIfClause(v.Body, occs)
			addOccurrenceFromElseClause(v.Else, occs)
		}
	}
	return occs
}

type assignmentSides struct {
	Lh *ast.Ident
	Rh ast.Expr
}

func trimAssignmentSides(lhs, rhs []ast.Expr) []assignmentSides {
	if len(lhs) != len(rhs) {
		return nil
	}

	res := make([]assignmentSides, 0, len(lhs))

	for i := 0; i < len(lhs); i++ {
		if lhsIdent, ok := lhs[i].(*ast.Ident); ok && lhsIdent.Name != "_" {
			res = append(res, assignmentSides{lhsIdent, rhs[i]})
		}
	}
	return res
}

func areExtraConditionsSatisfied(pair assignmentSides, pass *analysis.Pass) bool {
	if pass.Fset.Position(pair.Rh.End()).Line-pass.Fset.Position(pair.Rh.Pos()).Line > maxDeclHeight {
		return false
	}
	if pair.Rh.End()-pair.Lh.Pos() > maxDeclLength {
		return false
	}
	return true
}

func addOccurrenceFromCondition(stmt ast.Expr, occs map[string]occurrenceInfo) {
	switch v := stmt.(type) {
	case *ast.BinaryExpr:
		for _, v := range [2]ast.Expr{v.X, v.Y} {
			switch e := v.(type) {
			case *ast.SelectorExpr:
				processIdents(occs, e.X)
			case *ast.Ident:
				processIdents(occs, e)
			}
		}
	case *ast.CallExpr:
		for _, a := range v.Args {
			switch e := a.(type) {
			case *ast.CallExpr:
				// TODO
			case *ast.Ident:
				processIdents(occs, e)
			}
		}
	}
}

func addOccurrenceFromIfClause(stmt ast.Stmt, occs map[string]occurrenceInfo) {
	addOccurrenceFromBlockStmt(stmt, occs)
}

func addOccurrenceFromElseClause(stmt ast.Stmt, occs map[string]occurrenceInfo) {
	addOccurrenceFromBlockStmt(stmt, occs)
}

func addOccurrenceFromBlockStmt(stmt ast.Stmt, occs map[string]occurrenceInfo) {
	blockStmt, ok := stmt.(*ast.BlockStmt)
	if !ok {
		return
	}

	for _, el := range blockStmt.List {
		exptStmt, ok := el.(*ast.ExprStmt)
		if !ok {
			continue
		}

		callExpr, ok := exptStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		for _, arg := range callExpr.Args {
			if ident, ok := arg.(*ast.Ident); ok {
				if oi, ok := occs[ident.Name]; ok {
					if oi.ifStmtPos != 0 && oi.declarationPos != 0 {
						continue
					}

					oi.ifStmtPos = arg.Pos()
					occs[ident.Name] = oi
				} else if ident.Name != "nil" {
					occs[ident.Name] = occurrenceInfo{ifStmtPos: arg.Pos()}
				}
			}
		}
	}
}
