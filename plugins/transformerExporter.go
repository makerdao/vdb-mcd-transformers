// This is a plugin generated to export the configured transformer initializers

package main

import (
	bite "github.com/vulcanize/mcd_transformers/transformers/events/bite/initializer"
	cat_file_chop_lump "github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump/initializer"
	cat_file_flip "github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip/initializer"
	cat_file_vow "github.com/vulcanize/mcd_transformers/transformers/events/cat_file/vow/initializer"
	deal "github.com/vulcanize/mcd_transformers/transformers/events/deal/initializer"
	dent "github.com/vulcanize/mcd_transformers/transformers/events/dent/initializer"
	flap_kick "github.com/vulcanize/mcd_transformers/transformers/events/flap_kick/initializer"
	flip_kick "github.com/vulcanize/mcd_transformers/transformers/events/flip_kick/initializer"
	flop_kick "github.com/vulcanize/mcd_transformers/transformers/events/flop_kick/initializer"
	jug_drip "github.com/vulcanize/mcd_transformers/transformers/events/jug_drip/initializer"
	jug_file_base "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/base/initializer"
	jug_file_ilk "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/ilk/initializer"
	jug_file_vow "github.com/vulcanize/mcd_transformers/transformers/events/jug_file/vow/initializer"
	jug_init "github.com/vulcanize/mcd_transformers/transformers/events/jug_init/initializer"
	new_cdp "github.com/vulcanize/mcd_transformers/transformers/events/new_cdp/initializer"
	spot_file_mat "github.com/vulcanize/mcd_transformers/transformers/events/spot_file/mat/initializer"
	spot_file_pip "github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip/initializer"
	spot_poke "github.com/vulcanize/mcd_transformers/transformers/events/spot_poke/initializer"
	tend "github.com/vulcanize/mcd_transformers/transformers/events/tend/initializer"
	tick "github.com/vulcanize/mcd_transformers/transformers/events/tick/initializer"
	vat_file_debt_ceiling "github.com/vulcanize/mcd_transformers/transformers/events/vat_file/debt_ceiling/initializer"
	vat_file_ilk "github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk/initializer"
	vat_flux "github.com/vulcanize/mcd_transformers/transformers/events/vat_flux/initializer"
	vat_fold "github.com/vulcanize/mcd_transformers/transformers/events/vat_fold/initializer"
	vat_fork "github.com/vulcanize/mcd_transformers/transformers/events/vat_fork/initializer"
	vat_frob "github.com/vulcanize/mcd_transformers/transformers/events/vat_frob/initializer"
	vat_grab "github.com/vulcanize/mcd_transformers/transformers/events/vat_grab/initializer"
	vat_heal "github.com/vulcanize/mcd_transformers/transformers/events/vat_heal/initializer"
	vat_init "github.com/vulcanize/mcd_transformers/transformers/events/vat_init/initializer"
	vat_move "github.com/vulcanize/mcd_transformers/transformers/events/vat_move/initializer"
	vat_slip "github.com/vulcanize/mcd_transformers/transformers/events/vat_slip/initializer"
	vat_suck "github.com/vulcanize/mcd_transformers/transformers/events/vat_suck/initializer"
	vow_fess "github.com/vulcanize/mcd_transformers/transformers/events/vow_fess/initializer"
	vow_file "github.com/vulcanize/mcd_transformers/transformers/events/vow_file/initializer"
	vow_flog "github.com/vulcanize/mcd_transformers/transformers/events/vow_flog/initializer"
	yank "github.com/vulcanize/mcd_transformers/transformers/events/yank/initializer"
	cat "github.com/vulcanize/mcd_transformers/transformers/storage/cat/initializer"
	cdp_manager "github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager/initializer"
	flap_storage "github.com/vulcanize/mcd_transformers/transformers/storage/flap/initializer"
	bat_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/bat_flip"
	dgd_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/dgd_flip"
	eth_flip_a "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/eth_flip_a"
	eth_flip_b "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/eth_flip_b"
	eth_flip_c "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/eth_flip_c"
	gnt_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/gnt_flip"
	omg_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/omg_flip"
	rep_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/rep_flip"
	zrx_flip "github.com/vulcanize/mcd_transformers/transformers/storage/flip/initializers/zrx_flip"
	flop_storage "github.com/vulcanize/mcd_transformers/transformers/storage/flop/initializer"
	jug "github.com/vulcanize/mcd_transformers/transformers/storage/jug/initializer"
	spot "github.com/vulcanize/mcd_transformers/transformers/storage/spot/initializer"
	vat "github.com/vulcanize/mcd_transformers/transformers/storage/vat/initializer"
	vow "github.com/vulcanize/mcd_transformers/transformers/storage/vow/initializer"
	interface1 "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]interface1.EventTransformerInitializer, []interface1.StorageTransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []interface1.EventTransformerInitializer{flap_kick.EventTransformerInitializer, flip_kick.EventTransformerInitializer, jug_file_ilk.EventTransformerInitializer, jug_file_vow.EventTransformerInitializer, yank.EventTransformerInitializer, jug_drip.EventTransformerInitializer, vat_file_debt_ceiling.EventTransformerInitializer, vat_fold.EventTransformerInitializer, vow_file.EventTransformerInitializer, vat_move.EventTransformerInitializer, flop_kick.EventTransformerInitializer, vow_flog.EventTransformerInitializer, vat_file_ilk.EventTransformerInitializer, jug_file_base.EventTransformerInitializer, tick.EventTransformerInitializer, spot_poke.EventTransformerInitializer, vat_suck.EventTransformerInitializer, vat_fork.EventTransformerInitializer, cat_file_flip.EventTransformerInitializer, spot_file_pip.EventTransformerInitializer, vat_flux.EventTransformerInitializer, deal.EventTransformerInitializer, jug_init.EventTransformerInitializer, dent.EventTransformerInitializer, vow_fess.EventTransformerInitializer, vat_heal.EventTransformerInitializer, cat_file_vow.EventTransformerInitializer, new_cdp.EventTransformerInitializer, vat_init.EventTransformerInitializer, bite.EventTransformerInitializer, cat_file_chop_lump.EventTransformerInitializer, vat_slip.EventTransformerInitializer, vat_frob.EventTransformerInitializer, vat_grab.EventTransformerInitializer, spot_file_mat.EventTransformerInitializer, tend.EventTransformerInitializer}, []interface1.StorageTransformerInitializer{zrx_flip.StorageTransformerInitializer, eth_flip_c.StorageTransformerInitializer, gnt_flip.StorageTransformerInitializer, jug.StorageTransformerInitializer, eth_flip_a.StorageTransformerInitializer, eth_flip_b.StorageTransformerInitializer, flap_storage.StorageTransformerInitializer, flop_storage.StorageTransformerInitializer, rep_flip.StorageTransformerInitializer, omg_flip.StorageTransformerInitializer, bat_flip.StorageTransformerInitializer, vow.StorageTransformerInitializer, vat.StorageTransformerInitializer, cdp_manager.StorageTransformerInitializer, cat.StorageTransformerInitializer, spot.StorageTransformerInitializer, dgd_flip.StorageTransformerInitializer}, []interface1.ContractTransformerInitializer{}
}
