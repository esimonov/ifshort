package analyzer

import (
	"go/ast"
	"go/token"
	"time"

	"golang.org/x/tools/go/analysis"
)

// occurrence is a variable occurrence.
type occurrence struct {
	declarationPos token.Pos
	ifStmtPos      token.Pos
}

// lhsMarkeredOccurences is a map of left-hand side markers to occurrence.
type lhsMarkeredOccurences map[int64]occurrence

// namedOccurrenceMap is a map of variable names to lhsMarkeredOccurences.
type namedOccurrenceMap map[string]lhsMarkeredOccurences

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

func (lmo lhsMarkeredOccurences) getLatestLhsMarker() int64 {
	var maxLhsMarker int64

	for marker := range lmo {
		if marker > maxLhsMarker {
			maxLhsMarker = marker
		}
	}
	return maxLhsMarker
}

func getNamedOccurrenceMap(fdecl *ast.FuncDecl, pass *analysis.Pass) namedOccurrenceMap {
	occs := namedOccurrenceMap(map[string]lhsMarkeredOccurences{})

	for _, stmt := range fdecl.Body.List {
		switch v := stmt.(type) {
		case *ast.AssignStmt:
			occs.addFromAssignment(pass, v)
		case *ast.IfStmt:
			occs.addFromCondition(v)
			occs.addFromIfClause(v)
			occs.addFromElseClause(v)
		}
	}

	candidates := namedOccurrenceMap(map[string]lhsMarkeredOccurences{})

	for varName, occ := range occs {
		for lhs, o := range occ {
			if o.declarationPos != token.NoPos || occs.isFoundByLhsMarker(lhs) {
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

func (nom namedOccurrenceMap) isFoundByLhsMarker(lhsMarker int64) bool {
	var i int
	for _, markeredOcc := range nom {
		for marker := range markeredOcc {
			if marker == lhsMarker {
				i++
			}
		}
	}
	return i >= 2
}

func (nom namedOccurrenceMap) addFromAssignment(pass *analysis.Pass, assignment *ast.AssignStmt) {
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
			if oi, ok := nom[lhsIdent.Name]; ok {
				oi[lhsMarker] = occurrence{
					declarationPos: lhsIdent.Pos(),
				}
				nom[lhsIdent.Name] = oi
			} else {
				newOcc := occurrence{}
				if areFlagSettingsSatisfied(pass, assignment, i) {
					newOcc.declarationPos = lhsIdent.Pos()
				}
				nom[lhsIdent.Name] = lhsMarkeredOccurences{lhsMarker: newOcc}
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

func (nom namedOccurrenceMap) addFromCondition(stmt *ast.IfStmt) {
	switch v := stmt.Cond.(type) {
	case *ast.BinaryExpr:
		for _, v := range [2]ast.Expr{v.X, v.Y} {
			switch e := v.(type) {
			case *ast.Ident:
				nom.addFromIdent(stmt.If, e)
			case *ast.SelectorExpr:
				nom.addFromIdent(stmt.If, e.X)
			}
		}
	case *ast.CallExpr:
		for _, a := range v.Args {
			switch e := a.(type) {
			case *ast.Ident:
				nom.addFromIdent(stmt.If, e)
			case *ast.CallExpr:
				nom.addFromCallExpr(stmt.If, e)
			}
		}
	}
}

func (nom namedOccurrenceMap) addFromIfClause(stmt *ast.IfStmt) {
	nom.addFromBlockStmt(stmt.Body, stmt.If)
}

func (nom namedOccurrenceMap) addFromElseClause(stmt *ast.IfStmt) {
	nom.addFromBlockStmt(stmt.Else, stmt.If)
}

func (nom namedOccurrenceMap) addFromBlockStmt(stmt ast.Stmt, ifPos token.Pos) {
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
			nom.addFromCallExpr(ifPos, callExpr)
		}
	}
}

func (nom namedOccurrenceMap) addFromCallExpr(ifPos token.Pos, callExpr *ast.CallExpr) {
	for _, arg := range callExpr.Args {
		nom.addFromIdent(ifPos, arg)
	}
}

func (nom namedOccurrenceMap) addFromIdent(ifPos token.Pos, v ast.Expr) {
	ident, ok := v.(*ast.Ident)
	if !ok {
		return
	}

	if markeredOcc, ok := nom[ident.Name]; ok {
		marker := nom[ident.Name].getLatestLhsMarker()

		occ := markeredOcc[marker]
		if occ.ifStmtPos != token.NoPos && occ.declarationPos != token.NoPos {
			return
		}

		occ.ifStmtPos = ifPos
		nom[ident.Name][marker] = occ
	}
}
