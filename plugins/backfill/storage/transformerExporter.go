// This is a plugin generated to export the configured transformer initializers

package main

import (
	cat_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_1_0/initializer"
	flip_bat_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_1_0"
	flip_eth_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_a/v1_1_0"
	flip_knc_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_a/v1_1_0"
	flip_mana_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/mana_a/v1_1_0"
	flip_paxusd_a_v1_1_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/paxusd_a/v1_1_1"
	flip_tusd_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/tusd_a/v1_1_0"
	flip_usdc_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a/v1_1_0"
	flip_usdc_b_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b/v1_1_0"
	flip_usdt_a_v1_1_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdt_a/v1_1_1"
	flip_wbtc_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_a/v1_1_0"
	flip_zrx_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_a/v1_1_0"
	median_usdt "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_usdt"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			cat_v1_1_0.StorageTransformerInitializer,
			flip_bat_a_v1_1_0.StorageTransformerInitializer,
			flip_eth_a_v1_1_0.StorageTransformerInitializer,
			flip_knc_a_v1_1_0.StorageTransformerInitializer,
			flip_mana_a_v1_1_0.StorageTransformerInitializer,
			flip_paxusd_a_v1_1_1.StorageTransformerInitializer,
			flip_tusd_a_v1_1_0.StorageTransformerInitializer,
			flip_usdc_a_v1_1_0.StorageTransformerInitializer,
			flip_usdc_b_v1_1_0.StorageTransformerInitializer,
			flip_usdt_a_v1_1_1.StorageTransformerInitializer,
			flip_wbtc_a_v1_1_0.StorageTransformerInitializer,
			flip_zrx_a_v1_1_0.StorageTransformerInitializer,
			median_usdt.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
