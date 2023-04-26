package multimeter

type Multimeter interface {
	ProccessArray(bytearray []byte, printArray bool) (value float64, unit string, flags []string)
}
