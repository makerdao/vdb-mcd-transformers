// This is a plugin generated to export the configured transformer initializers

package main

import (
	clip_aave_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/aave_a/v1_6_0"
	clip_bal_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/bal_a/v1_6_0"
	clip_bat_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/bat_a/v1_6_0"
	clip_comp_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/comp_a/v1_6_0"
	clip_eth_a_v1_5_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/eth_a/v1_5_0"
	clip_eth_b_v1_5_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/eth_b/v1_5_0"
	clip_eth_c_v1_5_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/eth_c/v1_5_0"
	clip_knc_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/knc_a/v1_6_0"
	clip_link_a_v1_3_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/link_a/v1_3_0"
	clip_lrc_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/lrc_a/v1_6_0"
	clip_mana_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/mana_a/v1_6_0"
	clip_matic_a_v1_9_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/matic_a/v1_9_4"
	clip_renbtc_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/renbtc_a/v1_6_0"
	clip_uni_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/uni_a/v1_6_0"
	clip_univ2aaveeth_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2aaveeth_a/v1_8_0"
	clip_univ2daieth_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2daieth_a/v1_8_0"
	clip_univ2daiusdt_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2daiusdt_a/v1_8_0"
	clip_univ2ethusdt_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2ethusdt_a/v1_8_0"
	clip_univ2linketh_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2linketh_a/v1_8_0"
	clip_univ2unieth_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2unieth_a/v1_8_0"
	clip_univ2usdceth_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2usdceth_a/v1_8_0"
	clip_univ2wbtcdai_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2wbtcdai_a/v1_8_0"
	clip_univ2wbtceth_a_v1_8_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/univ2wbtceth_a/v1_8_0"
	clip_wbtc_a_v1_5_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/wbtc_a/v1_5_0"
	clip_wbtc_b_v1_9_10 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/wbtc_b/v1_9_10"
	clip_wbtc_c_v1_9_11 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/wbtc_c/v1_9_11"
	clip_yfi_a_v1_5_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/yfi_a/v1_5_0"
	clip_zrx_a_v1_6_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/zrx_a/v1_6_0"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			clip_aave_a_v1_6_0.StorageTransformerInitializer,
			clip_bal_a_v1_6_0.StorageTransformerInitializer,
			clip_bat_a_v1_6_0.StorageTransformerInitializer,
			clip_comp_a_v1_6_0.StorageTransformerInitializer,
			clip_eth_a_v1_5_0.StorageTransformerInitializer,
			clip_eth_b_v1_5_0.StorageTransformerInitializer,
			clip_eth_c_v1_5_0.StorageTransformerInitializer,
			clip_knc_a_v1_6_0.StorageTransformerInitializer,
			clip_link_a_v1_3_0.StorageTransformerInitializer,
			clip_lrc_a_v1_6_0.StorageTransformerInitializer,
			clip_mana_a_v1_6_0.StorageTransformerInitializer,
			clip_matic_a_v1_9_4.StorageTransformerInitializer,
			clip_renbtc_a_v1_6_0.StorageTransformerInitializer,
			clip_uni_a_v1_6_0.StorageTransformerInitializer,
			clip_univ2aaveeth_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2daieth_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2daiusdt_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2ethusdt_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2linketh_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2unieth_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2usdceth_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2wbtcdai_a_v1_8_0.StorageTransformerInitializer,
			clip_univ2wbtceth_a_v1_8_0.StorageTransformerInitializer,
			clip_wbtc_a_v1_5_0.StorageTransformerInitializer,
			clip_wbtc_b_v1_9_10.StorageTransformerInitializer,
			clip_wbtc_c_v1_9_11.StorageTransformerInitializer,
			clip_yfi_a_v1_5_0.StorageTransformerInitializer,
			clip_zrx_a_v1_6_0.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
