package testdata

func noOp1(...interface{}) {}

func noOp2(...interface{}) {}

func getValue(...interface{}) interface{} { return nil }

func getTwoValues(...interface{}) (interface{}, interface{}) { return nil, nil }

func getBool(...interface{}) bool { return false }

func getInt(...interface{}) int { return 0 }

func getChan(...interface{}) chan interface{} { return nil }

type dummyType struct {
	interf interface{}
	slice  []interface{}
}

func getDummy(...interface{}) dummyType { return dummyType{} }

func getDummyPtr(...interface{}) *dummyType { return &dummyType{} }

func getTwoDummies(...interface{}) (dummyType, dummyType) { return dummyType{}, dummyType{} }

func (dt dummyType) noOp(...interface{}) {}

func (dt dummyType) getValue(...interface{}) interface{} { return nil }
