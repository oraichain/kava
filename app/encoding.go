package app

import (
	enccodec "github.com/evmos/ethermint/encoding/codec"

	"github.com/kava-labs/kava/app/params"
)

// MakeEncodingConfig creates an EncodingConfig and registers the app's types on it.
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	enccodec.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
