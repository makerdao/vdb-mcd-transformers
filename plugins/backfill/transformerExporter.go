// This is a plugin generated to export the configured transformer initializers

package main

import (
	dog_v1_3_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog/initializers/v1_3_0"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			dog_v1_3_0.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
