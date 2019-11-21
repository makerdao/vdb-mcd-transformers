// This is a plugin generated to export the configured transformer initializers

package main

import (
	bite "github.com/makerdao/vdb-mcd-transformers/transformers/events/bite/initializer"
	cat_file_chop_lump "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/chop_lump/initializer"
	cat_file_flip "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/flip/initializer"
	cat_file_vow "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/vow/initializer"
	deal "github.com/makerdao/vdb-mcd-transformers/transformers/events/deal/initializer"
	dent "github.com/makerdao/vdb-mcd-transformers/transformers/events/dent/initializer"
	flap_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flap_kick/initializer"
	flip_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick/initializer"
	flop_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick/initializer"
	jug_drip "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_drip/initializer"
	jug_file_base "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/base/initializer"
	jug_file_ilk "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/ilk/initializer"
	jug_file_vow "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/vow/initializer"
	jug_init "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_init/initializer"
	new_cdp "github.com/makerdao/vdb-mcd-transformers/transformers/events/new_cdp/initializer"
	spot_file_mat "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/mat/initializer"
	spot_file_pip "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/pip/initializer"
	spot_poke "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_poke/initializer"
	tend "github.com/makerdao/vdb-mcd-transformers/transformers/events/tend/initializer"
	tick "github.com/makerdao/vdb-mcd-transformers/transformers/events/tick/initializer"
	vat_file_debt_ceiling "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/debt_ceiling/initializer"
	vat_file_ilk "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/ilk/initializer"
	vat_flux "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_flux/initializer"
	vat_fold "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_fold/initializer"
	vat_fork "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_fork/initializer"
	vat_frob "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob/initializer"
	vat_grab "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_grab/initializer"
	vat_heal "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_heal/initializer"
	vat_init "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_init/initializer"
	vat_move "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_move/initializer"
	vat_slip "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_slip/initializer"
	vat_suck "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_suck/initializer"
	vow_fess "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_fess/initializer"
	vow_file "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/initializer"
	vow_flog "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_flog/initializer"
	yank "github.com/makerdao/vdb-mcd-transformers/transformers/events/yank/initializer"
	cat "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/initializer"
	cdp_manager "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager/initializer"
	flap_storage "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializer"
	bat_flip "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_flip"
	eth_flip_a "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_flip_a"
	sai_flip "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/sai_flip"
	flop_storage "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop/initializer"
	jug "github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug/initializer"
	spot "github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot/initializer"
	vat "github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat/initializer"
	vow "github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow/initializer"
	interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

type exporter string

var Exporter exporter

func (e exporter) Export() ([]interface1.EventTransformerInitializer, []interface1.StorageTransformerInitializer, []interface1.ContractTransformerInitializer) {
	return []interface1.EventTransformerInitializer{
			bite.EventTransformerInitializer,
			cat_file_chop_lump.EventTransformerInitializer,
			cat_file_flip.EventTransformerInitializer,
			cat_file_vow.EventTransformerInitializer,
			deal.EventTransformerInitializer,
			dent.EventTransformerInitializer,
			flap_kick.EventTransformerInitializer,
			flip_kick.EventTransformerInitializer,
			flop_kick.EventTransformerInitializer,
			jug_drip.EventTransformerInitializer,
			jug_file_base.EventTransformerInitializer,
			jug_file_ilk.EventTransformerInitializer,
			jug_init.EventTransformerInitializer,
			jug_file_vow.EventTransformerInitializer,
			new_cdp.EventTransformerInitializer,
			spot_file_mat.EventTransformerInitializer,
			spot_file_pip.EventTransformerInitializer,
			spot_poke.EventTransformerInitializer,
			tend.EventTransformerInitializer,
			tick.EventTransformerInitializer,
			vat_flux.EventTransformerInitializer,
			vat_fold.EventTransformerInitializer,
			vat_fork.EventTransformerInitializer,
			vat_frob.EventTransformerInitializer,
			vat_grab.EventTransformerInitializer,
			vat_heal.EventTransformerInitializer,
			vat_init.EventTransformerInitializer,
			vat_move.EventTransformerInitializer,
			vat_slip.EventTransformerInitializer,
			vat_suck.EventTransformerInitializer,
			vow_fess.EventTransformerInitializer,
			vow_file.EventTransformerInitializer,
			vat_file_debt_ceiling.EventTransformerInitializer,
			vat_file_ilk.EventTransformerInitializer,
			vow_flog.EventTransformerInitializer,
			yank.EventTransformerInitializer,
		},
		[]interface1.StorageTransformerInitializer{
			bat_flip.StorageTransformerInitializer,
			cat.StorageTransformerInitializer,
			cdp_manager.StorageTransformerInitializer,
			eth_flip_a.StorageTransformerInitializer,
			flap_storage.StorageTransformerInitializer,
			flop_storage.StorageTransformerInitializer,
			jug.StorageTransformerInitializer,
			sai_flip.StorageTransformerInitializer,
			spot.StorageTransformerInitializer,
			vat.StorageTransformerInitializer,
			vow.StorageTransformerInitializer,
		},
		[]interface1.ContractTransformerInitializer{}
}
