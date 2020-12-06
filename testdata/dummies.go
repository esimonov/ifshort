package testdata

func noOp1(interface{}) {}

func noOp2(interface{}) {}

func returnValue() interface{} { return nil }

func returnTwoValues() (interface{}, interface{}) { return nil, nil }

func callWithOneArgAndReturn(interface{}) interface{} { return nil }

func callWithVariadicArgsAndReturn(...interface{}) interface{} { return nil }

type dummyType struct{ v interface{} }

func returnDummyType() dummyType { return dummyType{} }

func (dt dummyType) noOp() {}

func (dt dummyType) returnValue() interface{} { return nil }
