package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/osm_change"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.OsmChangeTable, constants.OsmChangeSignature()),
	Transformer: osm_change.Transformer{},
}.NewTransformer
