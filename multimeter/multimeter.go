package multimeter

type Multimeter interface {
	ProccessArray(bytearray []byte, printArray bool) (float64, string, []string)
}
