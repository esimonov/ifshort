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
	v := longCallWithReturnValue(
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
		noOp2(b)
	}
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
	err := errors.New("")
	if str := err.Error(); str != "" {
		noOp1(err)
	}
	if err != nil {
		noOp2(err)
	}
}

func notUsed_IfStmt_CondCallExpr_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	if errors.Is(err, errors.New("")) {
		noOp2(err)
	}
}

func notUsed_GoStmt_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	go noOp2(err)
}

func notUsed_ReturnStmt_OK() error {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	return err
}

func notUsed_SendStmt_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}

	errChan := make(chan error, 1)
	errChan <- err
}

func notUsed_SwitchStmt_Tag_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	switch err {
	case nil:
	}
}

func notUsed_SwitchStmt_CaseList_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	switch {
	case err == nil:
	}
}

func notUsed_SwitchStmt_CaseBody_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	switch {
	case true:
		noOp2(err)
	}
}

func notUsed_SwitchStmt_Body_OK() {
	err := errors.New("")
	if err != nil {
		noOp1(err)
	}
	err2 := errors.New("")
	switch err2 {
	case err:
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
	err := errors.New("Long long long long long declaration, linter shouldn't force short syntax for it")
	if err != nil {
		noOp1(err)
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

	err := dt.returnValue()
	noOp1(err)
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
