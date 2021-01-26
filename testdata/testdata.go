package testdata

import (
	"sync"
)

// Cases where short syntax SHOULD be used AND IS used.

func used_CondBinary_UsedInBody_OK() {
	if v := getValue(); v != nil {
		noOp1(v)
	}
}

func used_CondBinary_UsedInElse_OK() {
	if v := getValue(); v != nil {
	} else {
		noOp2(v)
	}
}

func used_CondBinary_UsedInBodyAndElse_OK() {
	if v := getValue(); v != nil {
		noOp1(v)
	} else {
		noOp2(v)
	}
}

// Cases where short syntax SHOULD be used BUT is NOT used.

func notUsed_CondBinaryExpr_NotOK() {
	v := getValue() // want "variable '.+' is only used in the if-statement"
	if v != nil {
		noOp1(v)
	}
}

func notUsed_CondCallExpr_NotOK() {
	a := getValue() // want "variable '.+' is only used in the if-statement"
	if getBool(a) {
	}

	b := getValue() // want "variable '.+' is only used in the if-statement"
	if getBool(getValue(b)) {
	}
}

func notUsed_Body_NotOK() {
	v := getValue() // want "variable '.+' is only used in the if-statement"
	if true {
		noOp1(v)
	}
}

func notUsed_Else_NotOK() {
	v := getValue() // want "variable '.+' is only used in the if-statement"
	if true {
	} else {
		noOp2(v)
	}
}

func notUsed_DifferentVarsWithSameName_NotOK() {
	_, b := getTwoValues() // want "variable '.+' is only used in the if-statement"
	if b != nil {
		noOp1(b)
	}

	a, b := getTwoValues()
	if b != nil {
		noOp1(a)
	}
	noOp2(b)
}

// Cases where short syntax SHOULD NOT be used AND IS NOT used.

func notUsed_DeferStmt_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	defer noOp2(v)
}

func notUsed_IfStmt_CondBinaryExpr_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	if v == nil {
		noOp2(v)
	}
}

func notUsed_IfStmt_CondBinaryExpr_MethodCall_OK() {
	dt := getDummy()
	if v := dt.getValue(); v == nil {
		noOp1(v)
	}
	if dt.v != nil {
		noOp2(dt)
	}
}

func notUsed_IfStmt_CondCallExpr_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	if getValue(v) != nil {
		noOp2(v)
	}
}

func notUsed_IfStmt_Body_OK() {
	a, b := 0, 0

	if a <= 0 {
		noOp1(a)
	} else if a > 0 {
		noOp2(a)
	}
	if b > 0 {
		noOp1(a)
		for a < 0 {
			noOp2(a)
			a++
		}
		noOp1(a)
	}
}

func notUsed_IfStmt_AssignInLeftHandSide_OK() {
	a := 0

	if b := 0; b > 0 {
		a = 0
	}
	if a > 0 {
		return
	}
}

func notUsed_GoStmt_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	go noOp2(v)
}

func notUsed_ReturnStmt_OK() interface{} {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	return v
}

func notUsed_SendStmt_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}

	vChan := make(chan interface{}, 1)
	vChan <- v
}

func notUsed_SwitchStmt_Tag_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	switch v {
	case nil:
	}
}

func notUsed_SwitchStmt_CaseList_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	switch {
	case v == nil:
	}
}

func notUsed_SwitchStmt_CaseBody_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	switch {
	case true:
		noOp2(v)
	}
}

func notUsed_SwitchStmt_Body_OK() {
	a := getValue()
	if a != nil {
		noOp1(a)
	}
	b := getValue()
	switch b {
	case a:
	}
}

func notUsed_SwitchStmt_Body_Assignment_OK() {
	size := 0
	if size == 0 {
		return
	}

	switch 0 {
	case 0:
		a := make([]byte, size-2)
		noOp1(a)
	}
}

func notUsed_CompositeLiteral_OK() map[int]struct{} {
	a := 0
	if a != 0 {
		return nil
	}

	b := struct{}{}

	return map[int]struct{}{a: b}
}

func notUsed_MultipleAssignments_OK() interface{} {
	a, b := getTwoValues()
	if a != nil {
		return a
	}
	return b
}

func notUsed_MultipleAssignments_AllUsesInIfs_OK() interface{} {
	a, b := getTwoValues()
	if a != nil {
		return a
	}
	if b != nil {
		return b
	}
	return nil
}

func notUsed_MultipleAssignments_WhenFlagSettingsAreNotSatisfied_OK() {
	longDeclarationToDissatisfyFlagSettings, b := getTwoValues()
	if b != nil {
		return
	}

	c := getValue(longDeclarationToDissatisfyFlagSettings)
	noOp1(c)
}

func notUsed_LongDecl_OK() {
	v := getValue("Long long long long long declaration, linter shouldn't force short syntax for it, at least I hope so.")
	if v != nil {
		noOp1(v)
	}
}

func notUsed_HighDecl_OK() {
	v := getValue(
		nil,
		nil,
		nil,
	)
	if v != nil {
		noOp1(v)
	}
}

func notUsed_MethodCall_OK() {
	d := dummyType{}
	if d.v == nil {
	}
	d.noOp()
}

func notUsed_MethodCallWithAssignment_OK() {
	d := dummyType{}
	if d.v != nil {
	}

	v := d.getValue()
	noOp1(v)
}

func notUsed_MethodCall_Nested_OK() {
	d := dummyType{}
	if d.v != nil {
	}
	noOp1(d.getValue())
}

func notUsed_Pointer_OK() {
	d := dummyType{}
	if d.v != nil {
	}
	noOp1(&d)
}

func notUsed_CondMethodCall_OK() {
	d := dummyType{}
	if d.getValue() == nil {
	}
}

func notUsed_Range_OK() {
	ds := []dummyType{}
	if ds == nil {
	}

	for _, d := range ds {
		d.noOp()
	}
}

func notUsed_ForCond_OK(i int) {
	for i < 0 {
		break
	}
	if i == 0 {
		return
	}
}

func notUsed_ForBody_OK(t int) {
	d := getDummy()
	if d.v == nil {
		return
	}

	for i := 0; i < t; i++ {
		noOp1(d.v)
	}
}

func notUsed_ForPost_OK() {
	i := 0
	for ; ; i++ {
		break
	}

	if i == 0 {
		return
	}
}

func notUsed_IncrementDecrement_OK() {
	i := 0
	i++
	if i == 0 {
		return
	}
	i--
}

func notUsed_AssignToField_OK() {
	d := dummyType{}
	d.v = getValue()
}

func notUsed_ReferenceToFields_OK() {
	a, b := getTwoDummies()
	if getBool(a.v) {
		return
	}
	noOp1(b.v)
}

func notUsed_IndexExpression_Index_OK() {
	s := []int{}

	length := len(s)
	if length == 0 {
		return
	}

	last := s[length-1]
	noOp1(last)
}

func notUsed_IndexExpression_Indexed_OK() {
	s := []int{}
	if s == nil {
		return
	}

	first := s[0]
	noOp1(first)
}

func notUsed_BinaryExpressionInIndex_OK(size int) {
	if size == 0 {
		return
	}

	a := make([]byte, size-1)
	noOp1(a)
}

func notUsed_SliceExpression_Low_OK(size int) {
	s := []int{}
	if size != 0 {
		return
	}
	noOp2(s[size:])
}

func notUsed_SliceExpression_High_OK(size int) {
	s := []int{}
	if size != 0 {
		return
	}
	noOp1(s[:size])
}

func notUsed_FuncLitReturn_OK() {
	s := ""
	tp := func() string { return s }

	tp()

	if s != "" {
		return
	}
}

func loopVar_OK() {
	wg := &sync.WaitGroup{}
	for range []int{1, 2, 3} {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}

	if func(wg *sync.WaitGroup) bool {
		wg.Wait()
		return true
	}(wg) {
		return
	}
}
