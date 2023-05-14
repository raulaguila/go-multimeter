package multimeter

type Multimeter interface {
	ProccessArray(bytearray []byte, printArray bool) (value float64, unit string, flags []string)
}

type MultimeterButtons interface {
	Select() []byte
	Auto() []byte
	Range() []byte
	Led() []byte
	Relative() []byte
}
