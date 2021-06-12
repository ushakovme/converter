package application

import "io"

type Converter interface {
	PNGToJPG(reader io.Reader, w io.Writer) error
}
