package op

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_, _ codectypes.UnpackInterfacesMessage = &Op{}, &OpBody{}
	_    sdk.Op                             = &Op{}
)

// GetMsgs implements the GetMsgs method on sdk.Tx.
func (m *Op) GetMsgs() []sdk.OpMsg {
	if m == nil || m.Body == nil {
		return nil
	}

	anys := m.Body.Messages
	res, err := GetMsgs(anys, "operation")
	if err != nil {
		panic(err)
	}
	return res
}

func (m *Op) ValidateBasic() error {
	if m == nil {
		return fmt.Errorf("bad Op")
	}

	body := m.Body
	if body == nil {
		return fmt.Errorf("missing OpBody")
	}

	return nil
}

func (m *Op) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if m.Body != nil {
		if err := m.Body.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}

	return nil
}

func (m *OpBody) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := UnpackInterfaces(unpacker, m.Messages); err != nil {
		return err
	}

	if err := unpackExtensionOptionsI(unpacker, m.ExtensionOptions); err != nil {
		return err
	}

	if err := unpackExtensionOptionsI(unpacker, m.NonCriticalExtensionOptions); err != nil {
		return err
	}

	return nil
}
