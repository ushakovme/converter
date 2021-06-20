package infrastructure

import (
	"bytes"
	"fmt"
	"github.com/ushakovme/converter/pkg/converter/application"
	proto "github.com/ushakovme/converter/proto/gen/go"
	"os"
)

type Server struct {
	proto.UnimplementedConverterServer
	converter application.Converter
}

func NewServer(converter application.Converter) *Server {
	return &Server{converter: converter}
}

func (s *Server) PNGToJPG(server proto.Converter_PNGToJPGServer) error {
	fmt.Println("PNGToJPG")
	req, err := server.Recv()
	if err != nil {
		return err
	}

	resp := proto.ConverterResponse{}
	resp.ImageID = "4"

	r := bytes.NewReader(req.Content)

	f, err := os.Create("test/files/img.jpg")
	if err != nil {
		return err
	}

	err = s.converter.PNGToJPG(r, f)
	if err != nil {
		return err
	}
	return server.SendAndClose(&resp)
}
