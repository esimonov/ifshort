package testdata

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

func notUsed_CondIdent_NotOK() {
	v := getBool() // want "variable '.+' is only used in the if-statement"
	if v {
		return
	}
}

func notUsed_CondCallWithIndentExpr_NotOK() {
	v := getInt() // want "variable '.+' is only used in the if-statement"
	if int(v) != 2 {
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

func notUsed_UnaryOpIfStatement_NotOK() {
	shouldRun := false // want "variable '.+' is only used in the if-statement"
	if !shouldRun {
		return
	}
	noOp1(0)
}

// Cases where short syntax SHOULD NOT be used AND IS NOT used.

func notUsed_DeferStmt_OK() {
	v := getValue()
	if v != nil {
		noOp1(v)
	}
	defer noOp2(v)
}

func notUsed_FuncDecl_OK(a func()) {
	v := a
	if v == nil {
		return
	}
	v()
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
	if dt.interf != nil {
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

func notUsed_SendStmt_Chan_OK(v interface{}) {
	ch := make(chan interface{})

	if ch == nil {
		return
	}
	ch <- v
}

func notUsed_SendStmt_Value_OK() {
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

func notUsed_CompositeLiteral_Map_OK() map[int]struct{} {
	a := 0
	if a != 0 {
		return nil
	}

	b := struct{}{}

	return map[int]struct{}{a: b}
}

func notUsed_CompositeLiteral_Struct_OK() dummyType {
	d := getDummy()
	if d.interf == 0 {
		return d
	}

	return dummyType{
		interf: getValue(d),
	}
}

func notUsed_CompositeLiteral_Array_OK() []interface{} {
	v := getValue()
	if v == nil {
		return nil
	}
	return []interface{}{v}
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
	if d.interf == nil {
	}
	d.noOp()
}

func notUsed_MethodCallWithAssignment_OK() {
	d := dummyType{}
	if d.interf != nil {
	}

	v := d.getValue()
	noOp1(v)
}

func notUsed_MethodCall_Nested_OK() {
	d := dummyType{}
	if d.interf != nil {
	}
	noOp1(d.getValue())
}

func notUsed_Pointer_OK() {
	d := dummyType{}
	if d.interf != nil {
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
	if d.interf == nil {
		return
	}

	for i := 0; i < t; i++ {
		noOp1(d.interf)
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
	d.interf = getValue()
}

func notUsed_ReferenceToFields_OK() {
	a, b := getTwoDummies()
	if getBool(a.interf) {
		return
	}
	noOp1(b.interf)
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

func notUsed_IndexExpression_StructField_Index_OK() interface{} {
	dummy := getDummy()

	idx := getInt()

	if idx < 0 {
		return nil
	}
	return dummy.slice[idx]
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

func notUsed_AssignmentToPointer_OK() {
	v := &dummyType{}

	if v != nil { // cannot be `if v := &dummyType{}; v != nil
		return
	}
}

func notUsed_TypeAssertion_OK() {
	v := getValue()
	if v == nil {
		noOp1(v)
	}

	w, ok := v.(*dummyType)
	if !ok {
		noOp2(w)
	}
}

func notUsed_BinaryExprInAssign_OK() {
	v1 := "v1"

	_ = "v2" + v1

	if false {
		noOp1(v1)
	}
}

func notUsed_ReferenceInSelect_OK() {
	v := 0
	if getBool(v) {
		v = 1
	}
	select {
	case <-getChan(v):
	}
}

func notUsed_AssignInSwitch_OK() {
	y := 100
	switch {
	case true:
		y = 1
	}
	if y < 5 {
	}
}

func notUsed_PassPointer_OK() {
	a := getDummyPtr()
	if a == nil {
	}
	noOp1(*a)
}

func notUsed_Labels_OK() {
	foo := true
BREAKOUT:
	if getBool() {
		foo = false
		goto BREAKOUT
	}
	if foo {
		return
	}
}

func notUsed_UnnerStruct_OK() {
	v := getInt()
	if v == 0 {
		v = 1
	}
	type Wrapper struct{ V int }
	t := []Wrapper{{v}}
	noOp1(t)
}

func notUsed_ReturnInSlice_Selector_OK(d dummyType) ([]interface{}, interface{}) {
	v := d
	if v.interf != nil {
		return nil, v.interf
	}
	return []interface{}{v.interf}, nil
}

func notUsed_Multiple_If_Statements_OK() {
	shouldRun := false
	if shouldRun {
		noOp1(0)
	}
	if !shouldRun {
		return
	}
	noOp2(0)
}

func notUsed_Also_Used_In_Else_Body_OK() {
	x := 0
	if x > 0 {
		noOp1(0)
	}

	if y := getInt(0); y > 0 {
		noOp1(y)
	} else {
		noOp1(x)
	}
}
