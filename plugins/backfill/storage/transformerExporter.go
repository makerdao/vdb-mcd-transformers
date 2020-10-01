// This is a plugin generated to export the configured transformer initializers

package main

import (
	flip_comp_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/comp_a/v1_1_2"
	flip_link_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/link_a/v1_1_2"
	flip_lrc_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/lrc_a/v1_1_2"
	median_comp "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_comp"
	median_link "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_link"
	median_lrc "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_lrc"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		   []storage.TransformerInitializer{
				flip_comp_a_v1_1_2.StorageTransformerInitializer,
				flip_link_a_v1_1_2.StorageTransformerInitializer,
				flip_lrc_a_v1_1_2.StorageTransformerInitializer,
				median_comp.StorageTransformerInitializer,
				median_link.StorageTransformerInitializer,
				median_lrc.StorageTransformerInitializer,
		   },
		   []interface1.ContractTransformerInitializer{}
}
