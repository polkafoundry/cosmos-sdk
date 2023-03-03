package types

import "github.com/gogo/protobuf/proto"

type OpMsg interface {
	proto.Message

	ValidateBasic() error
}

// Op defines the interface an operation must fulfill.
type Op interface {
	GetMsgs() []OpMsg

	// ValidateBasic does a simple and lightweight validation check that doesn't
	// require access to any other information.
	ValidateBasic() error
}

type OpWithMemo interface {
	Op
	GetMemo() string
}

type OpWithTimeoutHeight interface {
	Op
	GetTimeoutHeight() uint64
}

// OpMsgTypeURL returns the TypeURL of a `sdk.OpMsg`.
func OpMsgTypeURL(opMsg OpMsg) string {
	return "/" + proto.MessageName(opMsg)
}

// OpDecoder unmarshals transaction bytes
type OpDecoder func(opBytes []byte) (Op, error)

// OpEncoder marshals transaction to bytes
type OpEncoder func(op Op) ([]byte, error)
