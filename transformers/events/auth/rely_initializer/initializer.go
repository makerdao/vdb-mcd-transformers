package rely_initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.RelyTable, constants.RelySignature()),
	Transformer: auth.Transformer{TableName: constants.RelyTable},
}.NewTransformer
