// This is a plugin generated to export the configured transformer initializers

package main

import (
	flip_aave_a_v1_2_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/aave_a/v1_2_2"
	flip_bal_a_v1_1_14 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bal_a/v1_1_14"
	flip_eth_b_v1_1_3 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_b/v1_1_3"
	flip_gusd_a_v1_1_5 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/gusd_a/v1_1_5"
	flip_renbtc_a_v1_2_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/renbtc_a/v1_2_1"
	flip_uni_a_v1_2_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/uni_a/v1_2_1"
	flip_univ2daieth_a_v1_2_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2daieth_a/v1_2_2"
	flip_yfi_a_v1_1_14 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/yfi_a/v1_1_14"
	median_aave_v1_2_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_aave/v1_2_2"
	median_bal_v1_1_14 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_bal/v1_1_14"
	median_uni_v1_2_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_uni/v1_2_1"
	median_yfi_v1_1_14 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers/median_yfi/v1_1_14"
	event "github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	storage "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]event.TransformerInitializer, []storage.TransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []event.TransformerInitializer{},
		[]storage.TransformerInitializer{
			flip_aave_a_v1_2_2.StorageTransformerInitializer,
			flip_bal_a_v1_1_14.StorageTransformerInitializer,
			flip_eth_b_v1_1_3.StorageTransformerInitializer,
			flip_gusd_a_v1_1_5.StorageTransformerInitializer,
			flip_renbtc_a_v1_2_1.StorageTransformerInitializer,
			flip_uni_a_v1_2_1.StorageTransformerInitializer,
			flip_univ2daieth_a_v1_2_2.StorageTransformerInitializer,
			flip_yfi_a_v1_1_14.StorageTransformerInitializer,
			median_aave_v1_2_2.StorageTransformerInitializer,
			median_bal_v1_1_14.StorageTransformerInitializer,
			median_uni_v1_2_1.StorageTransformerInitializer,
			median_yfi_v1_1_14.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
