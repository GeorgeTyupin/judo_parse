package locationsplitter

import (
	"judo/internal/interfaces"
	lsio "judo/internal/io/excel/location_splitter"
)

type LocationSplitter struct {
	Reader interfaces.Reader
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
