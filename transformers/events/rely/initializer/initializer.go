package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/rely"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = event.Transformer{
	Config:    shared.GetEventTransformerConfig(constants.RelyTable, constants.RelySignature()),
	Converter: rely.Converter{},
}.NewTransformer
