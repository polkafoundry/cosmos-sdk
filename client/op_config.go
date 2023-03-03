package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	OpEncodingConfig interface {
		OpEncoder() sdk.OpEncoder
		OpDecoder() sdk.OpDecoder
		OpJSONEncoder() sdk.OpEncoder
		OpJSONDecoder() sdk.OpDecoder
	}

	// OpConfig defines an interface a client can utilize to generate an
	// application-defined concrete operation type. The type returned must
	// implement OpBuilder.
	OpConfig interface {
		OpEncodingConfig
		NewTxBuilder() OpBuilder
		WrapTxBuilder(sdk.Op) (OpBuilder, error)
	}

	OpBuilder interface {
		GetOp() sdk.Op

		SetMsgs(msgs ...sdk.OpMsg) error
		SetMemo(memo string)
	}
)
