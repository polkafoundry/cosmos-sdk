package op

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

type wrapper struct {
	cdc codec.Codec

	op *Op

	// bodyBz represents the protobuf encoding of OpBody. This should be encoding
	// from the client using OpRaw if the op was decoded from the wire
	bodyBz []byte

	opBodyHasUnknownNonCriticals bool
}

var (
	_ client.OpBuilder          = &wrapper{}
	_ sdk.Op                    = &wrapper{}
	_ ExtensionOptionsOpBuilder = &wrapper{}
	_ HasExtensionOptionsOp     = &wrapper{}
)

type HasExtensionOptionsOp interface {
	GetExtensionOptions() []*codectypes.Any
	GetNonCriticalExtensionOptions() []*codectypes.Any
}

type ExtensionOptionsOpBuilder interface {
	client.OpBuilder

	SetExtensionOptions(...*codectypes.Any)
	SetNonCriticalExtensionOptions(...*codectypes.Any)
}

func newBuilder(cdc codec.Codec) *wrapper {
	return &wrapper{
		cdc: cdc,
		op: &Op{
			Body: &OpBody{},
		},
	}
}

func (w *wrapper) GetMsgs() []sdk.OpMsg {
	return w.op.GetMsgs()
}

func (w *wrapper) ValidateBasic() error {
	return w.op.ValidateBasic()
}

func (w *wrapper) getBodyBytes() []byte {
	if len(w.bodyBz) == 0 {
		// if bodyBz is empty, then marshal the body. bodyBz will generally
		// be set to nil whenever SetBody is called so the result of calling
		// this method should always return the correct bytes. Note that after
		// decoding bodyBz is derived from TxRaw so that it matches what was
		// transmitted over the wire
		var err error
		w.bodyBz, err = proto.Marshal(w.op.Body)
		if err != nil {
			panic(err)
		}
	}
	return w.bodyBz
}

func (w *wrapper) GetOp() sdk.Op {
	return w
}

func (w *wrapper) GetProtoOp() *Op {
	return w.op
}

func (w *wrapper) SetMsgs(msgs ...sdk.OpMsg) error {
	anys, err := SetMsgs(msgs)
	if err != nil {
		return err
	}

	w.op.Body.Messages = anys

	// set bodyBz to nil because the cached bodyBz no longer matches op.Body
	w.bodyBz = nil

	return nil
}

func (w *wrapper) SetMemo(memo string) {
	w.op.Body.Memo = memo

	// set bodyBz to nil because the cached bodyBz no longer matches op.Body
	w.bodyBz = nil
}

func WrapOp(protoOp *Op) client.OpBuilder {
	return &wrapper{
		op: protoOp,
	}
}

func (w *wrapper) SetExtensionOptions(extOpts ...*codectypes.Any) {
	w.op.Body.ExtensionOptions = extOpts
	w.bodyBz = nil
}

func (w *wrapper) SetNonCriticalExtensionOptions(extOpts ...*codectypes.Any) {
	w.op.Body.NonCriticalExtensionOptions = extOpts
	w.bodyBz = nil
}

func (w *wrapper) GetExtensionOptions() []*codectypes.Any {
	return w.op.Body.ExtensionOptions
}

func (w *wrapper) GetNonCriticalExtensionOptions() []*codectypes.Any {
	return w.op.Body.NonCriticalExtensionOptions
}
