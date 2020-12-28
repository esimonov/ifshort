package testdata

import "errors"

// Cases where short syntax SHOULD be used AND IS used.

func used_CondBinary_UsedInBody_OK() {
	if v := returnValue(); v != nil {
		noOp1(v)
	}
}

func used_CondBinary_UsedInElse_OK() {
	if v := returnValue(); v != nil {
	} else {
		noOp2(v)
	}
}

func used_CondBinary_UsedInBodyAndElse_OK() {
	if v := returnValue(); v != nil {
		noOp1(v)
	} else {
		noOp2(v)
	}
}

// Cases where short syntax SHOULD be used BUT is NOT used.

func notUsed_CondBinaryExpr_NotOK() {
	v := returnValue() // want "variable '.+' is only used in the if-statement"
	if v != nil {
		noOp1(v)
	}
}

func notUsed_CondCallExpr_NotOK() {
	err := errors.New("") // want "variable '.+' is only used in the if-statement"
	if errors.Is(err, errors.New("")) {
	}
}

func notUsed_Body_NotOK() {
	v := returnValue() // want "variable '.+' is only used in the if-statement"
	if true {
		noOp1(v)
	}
}

func notUsed_Else_NotOK() {
	v := returnValue() // want "variable '.+' is only used in the if-statement"
	if true {
	} else {
		noOp2(v)
	}
}

func notUsed_DifferentVarsWithSameName_NotOK() {
	_, b := returnTwoValues() // want "variable '.+' is only used in the if-statement"
	if b != nil {
		noOp1(b)
	}

	a, b := returnTwoValues()
	if b != nil {
		noOp1(a)
	}
	noOp2(b)
}

// Cases where short syntax SHOULD NOT be used AND IS NOT used.

func notUsed_DeferStmt_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	defer noOp2(v)
}

func notUsed_IfStmt_CondBinaryExpr_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	if v == nil {
		noOp2(v)
	}
}

func notUsed_IfStmt_CondBinaryExpr_MethodCall_OK() {
	dt := returnDummy()
	if v := dt.returnValue(); v == nil {
		noOp1(v)
	}
	if dt.v != nil {
		noOp2(dt)
	}
}

func notUsed_IfStmt_CondCallExpr_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	if callWithOneArgAndReturnValue(v) != nil {
		noOp2(v)
	}
}

func notUsed_IfStmt_Body_OK(scale int) {
	pos := 0
	if pos <= 0 {
		noOp1(pos)
	} else if pos > 0 {
		noOp2(pos)
	}
	if scale > 0 {
		noOp1(pos)
		for pos < 0 {
			noOp2(pos)
			pos++
		}
		noOp1(pos)
	}
}

func notUsed_GoStmt_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	go noOp2(v)
}

func notUsed_ReturnStmt_OK() interface{} {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	return v
}

func notUsed_SendStmt_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}

	vChan := make(chan interface{}, 1)
	vChan <- v
}

func notUsed_SwitchStmt_Tag_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	switch v {
	case nil:
	}
}

func notUsed_SwitchStmt_CaseList_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	switch {
	case v == nil:
	}
}

func notUsed_SwitchStmt_CaseBody_OK() {
	v := returnValue()
	if v != nil {
		noOp1(v)
	}
	switch {
	case true:
		noOp2(v)
	}
}

func notUsed_SwitchStmt_Body_OK() {
	a := returnValue()
	if a != nil {
		noOp1(a)
	}
	b := returnValue()
	switch b {
	case a:
	}
}

func notUsed_MultipleAssignments_OK() interface{} {
	a, b := returnTwoValues()
	if a != nil {
		return a
	}
	return b
}

func notUsed_MultipleAssignments_AllUsesInIfs_OK() interface{} {
	a, b := returnTwoValues()
	if a != nil {
		return a
	}
	if b != nil {
		return b
	}
	return nil
}

func notUsed_MultipleAssignments_WhenFlagSettingsAreNotSatisfied_OK() {
	longDeclarationToDissatisfyFlagSettings, b := returnTwoValues()
	if b != nil {
		return
	}

	c := callWithOneArgAndReturnValue(longDeclarationToDissatisfyFlagSettings)
	noOp1(c)
}

func notUsed_LongDecl_OK() {
	v := callWithVariadicArgsAndReturnValue("Long long long long long declaration, linter shouldn't force short syntax for it")
	if v != nil {
		noOp1(v)
	}
}

func notUsed_HighDecl_OK() {
	v := callWithVariadicArgsAndReturnValue(
		nil,
		nil,
		nil,
	)
	if v != nil {
		noOp1(v)
	}
}

func notUsed_MethodCall_OK() {
	dt := dummyType{}
	if dt.v == nil {
	}
	dt.noOp()
}

func notUsed_MethodCallWithAssignment_OK() {
	dt := dummyType{}
	if dt.v != nil {
	}

	v := dt.returnValue()
	noOp1(v)
}

func notUsed_MethodCall_Nested_OK() {
	dt := dummyType{}
	if dt.v != nil {
	}
	noOp1(dt.returnValue())
}

func notUsed_Pointer_OK() {
	dt := dummyType{}
	if dt.v != nil {
	}
	noOp1(&dt)
}

func notUsed_CondMethodCall_OK() {
	dt := dummyType{}
	if dt.returnValue() == nil {
	}
}

func notUsed_Range_OK() {
	dts := []dummyType{}
	if dts == nil {
	}

	for _, dt := range dts {
		dt.noOp()
	}
}

func notUsed_ForCond_OK() {
	i := 0
	for i < 0 {
		break
	}

	if i == 0 {
		return
	}
}

func notUsed_ForBody_OK() {
	s := ""

	dt := returnDummy()
	if dt.v == nil {
		return
	}

	for i := 0; i < len(s); i++ {
		noOp1(dt.v)
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
	dt := dummyType{}
	dt.v = returnValue()
}

func notUsed_ReferenceToFields_OK() {
	a, b := returnTwoDummies()
	if a.v != nil {
		return
	}
	defer noOp1(b.v)
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

func notUsed_BinaryExpressionInIndex_OK() {
	s := []int{}
	size := 0
	if size != 0 {
		return
	}
	noOp1(s[size-1:])
}

func notUsed_SliceExpression_Low_OK() {
	s := []int{}
	size := 0
	if size != 0 {
		return
	}
	noOp2(s[size:])
}

func notUsed_SliceExpression_High_OK() {
	s := []int{}
	size := 0
	if size != 0 {
		return
	}
	noOp1(s[:size])
}
