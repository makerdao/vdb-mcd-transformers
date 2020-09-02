// This is a plugin generated to export the configured transformer initializers

package main

import (
	cat_file_box "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/box/initializer"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{
			cat_file_box.EventTransformerInitializer,
		},
		[]storage.TransformerInitializer{},
		[]interface1.ContractTransformerInitializer{}
}
