package model

// Foo obj
type Foo struct {
	// Name of foo
	// @gen foo name
	Name string `json:"tag"`
}

type Bar struct {
	Bool bool
	Byte byte
	I1   int
	I2   int8
	I3   int16
	I4   int32
	I5   int64
	UI1  uint
	UI2  uint8
	UI3  uint16
	UI4  uint32
	UI5  uint64
	F1   float32
	F2   float64
	C1   complex64
	C2   complex128
	S    string
	R    rune
}
