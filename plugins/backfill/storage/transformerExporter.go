// This is a plugin generated to export the configured transformer initializers

package main

import (
	flip_univ2aaveeth_a_v1_2_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2aaveeth_a/v1_2_7"
	flip_univ2wbtcdai_a_v1_2_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2wbtcdai_a/v1_2_7"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
				flip_univ2aaveeth_a_v1_2_7.StorageTransformerInitializer,
				flip_univ2wbtcdai_a_v1_2_7.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
