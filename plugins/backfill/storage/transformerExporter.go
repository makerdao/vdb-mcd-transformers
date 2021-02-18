// This is a plugin generated to export the configured transformer initializers

package main

import (
	flip_univ2daiusdc_a_v1_2_5 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2daiusdc_a/v1_2_5"
	flip_univ2ethusdt_a_v1_2_5 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2ethusdt_a/v1_2_5"
	flip_univ2linketh_a_v1_2_6 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2linketh_a/v1_2_6"
	flip_univ2unieth_a_v1_2_6 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2unieth_a/v1_2_6"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			flip_univ2daiusdc_a_v1_2_5.StorageTransformerInitializer,
			flip_univ2ethusdt_a_v1_2_5.StorageTransformerInitializer,
			flip_univ2linketh_a_v1_2_6.StorageTransformerInitializer,
			flip_univ2unieth_a_v1_2_6.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
