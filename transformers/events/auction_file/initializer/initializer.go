package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auction_file"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.AuctionFileTable, constants.AuctionFileSignature()),
	Transformer: auction_file.Transformer{},
}.NewTransformer
