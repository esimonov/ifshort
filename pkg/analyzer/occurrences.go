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

// lhsMarkeredOccurences is a map of left-hand side markers to occurrences.
type lhsMarkeredOccurences map[int64]occurrence

func (lmo lhsMarkeredOccurences) getLatestLhsMarker() int64 {
	var maxLhsMarker int64

	for marker := range lmo {
		if marker > maxLhsMarker {
			maxLhsMarker = marker
		}
	}
	return maxLhsMarker
}

// find lhs marker of the greatest token.Pos that is smaller than provided.
func (lmo lhsMarkeredOccurences) getLhsMarkerForPos(pos token.Pos) int64 {
	var m int64
	var foundPos token.Pos

	for marker, occ := range lmo {
		if occ.declarationPos < pos && occ.declarationPos >= foundPos {
			m = marker
			foundPos = occ.declarationPos
		}
	}
	return m
}

func (lmo lhsMarkeredOccurences) isEmponymousKey(pos token.Pos) bool {
	for _, occ := range lmo {
		if occ.ifStmtPos == pos {
			return true
		}
	}
	return false
}

// namedOccurrenceMap is a map of variable names to lhsMarkeredOccurences.
type namedOccurrenceMap map[string]lhsMarkeredOccurences

func getNamedOccurrenceMap(fdecl *ast.FuncDecl, pass *analysis.Pass) namedOccurrenceMap {
	nom := namedOccurrenceMap(map[string]lhsMarkeredOccurences{})

	for _, stmt := range fdecl.Body.List {
		switch v := stmt.(type) {
		case *ast.AssignStmt:
			nom.addFromAssignment(pass, v)
		case *ast.IfStmt:
			nom.addFromCondition(v)
			nom.addFromIfClause(v)
			nom.addFromElseClause(v)
		}
	}

	candidates := namedOccurrenceMap(map[string]lhsMarkeredOccurences{})

	for varName, markeredOccs := range nom {
		for marker, occ := range markeredOccs {
			if occ.declarationPos != token.NoPos || nom.isFoundByLhsMarker(marker) {
				if _, ok := candidates[varName]; !ok {
					candidates[varName] = lhsMarkeredOccurences{
						marker: occ,
					}
				} else {
					candidates[varName][marker] = occ
				}
			}
		}
	}
	return candidates
}

func (nom namedOccurrenceMap) isFoundByLhsMarker(lhsMarker int64) bool {
	var i int
	for _, markeredOccs := range nom {
		for marker := range markeredOccs {
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
		ident, ok := el.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name != "_" && ident.Obj != nil {
			if markeredOccs, ok := nom[ident.Name]; ok {
				markeredOccs[lhsMarker] = occurrence{
					declarationPos: ident.Pos(),
				}
				nom[ident.Name] = markeredOccs
			} else {
				newOcc := occurrence{}
				if areFlagSettingsSatisfied(pass, assignment, i) {
					newOcc.declarationPos = ident.Pos()
				}
				nom[ident.Name] = lhsMarkeredOccurences{lhsMarker: newOcc}
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

	if markeredOccs, ok := nom[ident.Name]; ok {
		marker := nom[ident.Name].getLatestLhsMarker()

		occ := markeredOccs[marker]
		if occ.ifStmtPos != token.NoPos && occ.declarationPos != token.NoPos {
			return
		}

		occ.ifStmtPos = ifPos
		nom[ident.Name][marker] = occ
	}
}
