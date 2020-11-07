package testdata

import "errors"

// Cases where short syntax SHOULD be used AND IS used.

func used_CondBinary_UsedInBody_OK() {
	if v := dummyWithReturn(); v != nil {
		dummyNoOp1(v)
	}
}

func used_CondBinary_UsedInElse_OK() {
	if v := dummyWithReturn(); v != nil {
	} else {
		dummyNoOp2(v)
	}
}

func used_CondBinary_UsedInBodyAndElse_OK() {
	if v := dummyWithReturn(); v != nil {
		dummyNoOp1(v)
	} else {
		dummyNoOp2(v)
	}
}

// Cases where short syntax SHOULD be used BUT is NOT used.

func notUsed_CondBinaryExpr_NotOK() {
	v := dummyWithReturn() // want "variable '.+' is only used in the if-statement"
	if v != nil {
		dummyNoOp1(v)
	}
}

func notUsed_Var2_CondBinaryExpr_NotOK() {
	v := dummyForLongCall(
		nil,
		nil,
		nil,
	)
	if v != nil {
		dummyNoOp1(v)
	}
}

func notUsed_CondCallExpr_NotOK() {
	err := errors.New("") // want "variable '.+' is only used in the if-statement"
	if errors.Is(err, errors.New("")) {
	}
}

func notUsed_Body_NotOK() {
	v := dummyWithReturn() // want "variable '.+' is only used in the if-statement"
	if true {
		dummyNoOp1(v)
	}
}

func notUsed_Else_NotOK() {
	v := dummyWithReturn() // want "variable '.+' is only used in the if-statement"
	if true {
	} else {
		dummyNoOp2(v)
	}
}

// Cases where short syntax SHOULD NOT be used AND IS NOT used.

func notUsed_DeferStmt_OK() {
	v := dummyWithReturn()
	if v != nil {
		dummyNoOp1(v)
	}
	defer dummyNoOp2(v)
}

func notUsed_IfStmt_CondBinaryExpr_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	if err == nil {
		dummyNoOp2(err)
	}
}

func notUsed_IfStmt_CondBinaryExpr_MethodCall_OK() {
	err := errors.New("")
	if str := err.Error(); str != "" {
		dummyNoOp1(err)
	}
	if err != nil {
		dummyNoOp2(err)
	}
}

func notUsed_IfStmt_CondCallExpr_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	if errors.Is(err, errors.New("")) {
		dummyNoOp2(err)
	}
}

func notUsed_GoStmt_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	go dummyNoOp2(err)
}

func notUsed_ReturnStmt_OK() error {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	return err
}

func notUsed_SendStmt_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}

	errChan := make(chan error, 1)
	errChan <- err
}

func notUsed_SwitchStmt_Tag_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	switch err {
	case nil:
	}
}

func notUsed_SwitchStmt_CaseList_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	switch {
	case err == nil:
	}
}

func notUsed_SwitchStmt_CaseBody_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	switch {
	case true:
		dummyNoOp2(err)
	}
}

func notUsed_SwitchStmt_Body_OK() {
	err := errors.New("")
	if err != nil {
		dummyNoOp1(err)
	}
	err2 := errors.New("")
	switch err2 {
	case err:
	}
}

func notUsed_LongLhs_OK() {
	f := func() (string, error) { return "", nil }
	str, err := f()
	if str != "" {
		return
	}
	dummyNoOp1(err)
}

func notUsed_LongDecl_OK() {
	err := errors.New("Long long long long long declaration, linter shouldn't force short syntax for it yeeeeeeah")
	if err != nil {
		dummyNoOp1(err)
	}
}

func notUsed_MethodCall_OK() {
	dt := dummyType{}
	if dt.v == nil {
	}
	dt.dummyWithReturn()
}

func notUsed_MethodCallWithAssignment_OK() {
	dt := dummyType{}
	if dt.v != nil {
	}

	err := dt.dummyWithReturn()
	dummyNoOp1(err)
}

func notUsed_MethodCall_Nested_OK() {
	dt := dummyType{}
	if dt.v != nil {
	}
	dummyNoOp1(dt.dummyWithReturn())
}

func notUsed_Pointer_OK() {
	dt := dummyType{}
	if dt.v != nil {
	}
	dummyNoOp1(&dt)
}

func notUsed_CondMethodCall_OK() {
	dt := dummyType{}
	if dt.dummyWithReturn() == nil {
	}
}

func notUsed_Range_OK() {
	dts := []dummyType{}
	if dts == nil {
	}

	for _, dt := range dts {
		dt.dummyNoOp()
	}
}
