// This is a plugin generated to export the configured transformer initializers

package main

import (
	flip_univ2usdceth_a_v1_2_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2usdceth_a/v1_2_4"
	flip_univ2wbtceth_a_v1_2_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2wbtceth_a/v1_2_4"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
	[]storage.TransformerInitializer{
		flip_univ2wbtceth_a_v1_2_4.StorageTransformerInitializer,
		flip_univ2usdceth_a_v1_2_4.StorageTransformerInitializer,
	},
	[]interface1.ContractTransformerInitializer{}
}
