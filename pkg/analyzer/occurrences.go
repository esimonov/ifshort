package analyzer

import (
	"go/ast"
	"go/token"
	"time"

	"golang.org/x/tools/go/analysis"
)

type occurrenceInfo struct {
	declarationPos token.Pos
	ifStmtPos      token.Pos
	lhsMarker      int64
}

func getOccurrenceMap(fdecl *ast.FuncDecl, pass *analysis.Pass) map[string]occurrenceInfo {
	occs := map[string]occurrenceInfo{}

	for _, stmt := range fdecl.Body.List {
		switch v := stmt.(type) {
		case *ast.AssignStmt:
			addOccurrencesFromAssignment(pass, v, occs)
		case *ast.IfStmt:
			addOccurrenceFromCondition(v.Cond, occs)
			addOccurrenceFromIfClause(v.Body, occs)
			addOccurrenceFromElseClause(v.Else, occs)
		}
	}
	return occs
}

func addOccurrencesFromAssignment(pass *analysis.Pass, assignment *ast.AssignStmt, occs map[string]occurrenceInfo) {
	lhsMarker := time.Now().UnixNano()

	for i, el := range assignment.Lhs {
		lhsIdent, ok := el.(*ast.Ident)
		if !ok {
			continue
		}

		if lhsIdent.Name != "_" && lhsIdent.Obj != nil && lhsIdent.Obj.Pos() == lhsIdent.Pos() {
			if oi, ok := occs[lhsIdent.Name]; ok {
				oi.declarationPos = lhsIdent.Pos()
				occs[lhsIdent.Name] = oi
			} else if areFlagSettingsSatisfied(pass, assignment, i) {
				occs[lhsIdent.Name] = occurrenceInfo{
					declarationPos: lhsIdent.Pos(),
					lhsMarker:      lhsMarker,
				}
			}
		}
	}
}

func areFlagSettingsSatisfied(pass *analysis.Pass, assignment *ast.AssignStmt, i int) bool {
	lh := assignment.Lhs[i]
	rh := assignment.Rhs[len(assignment.Rhs)-1]

	if len(assignment.Rhs) == len(assignment.Lhs) {
		rh = assignment.Rhs[i]
	}

	if pass.Fset.Position(rh.End()).Line-pass.Fset.Position(rh.Pos()).Line > maxDeclLines {
		return false
	}
	if int(rh.End()-lh.Pos()) > maxDeclChars {
		return false
	}
	return true
}

func addOccurrenceFromCondition(stmt ast.Expr, occs map[string]occurrenceInfo) {
	switch v := stmt.(type) {
	case *ast.BinaryExpr:
		for _, v := range [2]ast.Expr{v.X, v.Y} {
			switch e := v.(type) {
			case *ast.Ident:
				addOccurrenceFromIdent(occs, e)
			case *ast.SelectorExpr:
				addOccurrenceFromIdent(occs, e.X)
			}
		}
	case *ast.CallExpr:
		for _, a := range v.Args {
			switch e := a.(type) {
			case *ast.Ident:
				addOccurrenceFromIdent(occs, e)
			case *ast.CallExpr:
				addOccurrenceFromCallExpr(occs, e)
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

		if callExpr, ok := exptStmt.X.(*ast.CallExpr); ok {
			addOccurrenceFromCallExpr(occs, callExpr)
		}
	}
}

func addOccurrenceFromCallExpr(occs map[string]occurrenceInfo, callExpr *ast.CallExpr) {
	for _, arg := range callExpr.Args {
		addOccurrenceFromIdent(occs, arg)
	}
}

func addOccurrenceFromIdent(occs map[string]occurrenceInfo, v ast.Expr) {
	if ident, ok := v.(*ast.Ident); ok {
		if oi, ok := occs[ident.Name]; ok {
			if oi.ifStmtPos != 0 && oi.declarationPos != 0 {
				return
			}

			oi.ifStmtPos = v.Pos()
			occs[ident.Name] = oi
		}
	}
}
