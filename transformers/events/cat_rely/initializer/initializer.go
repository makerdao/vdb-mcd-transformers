package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_rely"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = event.Transformer{
	Config:    shared.GetEventTransformerConfig(constants.CatRelyTable, constants.CatRelySignature()),
	Converter: cat_rely.Converter{},
}.NewTransformer
