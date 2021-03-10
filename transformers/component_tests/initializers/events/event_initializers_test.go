package events

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	auction_file "github.com/makerdao/vdb-mcd-transformers/transformers/events/auction_file/initializer"
	deny "github.com/makerdao/vdb-mcd-transformers/transformers/events/auth/deny_initializer"
	rely "github.com/makerdao/vdb-mcd-transformers/transformers/events/auth/rely_initializer"
	bite "github.com/makerdao/vdb-mcd-transformers/transformers/events/bite/initializer"
	cat_claw "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_claw/initializer"
	cat_file_box "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/box/initializer"
	cat_file_chop_lump_dunk "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/chop_lump_dunk/initializer"
	cat_file_flip "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/flip/initializer"
	cat_file_vow "github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/vow/initializer"
	deal "github.com/makerdao/vdb-mcd-transformers/transformers/events/deal/initializer"
	dent "github.com/makerdao/vdb-mcd-transformers/transformers/events/dent/initializer"
	dog_bark "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_bark/initializer"
	dog_deny "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_deny/initializer"
	dog_digs "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_digs/initializer"
	dog_file_ilk_clip "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_clip/initializer"
	dog_file_ilk_uint "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_uint/initializer"
	dog_rely "github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_rely/initializer"
	flap_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flap_kick/initializer"
	flip_file_cat "github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_file/cat/initializer"
	flip_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flip_kick/initializer"
	flop_kick "github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick/initializer"
	jug_drip "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_drip/initializer"
	jug_file_base "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/base/initializer"
	jug_file_ilk "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/ilk/initializer"
	jug_file_vow "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_file/vow/initializer"
	jug_init "github.com/makerdao/vdb-mcd-transformers/transformers/events/jug_init/initializer"
	log_median_price "github.com/makerdao/vdb-mcd-transformers/transformers/events/log_median_price/initializer"
	log_value "github.com/makerdao/vdb-mcd-transformers/transformers/events/log_value/initializer"
	median_diss_batch "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/batch/initializer"
	median_diss_single "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/single/initializer"
	median_drop "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_drop/initializer"
	median_kiss_batch "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_kiss/batch/initializer"
	median_kiss_single "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_kiss/single/initializer"
	median_lift "github.com/makerdao/vdb-mcd-transformers/transformers/events/median_lift/initializer"
	new_cdp "github.com/makerdao/vdb-mcd-transformers/transformers/events/new_cdp/initializer"
	osm_change "github.com/makerdao/vdb-mcd-transformers/transformers/events/osm_change/initializer"
	pot_cage "github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_cage/initializer"
	pot_drip "github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_drip/initializer"
	pot_exit "github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_exit/initializer"
	pot_file_dsr "github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/dsr/initializer"
	pot_file_vow "github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/vow/initializer"
	pot_join "github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_join/initializer"
	spot_file_mat "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/mat/initializer"
	spot_file_par "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/par/initializer"
	spot_file_pip "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/pip/initializer"
	spot_poke "github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_poke/initializer"
	tend "github.com/makerdao/vdb-mcd-transformers/transformers/events/tend/initializer"
	tick "github.com/makerdao/vdb-mcd-transformers/transformers/events/tick/initializer"
	vat_deny "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth/deny_initializer"
	vat_hope "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth/hope_initializer"
	vat_nope "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth/nope_initializer"
	vat_rely "github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth/rely_initializer"
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
	vow_file_auction_address "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_address/initializer"
	vow_file_auction_attributes "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_attributes/initializer"
	vow_flog "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_flog/initializer"
	vow_heal "github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_heal/initializer"
	yank "github.com/makerdao/vdb-mcd-transformers/transformers/events/yank/initializer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event transformer initializers", func() {
	// NOTE: using raw values in these tests instead of variables (e.g. "auction_file" instead of
	// constants.AuctionFileTable) for transformer names and topics. Doing so to reduce the risk of copy-paste errors
	// (e.g. mis-configured transformer that passes tests because the incorrect var is used in both the src and tests).
	var db = test_config.NewTestDB(test_config.NewTestNode())

	It("configures auction file", func() {
		transformer := auction_file.EventTransformerInitializer(db)
		name := "auction_file"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928163, auctionsContracts())
	})

	It("configures bite", func() {
		transformer := bite.EventTransformerInitializer(db)
		topic := "0xa716da86bc1fb6d43d1493373f34d7a418b619681cd7b90f7ea667ba1489be28"
		assertCorrectInit(transformer, "bite", topic, 8928165, catContracts())
	})

	It("configures cat claw", func() {
		transformer := cat_claw.EventTransformerInitializer(db)
		name := "cat_claw"
		topic := "0xe66d279b00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 10742907, []string{test_data.Cat110Address()})
	})

	It("configures cat file box", func() {
		transformer := cat_file_box.EventTransformerInitializer(db)
		name := "cat_file_box"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 10742907, []string{test_data.Cat110Address()})
	})

	It("configures cat file chop lump dunk", func() {
		transformer := cat_file_chop_lump_dunk.EventTransformerInitializer(db)
		name := "cat_file_chop_lump_dunk"
		topic := "0x1a0b287e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928165, catContracts())
	})

	It("configures cat file flip", func() {
		transformer := cat_file_flip.EventTransformerInitializer(db)
		name := "cat_file_flip"
		topic := "0xebecb39d00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928165, catContracts())
	})

	It("configures cat file vow", func() {
		transformer := cat_file_vow.EventTransformerInitializer(db)
		name := "cat_file_vow"
		topic := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928165, catContracts())
	})

	It("configures deal", func() {
		transformer := deal.EventTransformerInitializer(db)
		topic := "0xc959c42b00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, "deal", topic, 8928163, auctionsContracts())
	})

	It("configures deny", func() {
		transformer := deny.EventTransformerInitializer(db)
		topic := "0x9c52a7f100000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, "deny", topic, 8925057, authContracts())
	})

	It("configures dent", func() {
		transformer := dent.EventTransformerInitializer(db)
		topic := "0x5ff3a38200000000000000000000000000000000000000000000000000000000"
		addresses := append(flipContracts(), flopContracts()...)
		assertCorrectInit(transformer, "dent", topic, 8928180, addresses)
	})

	It("configures dog bark", func() {
		transformer := dog_bark.EventTransformerInitializer(db)
		topic := "0x85258d09e1e4ef299ff3fc11e74af99563f022d21f3f940db982229dc2a3358c"
		addresses := []string{test_data.Dog1xxAddress()}
		assertCorrectInit(transformer, "dog_bark", topic, 100000000000000, addresses)
	})

	It("configures dog deny", func() {
		transformer := dog_deny.EventTransformerInitializer(db)
		topic := "0x9c52a7f1fa525c4063d456584f4b5939a57b269dace468175e052aa4c179ce7a"
		addresses := []string{test_data.Dog1xxAddress()}
		assertCorrectInit(transformer, "dog_deny", topic, 100000000000000, addresses)
	})

	It("configures dog dig", func() {
		transformer := dog_digs.EventTransformerInitializer(db)
		topic := "0x54f095dc7308776bf01e8580e4dd40fd959ea4bf50b069975768320ef8d77d8a"
		addresses := []string{test_data.Dog1xxAddress()}
		assertCorrectInit(transformer, "dog_digs", topic, 100000000000000, addresses)
	})

	It("configures dog_file_ilk_clip", func() {
		transformer := dog_file_ilk_clip.EventTransformerInitializer(db)
		topic := "0x1a0b287e7eb69c42f52dc88cb0bc5f2ecb5122b0b35c2c4b755d0eaf811ae0f8"
		addresses := []string{test_data.Dog1xxAddress()}
		assertCorrectInit(transformer, "dog_file_ilk_clip", topic, 100000000000000, addresses)
	})

	It("configures dog rely", func() {
		transformer := dog_rely.EventTransformerInitializer(db)
		topic := "0x65fae35ed06235c67d3076f28ca18323d5f077aaa8c2b759b78287ec32e69afd"
		addresses := []string{test_data.Dog1xxAddress()}
		assertCorrectInit(transformer, "dog_rely", topic, 100000000000000, addresses)
	})

	It("configures dog_file_ilk_uint", func() {
		transformer := dog_file_ilk_uint.EventTransformerInitializer(db)
		topic := "0x1a0b287e7eb69c42f52dc88cb0bc5f2ecb5122b0b35c2c4b755d0eaf811ae0f8"
		addresses := []string{test_data.Dog1xxAddress()}
		assertCorrectInit(transformer, "dog_file_ilk_uint", topic, 100000000000000, addresses)
	})

	It("configures flap kick", func() {
		transformer := flap_kick.EventTransformerInitializer(db)
		name := "flap_kick"
		topic := "0xe6dde59cbc017becba89714a037778d234a84ce7f0a137487142a007e580d609"
		assertCorrectInit(transformer, name, topic, 8928163, flapContracts())
	})

	It("configures flip file cat", func() {
		transformer := flip_file_cat.EventTransformerInitializer(db)
		name := "flip_file_cat"
		topic := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 10743070, test_data.FlipV110Addresses())
	})

	It("configures flip kick", func() {
		transformer := flip_kick.EventTransformerInitializer(db)
		name := "flip_kick"
		topic := "0xc84ce3a1172f0dec3173f04caaa6005151a4bfe40d4c9f3ea28dba5f719b2a7a"
		assertCorrectInit(transformer, name, topic, 8928180, flipContracts())
	})

	It("configures flop kick", func() {
		transformer := flop_kick.EventTransformerInitializer(db)
		name := "flop_kick"
		topic := "0x7e8881001566f9f89aedb9c5dc3d856a2b81e5235a8196413ed484be91cc0df6"
		assertCorrectInit(transformer, name, topic, 9006717, flopContracts())
	})

	It("configures jug drip", func() {
		transformer := jug_drip.EventTransformerInitializer(db)
		name := "jug_drip"
		topic := "0x44e2a5a800000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.JugAddress()})
	})

	It("configures jug file base", func() {
		transformer := jug_file_base.EventTransformerInitializer(db)
		name := "jug_file_base"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.JugAddress()})
	})

	It("configures jug file ilk", func() {
		transformer := jug_file_ilk.EventTransformerInitializer(db)
		name := "jug_file_ilk"
		topic := "0x1a0b287e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.JugAddress()})
	})

	It("configures jug file vow", func() {
		transformer := jug_file_vow.EventTransformerInitializer(db)
		name := "jug_file_vow"
		topic := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.JugAddress()})
	})

	It("configures jug init", func() {
		transformer := jug_init.EventTransformerInitializer(db)
		name := "jug_init"
		topic := "0x3b66319500000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.JugAddress()})
	})

	It("configures log median price", func() {
		transformer := log_median_price.EventTransformerInitializer(db)
		name := "log_median_price"
		topic := "0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures log value", func() {
		transformer := log_value.EventTransformerInitializer(db)
		name := "log_value"
		topic := "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"
		assertCorrectInit(transformer, name, topic, 8925094, test_data.OsmAddresses())
	})

	It("configures median diss batch", func() {
		transformer := median_diss_batch.EventTransformerInitializer(db)
		name := "median_diss_batch"
		topic := "0x46d4577d00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures median diss single", func() {
		transformer := median_diss_single.EventTransformerInitializer(db)
		name := "median_diss_single"
		topic := "0x65c4ce7a00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures median drop", func() {
		transformer := median_drop.EventTransformerInitializer(db)
		name := "median_drop"
		topic := "0x8ef5eaf000000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures median kiss batch", func() {
		transformer := median_kiss_batch.EventTransformerInitializer(db)
		name := "median_kiss_batch"
		topic := "0x1b25b65f00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures median kiss single", func() {
		transformer := median_kiss_single.EventTransformerInitializer(db)
		name := "median_kiss_single"
		topic := "0xf29c29c400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures median lift", func() {
		transformer := median_lift.EventTransformerInitializer(db)
		name := "median_lift"
		topic := "0x9431810600000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925057, test_data.MedianAddresses())
	})

	It("configures new cdp", func() {
		transformer := new_cdp.EventTransformerInitializer(db)
		name := "new_cdp"
		topic := "0xd6be0bc178658a382ff4f91c8c68b542aa6b71685b8fe427966b87745c3ea7a2"
		assertCorrectInit(transformer, name, topic, 8928198, []string{test_data.CdpManagerAddress()})
	})

	It("configures osm change", func() {
		transformer := osm_change.EventTransformerInitializer(db)
		name := "osm_change"
		topic := "0x1e77933e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8925094, test_data.OsmAddresses())
	})

	It("configures pot cage", func() {
		transformer := pot_cage.EventTransformerInitializer(db)
		name := "pot_cage"
		topic := "0x6924500900000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.PotAddress()})
	})

	It("configures pot drip", func() {
		transformer := pot_drip.EventTransformerInitializer(db)
		name := "pot_drip"
		topic := "0x9f678cca00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.PotAddress()})
	})

	It("configures pot exit", func() {
		transformer := pot_exit.EventTransformerInitializer(db)
		name := "pot_exit"
		topic := "0x7f8661a100000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.PotAddress()})
	})

	It("configures pot file dsr", func() {
		transformer := pot_file_dsr.EventTransformerInitializer(db)
		name := "pot_file_dsr"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.PotAddress()})
	})

	It("configures pot file vow", func() {
		transformer := pot_file_vow.EventTransformerInitializer(db)
		name := "pot_file_vow"
		topic := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.PotAddress()})
	})

	It("configures pot join", func() {
		transformer := pot_join.EventTransformerInitializer(db)
		name := "pot_join"
		topic := "0x049878f300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928160, []string{test_data.PotAddress()})
	})

	It("configures rely", func() {
		transformer := rely.EventTransformerInitializer(db)
		topic := "0x65fae35e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, "rely", topic, 8925057, authContracts())
	})

	It("configures spot file mat", func() {
		transformer := spot_file_mat.EventTransformerInitializer(db)
		name := "spot_file_mat"
		topic := "0x1a0b287e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.SpotAddress()})
	})

	It("configures spot file par", func() {
		transformer := spot_file_par.EventTransformerInitializer(db)
		name := "spot_file_par"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.SpotAddress()})
	})

	It("configures spot file pip", func() {
		transformer := spot_file_pip.EventTransformerInitializer(db)
		name := "spot_file_pip"
		topic := "0xebecb39d00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.SpotAddress()})
	})

	It("configures spot poke", func() {
		transformer := spot_poke.EventTransformerInitializer(db)
		name := "spot_poke"
		topic := "0xdfd7467e425a8107cfd368d159957692c25085aacbcf5228ce08f10f2146486e"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.SpotAddress()})
	})

	It("configures tend", func() {
		transformer := tend.EventTransformerInitializer(db)
		topic := "0x4b43ed1200000000000000000000000000000000000000000000000000000000"
		addresses := append(flapContracts(), flipContracts()...)
		assertCorrectInit(transformer, "tend", topic, 8928163, addresses)
	})

	It("configures tick", func() {
		transformer := tick.EventTransformerInitializer(db)
		topic := "0xfc7b6aee00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, "tick", topic, 8928163, auctionsContracts())
	})

	It("configures vat deny", func() {
		transformer := vat_deny.EventTransformerInitializer(db)
		name := "vat_deny"
		topic := "0x9c52a7f100000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat file debt ceiling", func() {
		transformer := vat_file_debt_ceiling.EventTransformerInitializer(db)
		name := "vat_file_debt_ceiling"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat file ilk", func() {
		transformer := vat_file_ilk.EventTransformerInitializer(db)
		name := "vat_file_ilk"
		topic := "0x1a0b287e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat flux", func() {
		transformer := vat_flux.EventTransformerInitializer(db)
		name := "vat_flux"
		topic := "0x6111be2e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat fold", func() {
		transformer := vat_fold.EventTransformerInitializer(db)
		name := "vat_fold"
		topic := "0xb65337df00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat fork", func() {
		transformer := vat_fork.EventTransformerInitializer(db)
		name := "vat_fork"
		topic := "0x870c616d00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat frob", func() {
		transformer := vat_frob.EventTransformerInitializer(db)
		name := "vat_frob"
		topic := "0x7608870300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat grab", func() {
		transformer := vat_grab.EventTransformerInitializer(db)
		name := "vat_grab"
		topic := "0x7bab3f4000000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat heal", func() {
		transformer := vat_heal.EventTransformerInitializer(db)
		name := "vat_heal"
		topic := "0xf37ac61c00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat hope", func() {
		transformer := vat_hope.EventTransformerInitializer(db)
		name := "vat_hope"
		topic := "0xa3b22fc400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat init", func() {
		transformer := vat_init.EventTransformerInitializer(db)
		name := "vat_init"
		topic := "0x3b66319500000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat move", func() {
		transformer := vat_move.EventTransformerInitializer(db)
		name := "vat_move"
		topic := "0xbb35783b00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat nope", func() {
		transformer := vat_nope.EventTransformerInitializer(db)
		name := "vat_nope"
		topic := "0xdc4d20fa00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat rely", func() {
		transformer := vat_rely.EventTransformerInitializer(db)
		name := "vat_rely"
		topic := "0x65fae35e00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat slip", func() {
		transformer := vat_slip.EventTransformerInitializer(db)
		name := "vat_slip"
		topic := "0x7cdd3fde00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vat suck", func() {
		transformer := vat_suck.EventTransformerInitializer(db)
		name := "vat_suck"
		topic := "0xf24e23eb00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928152, []string{test_data.VatAddress()})
	})

	It("configures vow fess", func() {
		transformer := vow_fess.EventTransformerInitializer(db)
		name := "vow_fess"
		topic := "0x697efb7800000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928163, []string{test_data.VowAddress()})
	})

	It("configures vow file auction address", func() {
		transformer := vow_file_auction_address.EventTransformerInitializer(db)
		name := "vow_file_auction_address"
		topic := "0xd4e8be8300000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928163, []string{test_data.VowAddress()})
	})

	It("configures vow file auction attributes", func() {
		transformer := vow_file_auction_attributes.EventTransformerInitializer(db)
		name := "vow_file_auction_attributes"
		topic := "0x29ae811400000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928163, []string{test_data.VowAddress()})
	})

	It("configures vow flog", func() {
		transformer := vow_flog.EventTransformerInitializer(db)
		name := "vow_flog"
		topic := "0xd7ee674b00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928163, []string{test_data.VowAddress()})
	})

	It("configures vow heal", func() {
		transformer := vow_heal.EventTransformerInitializer(db)
		name := "vow_heal"
		topic := "0xf37ac61c00000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, name, topic, 8928163, []string{test_data.VowAddress()})
	})

	It("configures yank", func() {
		transformer := yank.EventTransformerInitializer(db)
		topic := "0x26e027f100000000000000000000000000000000000000000000000000000000"
		assertCorrectInit(transformer, "yank", topic, 8928163, auctionsContracts())
	})
})

func assertCorrectInit(transformer event.ITransformer, name, topic string, deployment int64, addresses []string) {
	config := transformer.GetConfig()

	Expect(config.TransformerName).To(Equal(name))
	Expect(config.Topic).To(Equal(topic))
	Expect(config.StartingBlockNumber).To(Equal(deployment))
	Expect(config.ContractAddresses).To(ConsistOf(addresses))
}

func auctionsContracts() []string {
	flapAndFlipContracts := append(flapContracts(), flipContracts()...)
	return append(flapAndFlipContracts, flopContracts()...)
}

func authContracts() []string {
	individualContracts := []string{
		test_data.JugAddress(),
		test_data.PotAddress(),
		test_data.SpotAddress(),
		test_data.VowAddress(),
	}
	withCats := append(individualContracts, catContracts()...)
	withFlaps := append(withCats, flapContracts()...)
	withFlips := append(withFlaps, flipContracts()...)
	withFlops := append(withFlips, flopContracts()...)
	withMedians := append(withFlops, test_data.MedianAddresses()...)
	return append(withMedians, test_data.OsmAddresses()...)
}

func catContracts() []string {
	return []string{test_data.Cat100Address(), test_data.Cat110Address()}
}

func flapContracts() []string {
	return []string{test_data.FlapV100Address(), test_data.FlapV109Address()}
}

func flipContracts() []string {
	return append(test_data.FlipV100Addresses(), test_data.FlipV110Addresses()...)
}

func flopContracts() []string {
	return []string{test_data.FlopV101Address(), test_data.FlopV109Address()}
}
