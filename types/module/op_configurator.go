package module

import (
	"github.com/gogo/protobuf/grpc"
)

type OpConfigurator interface {
	OpServer() grpc.Server
}

type opConfigurator struct {
	opServer grpc.Server
}

func NewOpConfigurator(opServer grpc.Server) OpConfigurator {
	return &opConfigurator{opServer: opServer}
}

var _ OpConfigurator = opConfigurator{}

func (c opConfigurator) OpServer() grpc.Server {
	return c.opServer
}
