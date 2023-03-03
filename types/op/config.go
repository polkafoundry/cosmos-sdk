package op

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type config struct {
	encoder     sdk.OpEncoder
	decoder     sdk.OpDecoder
	jsonEncoder sdk.OpEncoder
	jsonDecoder sdk.OpDecoder
	protoCodec  codec.ProtoCodecMarshaler
}

func (c config) NewTxBuilder() client.OpBuilder {
	return newBuilder(c.protoCodec)
}

func (c config) WrapTxBuilder(newOp sdk.Op) (client.OpBuilder, error) {
	newBuilder, ok := newOp.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, newOp)
	}

	return newBuilder, nil
}

func (c config) OpEncoder() sdk.OpEncoder {
	return c.encoder
}

func (c config) OpDecoder() sdk.OpDecoder {
	return c.decoder
}

func (c config) OpJSONEncoder() sdk.OpEncoder {
	return c.jsonEncoder
}

func (c config) OpJSONDecoder() sdk.OpDecoder {
	return c.jsonDecoder
}

func NewOpConfig(protoCodec codec.ProtoCodecMarshaler) client.OpConfig {
	return &config{
		encoder:     DefaultOpEncoder(),
		decoder:     DefaultOpDecoder(protoCodec),
		jsonEncoder: DefaultJSONEncoder(protoCodec),
		jsonDecoder: DefaultJSONOpDecoder(protoCodec),
		protoCodec:  protoCodec,
	}
}
