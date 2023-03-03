package op

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

func DefaultOpEncoder() sdk.OpEncoder {
	return func(op sdk.Op) ([]byte, error) {
		opWrapper, ok := op.(*wrapper)
		if !ok {
			return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, op)
		}

		raw := &OpRaw{BodyBytes: opWrapper.getBodyBytes()}

		return proto.Marshal(raw)
	}
}

func DefaultJSONEncoder(cdc codec.ProtoCodecMarshaler) sdk.OpEncoder {
	return func(op sdk.Op) ([]byte, error) {
		opWrapper, ok := op.(*wrapper)
		if ok {
			return cdc.MarshalJSON(opWrapper.op)
		}

		protoOp, ok := op.(*Op)
		if ok {
			return cdc.MarshalJSON(protoOp)
		}

		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, op)
	}
}
