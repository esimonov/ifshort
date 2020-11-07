package testdata

func dummyNoOp1(v interface{}) {}

func dummyNoOp2(v interface{}) { dummyNoOp1(v) }

func dummyWithReturn() interface{} { return nil }

func dummyForLongCall(...interface{}) interface{} { return nil }

type dummyType struct{ v interface{} }

func (dt dummyType) dummyNoOp() {}

func (dt dummyType) dummyWithReturn() interface{} { return nil }
