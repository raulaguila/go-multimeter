package enum

type mode int8

const (
	Reader mode = iota
	Writer
)

func (s mode) String() string {
	switch s {
	case Reader:
		return "reader"
	case Writer:
		return "writer"
	}
	return "unknown"
}
