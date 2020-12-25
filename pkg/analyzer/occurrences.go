package analyzer

import (
	"go/ast"
	"go/token"
	"time"

	"golang.org/x/tools/go/analysis"
)

type occurrence struct {
	declarationPos token.Pos
	ifStmtPos      token.Pos
}

type lhsMarkeredOccurences map[int64]occurrence

// find lhs marker of the greatest token.Pos that is smaller than provided.
func (lmo lhsMarkeredOccurences) getLhsMarker(pos token.Pos) int64 {
	var m int64
	var foundPos token.Pos

	for lhsMarker, occ := range lmo {
		if occ.declarationPos < pos && occ.declarationPos >= foundPos {
			m = lhsMarker
			foundPos = occ.declarationPos
		}
	}
	return m
}

func getOccurrenceMap(fdecl *ast.FuncDecl, pass *analysis.Pass) map[string]lhsMarkeredOccurences {
	occs := map[string]lhsMarkeredOccurences{}

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

	candidates := map[string]lhsMarkeredOccurences{}

	for varName, occ := range occs {
		for lhs, o := range occ {
			if o.declarationPos != token.NoPos || isFoundByLhsMarker(occs, lhs) {
				if _, ok := candidates[varName]; !ok {
					candidates[varName] = lhsMarkeredOccurences{
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

func isFoundByLhsMarker(candidates map[string]lhsMarkeredOccurences, lhsMarker int64) bool {
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

func addOccurrencesFromAssignment(pass *analysis.Pass, assignment *ast.AssignStmt, occs map[string]lhsMarkeredOccurences) {
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
				oi[lhsMarker] = occurrence{
					declarationPos: lhsIdent.Pos(),
				}
				occs[lhsIdent.Name] = oi
			} else {
				newOcc := occurrence{}
				if areFlagSettingsSatisfied(pass, assignment, i) {
					newOcc.declarationPos = lhsIdent.Pos()
				}
				occs[lhsIdent.Name] = lhsMarkeredOccurences{lhsMarker: newOcc}
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

func addOccurrenceFromCondition(stmt *ast.IfStmt, occs map[string]lhsMarkeredOccurences) {
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

func addOccurrenceFromIfClause(stmt *ast.IfStmt, occs map[string]lhsMarkeredOccurences) {
	addOccurrenceFromBlockStmt(stmt.Body, stmt.If, occs)
}

func addOccurrenceFromElseClause(stmt *ast.IfStmt, occs map[string]lhsMarkeredOccurences) {
	addOccurrenceFromBlockStmt(stmt.Else, stmt.If, occs)
}

func addOccurrenceFromBlockStmt(stmt ast.Stmt, ifPos token.Pos, occs map[string]lhsMarkeredOccurences) {
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

func addOccurrenceFromCallExpr(occs map[string]lhsMarkeredOccurences, ifPos token.Pos, callExpr *ast.CallExpr) {
	for _, arg := range callExpr.Args {
		addOccurrenceFromIdent(occs, ifPos, arg)
	}
}

func addOccurrenceFromIdent(occs map[string]lhsMarkeredOccurences, ifPos token.Pos, v ast.Expr) {
	ident, ok := v.(*ast.Ident)
	if !ok {
		return
	}

	if oi, ok := occs[ident.Name]; ok {
		lhs := getLatestLhs(occs[ident.Name])

		o := oi[lhs]
		if o.ifStmtPos != token.NoPos && o.declarationPos != token.NoPos {
			return
		}

		o.ifStmtPos = ifPos
		occs[ident.Name][lhs] = o
	}
}

func getLatestLhs(o lhsMarkeredOccurences) int64 {
	var maxLhs int64
	for lhs := range o {
		if lhs > maxLhs {
			maxLhs = lhs
		}
	}
	return maxLhs
}
