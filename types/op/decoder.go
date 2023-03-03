package op

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/unknownproto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func DefaultOpDecoder(cdc codec.ProtoCodecMarshaler) sdk.OpDecoder {
	return func(opBytes []byte) (sdk.Op, error) {
		var raw OpRaw

		// reject all unknown proto fields in the root TxRaw
		err := unknownproto.RejectUnknownFieldsStrict(opBytes, &raw, cdc.InterfaceRegistry())
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrOpDecode, err.Error())
		}

		err = cdc.Unmarshal(opBytes, &raw)
		if err != nil {
			return nil, err
		}

		var body OpBody

		// allow non-critical unknown fields in TxBody
		opBodyHasUnknownNonCriticals, err := unknownproto.RejectUnknownFields(raw.BodyBytes, &body, true, cdc.InterfaceRegistry())
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrOpDecode, err.Error())
		}

		err = cdc.Unmarshal(raw.BodyBytes, &body)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrOpDecode, err.Error())
		}

		theOp := &Op{
			Body: &body,
		}

		return &wrapper{
			cdc:                          cdc,
			op:                           theOp,
			bodyBz:                       raw.BodyBytes,
			opBodyHasUnknownNonCriticals: opBodyHasUnknownNonCriticals,
		}, nil
	}
}

// DefaultJSONOpDecoder returns a default protobuf JSON OpDecoder using the provided Marshaler.
func DefaultJSONOpDecoder(cdc codec.ProtoCodecMarshaler) sdk.OpDecoder {
	return func(opBytes []byte) (sdk.Op, error) {
		var theOp Op
		err := cdc.UnmarshalJSON(opBytes, &theOp)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}
		return &wrapper{
			op: &theOp,
		}, nil
	}
}
