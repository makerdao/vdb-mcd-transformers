// This is a plugin generated to export the configured transformer initializers

package main

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/tusd_flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a_flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b_flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_flip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_knc"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_wbtc"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_zrx"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			knc_flip.StorageTransformerInitializer,
			median_knc.StorageTransformerInitializer,
			median_wbtc.StorageTransformerInitializer,
			median_zrx.StorageTransformerInitializer,
			tusd_flip.StorageTransformerInitializer,
			usdc_a_flip.StorageTransformerInitializer,
			usdc_b_flip.StorageTransformerInitializer,
			wbtc_flip.StorageTransformerInitializer,
			zrx_flip.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
