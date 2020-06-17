// This is a plugin generated to export the configured transformer initializers

package main

import (
	deny "github.com/makerdao/vdb-mcd-transformers/transformers/events/auth/deny_initializer"
	rely "github.com/makerdao/vdb-mcd-transformers/transformers/events/auth/rely_initializer"
	deal "github.com/makerdao/vdb-mcd-transformers/transformers/events/deal/initializer"
	dent "github.com/makerdao/vdb-mcd-transformers/transformers/events/dent/initializer"
	flip_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick/initializer"
	log_kill "github.com/makerdao/vdb-mcd-transformers/transformers/events/log_kill/initializer"
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
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{
			deal.EventTransformerInitializer,
			dent.EventTransformerInitializer,
			deny.EventTransformerInitializer,
			flip_kick.EventTransformerInitializer,
			log_kill.EventTransformerInitializer,
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
