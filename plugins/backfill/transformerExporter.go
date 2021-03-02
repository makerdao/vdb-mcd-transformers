// This is a plugin generated to export the configured transformer initializers

package main

import (
	auction_file "github.com/makerdao/vdb-mcd-transformers/transformers/events/auction_file/initializer"
	deny "github.com/makerdao/vdb-mcd-transformers/transformers/events/auth/deny_initializer"
	rely "github.com/makerdao/vdb-mcd-transformers/transformers/events/auth/rely_initializer"
	deal "github.com/makerdao/vdb-mcd-transformers/transformers/events/deal/initializer"
	dent "github.com/makerdao/vdb-mcd-transformers/transformers/events/dent/initializer"
	flip_file_cat "github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_file/cat/initializer"
	flip_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick/initializer"
	tend "github.com/makerdao/vdb-mcd-transformers/transformers/events/tend/initializer"
	tick "github.com/makerdao/vdb-mcd-transformers/transformers/events/tick/initializer"
	yank "github.com/makerdao/vdb-mcd-transformers/transformers/events/yank/initializer"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{
			yank.EventTransformerInitializer,
			tick.EventTransformerInitializer,
			tend.EventTransformerInitializer,
			rely.EventTransformerInitializer,
			flip_kick.EventTransformerInitializer,
			flip_file_cat.EventTransformerInitializer,
			deny.EventTransformerInitializer,
			dent.EventTransformerInitializer,
			deal.EventTransformerInitializer,
			auction_file.EventTransformerInitializer,
		},
		[]storage.TransformerInitializer{},
		[]interface1.ContractTransformerInitializer{}
}
