package rely_initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.VatRelyTable, constants.RelySignature()),
	Transformer: vat_auth.Transformer{TableName: constants.VatRelyTable},
}.NewTransformer
