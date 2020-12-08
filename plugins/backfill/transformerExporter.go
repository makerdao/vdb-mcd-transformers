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
	log_median_price "github.com/makerdao/vdb-mcd-transformers/transformers/events/log_median_price/initializer"
	log_value "github.com/makerdao/vdb-mcd-transformers/transformers/events/log_value/initializer"
	median_diss_batch "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/batch/initializer"
	median_diss_single "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/single/initializer"
	median_drop "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_drop/initializer"
	median_kiss_batch "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_kiss/batch/initializer"
	median_kiss_single "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_kiss/single/initializer"
	median_lift "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_lift/initializer"
	osm_change "github.com/makerdao/vdb-mcd-transformers/transformers/events/osm_change/initializer"
	tend "github.com/makerdao/vdb-mcd-transformers/transformers/events/tend/initializer"
	tick "github.com/makerdao/vdb-mcd-transformers/transformers/events/tick/initializer"
	yank "github.com/makerdao/vdb-mcd-transformers/transformers/events/yank/initializer"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{
			auction_file.EventTransformerInitializer,
			deal.EventTransformerInitializer,
			dent.EventTransformerInitializer,
			deny.EventTransformerInitializer,
			flip_file_cat.EventTransformerInitializer,
			flip_kick.EventTransformerInitializer,
			log_median_price.EventTransformerInitializer,
			log_value.EventTransformerInitializer,
			median_diss_batch.EventTransformerInitializer,
			median_diss_single.EventTransformerInitializer,
			median_drop.EventTransformerInitializer,
			median_kiss_batch.EventTransformerInitializer,
			median_kiss_single.EventTransformerInitializer,
			median_lift.EventTransformerInitializer,
			osm_change.EventTransformerInitializer,
			rely.EventTransformerInitializer,
			tend.EventTransformerInitializer,
			tick.EventTransformerInitializer,
			yank.EventTransformerInitializer,
		},
		[]storage.TransformerInitializer{},
		[]interface1.ContractTransformerInitializer{}
}
