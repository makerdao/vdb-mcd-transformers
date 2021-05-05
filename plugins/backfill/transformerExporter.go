// This is a plugin generated to export the configured transformer initializers

package main

import (
	clip_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_kick/initializer"
	clip_redo "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_redo/initializer"
	clip_take "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_take/initializer"
	clip_yank "github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_yank/initializer"
	dog_deny "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_deny/initializer"
	dog_digs "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_digs/initializer"
	dog_file_hole "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/hole/initializer"
	dog_file_ilk_chop_hole "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_chop_hole/initializer"
	dog_file_ilk_clip "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_clip/initializer"
	dog_file_vow "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/vow/initializer"
	dog_rely "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_rely/initializer"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{
			clip_kick.EventTransformerInitializer,
			clip_redo.EventTransformerInitializer,
			clip_take.EventTransformerInitializer,
			clip_yank.EventTransformerInitializer,
			dog_deny.EventTransformerInitializer,
			dog_digs.EventTransformerInitializer,
			dog_rely.EventTransformerInitializer,
			dog_file_ilk_clip.EventTransformerInitializer,
			dog_file_hole.EventTransformerInitializer,
			dog_file_ilk_chop_hole.EventTransformerInitializer,
			dog_file_vow.EventTransformerInitializer,
		},
		[]storage.TransformerInitializer{},
		[]interface1.ContractTransformerInitializer{}
}
