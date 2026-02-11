package locationsplitter

import (
	lsio "judo/internal/io/excel/location_splitter"
)

type Reader interface {
	Read() (any, error)
}

type LocationSplitter struct {
	Reader Reader
}

func NewLocationSplitter() (*LocationSplitter, error) {
	reader, err := lsio.NewReader("Codes")
	if err != nil {
		return nil, err
	}

	splitter := &LocationSplitter{
		Reader: reader,
	}
	return splitter, nil
}

func (ls *LocationSplitter) Split() {

}
