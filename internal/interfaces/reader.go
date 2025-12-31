package interfaces

type Reader interface {
	Read() (any, error)
}
