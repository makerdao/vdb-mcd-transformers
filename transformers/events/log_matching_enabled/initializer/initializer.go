package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_matching_enabled"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.LogMatchingEnabledTable, constants.LogMatchingEnabledSignature()),
	Transformer: log_matching_enabled.Transformer{},
}.NewTransformer
