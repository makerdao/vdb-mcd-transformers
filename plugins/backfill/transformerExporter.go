// This is a plugin generated to export the configured transformer initializers

package main

import (
	clip_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_kick/initializer"
	clip_redo "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_redo/initializer"
	clip_take "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_take/initializer"
	clip_yank "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_yank/initializer"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{
			clip_kick.EventTransformerInitializer,
			clip_take.EventTransformerInitializer,
			clip_redo.EventTransformerInitializer,
			clip_yank.EventTransformerInitializer,
		},
		[]storage.TransformerInitializer{},
		[]interface1.ContractTransformerInitializer{}
}
