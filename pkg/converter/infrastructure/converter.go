package infrastructure

import (
	"github.com/pkg/errors"
	"github.com/ushakovme/converter/pkg/converter/application"
	"image/jpeg"
	"image/png"
	"io"
)

type converter struct {
}

func NewConverter() application.Converter {
	return &converter{}
}

func (c *converter) PNGToJPG(reader io.Reader, w io.Writer) error {
	img, err := png.Decode(reader)
	if err != nil {
		return errors.WithStack(err)
	}
	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
