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

func getOccurrenceMap(fdecl *ast.FuncDecl, pass *analysis.Pass) map[string]map[int64]occurrenceInfo {
	occs := map[string]map[int64]occurrenceInfo{}

	for _, stmt := range fdecl.Body.List {
		switch v := stmt.(type) {
		case *ast.AssignStmt:
			addOccurrencesFromAssignment(pass, v, occs)
		case *ast.IfStmt:
			addOccurrenceFromCondition(v, occs)
			addOccurrenceFromIfClause(v, occs)
			addOccurrenceFromElseClause(v, occs)
		}
	}

	candidates := map[string]map[int64]occurrenceInfo{}

	for varName, occ := range occs {
		for lhs, o := range occ {
			if o.declarationPos != 0 || isFoundByLhsMarker(occs, lhs) {
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
	return candidates
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

func addOccurrencesFromAssignment(pass *analysis.Pass, assignment *ast.AssignStmt, occs map[string]map[int64]occurrenceInfo) {
	if assignment.Tok != token.DEFINE {
		return
	}

	lhsMarker := time.Now().UnixNano()

	for i, el := range assignment.Lhs {
		lhsIdent, ok := el.(*ast.Ident)
		if !ok {
			continue
		}

		if lhsIdent.Name != "_" && lhsIdent.Obj != nil { //&& lhsIdent.Obj.Pos() == lhsIdent.Pos() {
			if oi, ok := occs[lhsIdent.Name]; ok {
				oi[lhsMarker] = occurrenceInfo{
					declarationPos: lhsIdent.Pos(),
					lhsMarker:      lhsMarker,
				}
				occs[lhsIdent.Name] = oi
			} else {
				newOcc := occurrenceInfo{lhsMarker: lhsMarker}
				if areFlagSettingsSatisfied(pass, assignment, i) {
					newOcc.declarationPos = lhsIdent.Pos()
				}
				occs[lhsIdent.Name] = map[int64]occurrenceInfo{lhsMarker: newOcc}
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

func addOccurrenceFromCondition(stmt *ast.IfStmt, occs map[string]map[int64]occurrenceInfo) {
	switch v := stmt.Cond.(type) {
	case *ast.BinaryExpr:
		for _, v := range [2]ast.Expr{v.X, v.Y} {
			switch e := v.(type) {
			case *ast.Ident:
				addOccurrenceFromIdent(occs, stmt.If, e)
			case *ast.SelectorExpr:
				addOccurrenceFromIdent(occs, stmt.If, e.X)
			}
		}
	case *ast.CallExpr:
		for _, a := range v.Args {
			switch e := a.(type) {
			case *ast.Ident:
				addOccurrenceFromIdent(occs, stmt.If, e)
			case *ast.CallExpr:
				addOccurrenceFromCallExpr(occs, stmt.If, e)
			}
		}
	}
}

func addOccurrenceFromIfClause(stmt *ast.IfStmt, occs map[string]map[int64]occurrenceInfo) {
	addOccurrenceFromBlockStmt(stmt.Body, stmt.If, occs)
}

func addOccurrenceFromElseClause(stmt *ast.IfStmt, occs map[string]map[int64]occurrenceInfo) {
	addOccurrenceFromBlockStmt(stmt.Else, stmt.If, occs)
}

func addOccurrenceFromBlockStmt(stmt ast.Stmt, ifPos token.Pos, occs map[string]map[int64]occurrenceInfo) {
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
			addOccurrenceFromCallExpr(occs, ifPos, callExpr)
		}
	}
}

func addOccurrenceFromCallExpr(occs map[string]map[int64]occurrenceInfo, ifPos token.Pos, callExpr *ast.CallExpr) {
	for _, arg := range callExpr.Args {
		addOccurrenceFromIdent(occs, ifPos, arg)
	}
}

func addOccurrenceFromIdent(occs map[string]map[int64]occurrenceInfo, ifPos token.Pos, v ast.Expr) {
	if ident, ok := v.(*ast.Ident); ok {
		if oi, ok := occs[ident.Name]; ok {
			lhs := getLatestLhs(occs[ident.Name])
			o := oi[lhs]
			if o.ifStmtPos != 0 && o.declarationPos != 0 {
				return
			}

			o.ifStmtPos = ifPos
			occs[ident.Name][lhs] = o
		}
	}
}

func getLatestLhs(o map[int64]occurrenceInfo) int64 {
	var maxLhs int64
	for lhs := range o {
		if lhs > maxLhs {
			maxLhs = lhs
		}
	}
	return maxLhs
}
