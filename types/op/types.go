package op

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	msgResponseInterfaceProtoName = "cosmos.op.v1beta1.MsgResponse"
)

// MsgResponse is the interface all Msg server handlers' response types need to
// implement. It's the interface that's representing all Msg responses packed
// in Anys.
type MsgResponse interface{}

// SetMsgs takes a slice of sdk.OpMsg's and turn them into Any's.
func SetMsgs(msgs []sdk.OpMsg) ([]*types.Any, error) {
	anys := make([]*types.Any, len(msgs))
	for i, msg := range msgs {
		var err error
		anys[i], err = types.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
	}
	return anys, nil
}

// GetMsgs takes a slice of Any's and turn them into sdk.OpMsg's.
func GetMsgs(anys []*types.Any, name string) ([]sdk.OpMsg, error) {
	msgs := make([]sdk.OpMsg, len(anys))
	for i, any := range anys {
		cached := any.GetCachedValue()
		if cached == nil {
			return nil, fmt.Errorf("any cached value is nil, %s messages must be correctly packed any values", name)
		}
		msgs[i] = cached.(sdk.OpMsg)
	}
	return msgs, nil
}

// UnpackInterfaces unpacks Any's to sdk.Msg's.
func UnpackInterfaces(unpacker types.AnyUnpacker, anys []*types.Any) error {
	for _, any := range anys {
		var msg sdk.OpMsg
		err := unpacker.UnpackAny(any, &msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(msgResponseInterfaceProtoName, (*MsgResponse)(nil))

	registry.RegisterInterface("cosmos.op.v1beta1.Op", (*sdk.Op)(nil))
	registry.RegisterImplementations((*sdk.Op)(nil), &Op{})

	registry.RegisterInterface("cosmos.op.v1beta1.ExtensionOptionI", (*ExtensionOptionI)(nil))
}
