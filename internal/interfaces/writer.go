package interfaces

type Writer interface {
	Write(data any)
	SaveFile()
}
