package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_redo"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.ClipRedoTable, constants.ClipRedoSignature()),
	Transformer: clip_redo.Transformer{},
}.NewTransformer
