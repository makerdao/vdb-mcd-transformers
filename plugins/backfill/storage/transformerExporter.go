// This is a plugin generated to export the configured transformer initializers

package main

import (
	flip_eth_b_v1_1_3 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_b/v1_1_3"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			flip_eth_b_v1_1_3.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
