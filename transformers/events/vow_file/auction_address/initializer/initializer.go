package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_address"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.VowFileAuctionAddressTable, constants.VowFileAuctionAddressSignature()),
	Transformer: auction_address.Transformer{},
}.NewTransformer
