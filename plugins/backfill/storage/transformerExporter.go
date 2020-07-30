// This is a plugin generated to export the configured transformer initializers

package main

import (
	flap_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializers/v1_0_9"
	flip_bat_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_0_9"
	flip_eth_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_a/v1_0_9"
	flip_knc_a_v1_0_8 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_a/v1_0_8"
	flip_knc_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_a/v1_0_9"
	flip_mana_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/mana_a/v1_0_9"
	flip_tusd_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/tusd_a/v1_0_9"
	flip_usdc_a_v1_0_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a/v1_0_4"
	flip_usdc_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a/v1_0_9"
	flip_usdc_b_v1_0_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b/v1_0_7"
	flip_usdc_b_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b/v1_0_9"
	flip_wbtc_a_v1_0_6 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_a/v1_0_6"
	flip_wbtc_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_a/v1_0_9"
	flip_zrx_a_v1_0_8 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_a/v1_0_8"
	flip_zrx_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_a/v1_0_9"
	flop_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop/initializers/v1_0_9"
	median_knc "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_knc"
	median_mana "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_mana"
	median_wbtc "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_wbtc"
	median_zrx "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_zrx"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			flip_usdc_b_v1_0_9.StorageTransformerInitializer,
			median_zrx.StorageTransformerInitializer,
			flip_eth_a_v1_0_9.StorageTransformerInitializer,
			flip_usdc_a_v1_0_9.StorageTransformerInitializer,
			flip_zrx_a_v1_0_8.StorageTransformerInitializer,
			flap_v1_0_9.StorageTransformerInitializer,
			flip_knc_a_v1_0_9.StorageTransformerInitializer,
			flip_usdc_b_v1_0_7.StorageTransformerInitializer,
			flip_zrx_a_v1_0_9.StorageTransformerInitializer,
			flip_wbtc_a_v1_0_9.StorageTransformerInitializer,
			flop_v1_0_9.StorageTransformerInitializer,
			flip_bat_a_v1_0_9.StorageTransformerInitializer,
			flip_knc_a_v1_0_8.StorageTransformerInitializer,
			flip_mana_a_v1_0_9.StorageTransformerInitializer,
			flip_tusd_a_v1_0_9.StorageTransformerInitializer,
			flip_usdc_a_v1_0_4.StorageTransformerInitializer,
			flip_wbtc_a_v1_0_6.StorageTransformerInitializer,
			median_knc.StorageTransformerInitializer,
			median_mana.StorageTransformerInitializer,
			median_wbtc.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
