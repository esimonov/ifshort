package testdata

func noOp1(interface{}) {}

func noOp2(interface{}) {}

func returnValue() interface{} { return nil }

func returnTwoValues() (interface{}, interface{}) { return nil, nil }

func callWithOneArgAndReturnValue(interface{}) interface{} { return nil }

func callWithVariadicArgsAndReturnValue(...interface{}) interface{} { return nil }

type dummyType struct{ v interface{} }

func returnDummy() dummyType { return dummyType{} }

func returnTwoDummies() (dummyType, dummyType) { return dummyType{}, dummyType{} }

func (dt dummyType) noOp() {}

func (dt dummyType) returnValue() interface{} { return nil }
