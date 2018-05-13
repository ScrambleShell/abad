package types

type (
	Bool bool
)

var True = Bool(true)
var False = Bool(false)

func NewBool(b bool) Bool {
	return Bool(b)
}

func (_ Bool) Kind() Kind {
	return KindBool
}

func (b Bool) IsTrue() bool {
	return bool(b)
}

func (b Bool) IsFalse() bool {
	return bool(b)
}

func (b Bool) ToBool() Bool {
	return b
}

func (b Bool) ToNumber() Number {
	if b {
		return NewNumber(1)
	}
	return NewNumber(+0)
}

func (b Bool) ToString() String {
	if b {
		return NewString("true")
	}
	return NewString("false")
}

func (b Bool) Equal(a Bool) bool {
	return bool(b) == bool(a) 
}