package testdata

func noOp1(interface{}) {}

func noOp2(interface{}) {}

func getValue(...interface{}) interface{} { return nil }

func getTwoValues(...interface{}) (interface{}, interface{}) { return nil, nil }

func getBool(...interface{}) bool { return false }

type dummyType struct{ v interface{} }

func getDummy() dummyType { return dummyType{} }

func getTwoDummies() (dummyType, dummyType) { return dummyType{}, dummyType{} }

func (dt dummyType) noOp() {}

func (dt dummyType) getValue() interface{} { return nil }
