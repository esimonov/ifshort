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

func notUsed_Var2_CondBinaryExpr_NotOK() {
	v := callWithVariadicArgsAndReturn(
		nil,
		nil,
		nil,
	)
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
	dt := returnDummyType()
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
	if callWithOneArgAndReturn(v) != nil {
		noOp2(v)
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

func notUsed_MultipleAssignments_OK() {
	a, b := returnTwoValues()
	if a != nil {
		return
	}
	noOp1(b)
}

func notUsed_LongDecl_OK() {
	v := callWithVariadicArgsAndReturn("Long long long long long declaration, linter shouldn't force short syntax for it")
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
