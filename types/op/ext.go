package op

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
)

// ExtensionOptionI defines the interface for tx extension options
type ExtensionOptionI interface{}

// unpackExtensionOptionsI unpacks Any's to TxExtensionOptionI's.
func unpackExtensionOptionsI(unpacker types.AnyUnpacker, anys []*types.Any) error {
	for _, any := range anys {
		var opt ExtensionOptionI
		err := unpacker.UnpackAny(any, &opt)
		if err != nil {
			return err
		}
	}

	return nil
}
