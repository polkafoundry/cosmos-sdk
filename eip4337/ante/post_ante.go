package ante

import (
	"encoding/hex"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/op"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	tmtypes "github.com/tendermint/tendermint/types"
)

type UserOperationEmitEventDecorator struct {
	registry codectypes.InterfaceRegistry
}

func NewUserOperationEmitEventDecorator() UserOperationEmitEventDecorator {
	return UserOperationEmitEventDecorator{}
}

func (decorator UserOperationEmitEventDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
	if ok {
		opts := txWithExtensions.GetNonCriticalExtensionOptions()
		for _, opt := range opts {
			switch typeURL := opt.GetTypeUrl(); typeURL {
			case "/cosmos.op.v1beta1.ExtensionOptionsBundlingTx":
				var unpackedOpt op.ExtensionOptionsBundlingTx
				err := decorator.registry.UnpackAny(opt, &unpackedOpt)
				if err != nil {
					// FIXME(phinc275): reject or just leave it?
					continue
				}
				for _, o := range unpackedOpt.Ops {
					ctx.EventManager().EmitEvent(sdk.NewEvent(
						tmtypes.EventUserOperation,
						sdk.NewAttribute(tmtypes.EventKeyUserOperation, Encode(o)),
					))
				}
			}
		}
	}
	return next(ctx, tx, simulate)
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}
