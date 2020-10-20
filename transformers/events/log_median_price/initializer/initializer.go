package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_median_price"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.LogMedianPriceTable, constants.LogMedianPriceSignature()),
	Transformer: log_median_price.Transformer{},
}.NewTransformer
