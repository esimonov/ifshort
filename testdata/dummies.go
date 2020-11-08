package testdata

func noOp1(interface{}) {}

func noOp2(interface{}) {}

func returnValue() interface{} { return nil }

func returnTwoValues() (interface{}, interface{}) { return nil, nil }

func longCallWithReturnValue(...interface{}) interface{} { return nil }

type dummyType struct{ v interface{} }

func (dt dummyType) noOp() {}

func (dt dummyType) returnValue() interface{} { return nil }
