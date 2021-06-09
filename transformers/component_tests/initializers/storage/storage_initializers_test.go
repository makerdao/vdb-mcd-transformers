package storage

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	cat_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_0_0/initializer"
	cat_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_1_0/initializer"
	cdp_manager "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager/initializer"
	clip_link_a_v1_0_3 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip/initializers/link_a/v1_3_0"
	dog_v1_3_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog/initializers/v1_3_0"
	flap_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializers/v1_0_0"
	flap_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializers/v1_0_9"
	flip_aave_a_v1_2_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/aave_a/v1_2_2"
	flip_bal_a_v1_1_14 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bal_a/v1_1_14"
	flip_bat_a_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_0_0"
	flip_bat_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_0_9"
	flip_bat_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_1_0"
	flip_comp_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/comp_a/v1_1_2"
	flip_eth_a_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_a/v1_0_0"
	flip_eth_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_a/v1_0_9"
	flip_eth_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_a/v1_1_0"
	flip_eth_b_v1_1_3 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_b/v1_1_3"
	flip_gusd_a_v1_1_5 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/gusd_a/v1_1_5"
	flip_knc_a_v1_0_8 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_a/v1_0_8"
	flip_knc_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_a/v1_0_9"
	flip_knc_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/knc_a/v1_1_0"
	flip_link_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/link_a/v1_1_2"
	flip_lrc_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/lrc_a/v1_1_2"
	flip_mana_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/mana_a/v1_0_9"
	flip_mana_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/mana_a/v1_1_0"
	flip_paxusd_a_v1_1_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/paxusd_a/v1_1_1"
	flip_rentbtc_a_v1_2_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/renbtc_a/v1_2_1"
	flip_sai_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/sai/v1_0_0"
	flip_tusd_a_v1_0_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/tusd_a/v1_0_7"
	flip_tusd_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/tusd_a/v1_0_9"
	flip_tusd_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/tusd_a/v1_1_0"
	flip_uni_a_v1_2_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/uni_a/v1_2_1"
	flip_univ2aaveeth_a_v1_2_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2aaveeth_a/v1_2_7"
	flip_univ2daieth_a_v1_2_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2daieth_a/v1_2_2"
	flip_univ2daiusdc_a_v1_2_5 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2daiusdc_a/v1_2_5"
	flip_univ2daiusdt_a_v1_2_8 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2daiusdt_a/v1_2_8"
	flip_univ2ethusdt_a_v1_2_5 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2ethusdt_a/v1_2_5"
	flip_univ2linketh_a_v1_2_6 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2linketh_a/v1_2_6"
	flip_univ2unieth_a_v1_2_6 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2unieth_a/v1_2_6"
	flip_univ2usdceth_a_v1_2_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2usdceth_a/v1_2_4"
	flip_univ2wbtcdai_a_v1_2_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2wbtcdai_a/v1_2_7"
	flip_univ2wbtceth_a_v1_2_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/univ2wbtceth_a/v1_2_4"
	flip_usdc_a_v1_0_4 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a/v1_0_4"
	flip_usdc_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a/v1_0_9"
	flip_usdc_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_a/v1_1_0"
	flip_usdc_b_v1_0_7 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b/v1_0_7"
	flip_usdc_b_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b/v1_0_9"
	flip_usdc_b_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdc_b/v1_1_0"
	flip_usdt_a_v1_1_1 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/usdt_a/v1_1_1"
	flip_wbtc_a_v1_0_6 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_a/v1_0_6"
	flip_wbtc_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_a/v1_0_9"
	flip_wbtc_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/wbtc_a/v1_1_0"
	flip_yfi_a_v1_1_14 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/yfi_a/v1_1_14"
	flip_zrx_a_v1_0_8 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_a/v1_0_8"
	flip_zrx_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_a/v1_0_9"
	flip_zrx_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/zrx_a/v1_1_0"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage transformer initializers", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())

	It("configures the v1_0_0 cat", func() {
		address := "0x78F2c2AF65126834c51822F56Be0d7469D7A523E"
		transformer := cat_v1_0_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 cat", func() {
		address := "0xa5679C04fc3d9d8b0AaB1F0ab83555b301cA70Ea"
		transformer := cat_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the cdp_manager", func() {
		address := "0x5ef30b9986345249bc32d8928B7ee64DE9435E39"
		transformer := cdp_manager.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_3_0 clip link_a", func() {
		address := test_data.ClipLinkAV130Address()
		transformer := clip_link_a_v1_0_3.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the dog v1_3_0", func() {
		address := test_data.Dog130Address()
		transformer := dog_v1_3_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_0 flap", func() {
		address := "0xdfE0fb1bE2a52CDBf8FB962D5701d7fd0902db9f"
		transformer := flap_v1_0_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flap", func() {
		address := "0xC4269cC7acDEdC3794b221aA4D9205F564e27f0d"
		transformer := flap_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_2 flip aave", func() {
		address := test_data.FlipAaveAV122Address()
		transformer := flip_aave_a_v1_2_2.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_14 flip bal", func() {
		address := test_data.FlipBalAV1114Address()
		transformer := flip_bal_a_v1_1_14.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_0 flip bat_a", func() {
		address := "0xaA745404d55f88C108A28c86abE7b5A1E7817c07"
		transformer := flip_bat_a_v1_0_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip bat_a", func() {
		address := "0x5EdF770FC81E7b8C2c89f71F30f211226a4d7495"
		transformer := flip_bat_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip bat_a", func() {
		address := "0xF7C569B2B271354179AaCC9fF1e42390983110BA"
		transformer := flip_bat_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_2 flip comp_a", func() {
		address := "0x524826F84cB3A19B6593370a5889A58c00554739"
		transformer := flip_comp_a_v1_1_2.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_0 flip eth_a", func() {
		address := "0xd8a04F5412223F513DC55F839574430f5EC15531"
		transformer := flip_eth_a_v1_0_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip eth_a", func() {
		address := "0x0F398a2DaAa134621e4b687FCcfeE4CE47599Cc1"
		transformer := flip_eth_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip eth_a", func() {
		address := "0xF32836B9E1f47a0515c6Ec431592D5EbC276407f"
		transformer := flip_eth_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_3 flip eth_b", func() {
		address := "0xD499d71bE9e9E5D236A07ac562F7B6CeacCa624c"
		transformer := flip_eth_b_v1_1_3.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_15 flip gusd_a", func() {
		address := test_data.FlipGusdAV115Address()
		transformer := flip_gusd_a_v1_1_5.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_8 flip knc_a", func() {
		address := "0xAbBCB9Ae89cDD3C27E02D279480C7fF33083249b"
		transformer := flip_knc_a_v1_0_8.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip knc_a", func() {
		address := "0xAD4a0B5F3c6Deb13ADE106Ba6E80Ca6566538eE6"
		transformer := flip_knc_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip knc_a", func() {
		address := "0x57B01F1B3C59e2C0bdfF3EC9563B71EEc99a3f2f"
		transformer := flip_knc_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_2 flip link_a", func() {
		address := "0xB907EEdD63a30A3381E6D898e5815Ee8c9fd2c85"
		transformer := flip_link_a_v1_1_2.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_2 flip lrc_a", func() {
		address := "0x7FdDc36dcdC435D8F54FDCB3748adcbBF70f3dAC"
		transformer := flip_lrc_a_v1_1_2.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip mana_a", func() {
		address := "0x4bf9D2EBC4c57B9B783C12D30076507660B58b3a"
		transformer := flip_mana_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip mana_a", func() {
		address := "0x0a1D75B4f49BA80724a214599574080CD6B68357"
		transformer := flip_mana_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_1 flip paxusd_a", func() {
		address := "0x52D5D1C05CC79Fc24A629Cb24cB06C5BE5d766E7"
		transformer := flip_paxusd_a_v1_1_1.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_1 flip renbtc_a", func() {
		address := test_data.FlipRenbtcA121Address()
		transformer := flip_rentbtc_a_v1_2_1.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_0 flip sai", func() {
		address := "0x5432b2f3c0DFf95AA191C45E5cbd539E2820aE72"
		transformer := flip_sai_v1_0_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_7 flip tusd_a", func() {
		address := "0xba3f6a74BD12Cf1e48d4416c7b50963cA98AfD61"
		transformer := flip_tusd_a_v1_0_7.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip tusd_a", func() {
		address := "0x04C42fAC3e29Fd27118609a5c36fD0b3Cb8090b3"
		transformer := flip_tusd_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip tusd_a", func() {
		address := "0x9E4b213C4defbce7564F2Ac20B6E3bF40954C440"
		transformer := flip_tusd_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_1 flip uni_a", func() {
		address := test_data.FlipUniAV121Address()
		transformer := flip_uni_a_v1_2_1.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_7 flip univ2aaveeth_a", func() {
		address := test_data.FlipUniV2AaveEthAddress()
		transformer := flip_univ2aaveeth_a_v1_2_7.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_2 flip univ2daieth_a", func() {
		address := test_data.FlipUniV2DaiEthAddress()
		transformer := flip_univ2daieth_a_v1_2_2.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_5 flip univ2daiusdc_a", func() {
		address := test_data.FlipUniV2DaiUsdcAddress()
		transformer := flip_univ2daiusdc_a_v1_2_5.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_5 flip univ2ethusdt_a", func() {
		address := test_data.FlipUniV2EthUsdtAddress()
		transformer := flip_univ2ethusdt_a_v1_2_5.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_6 flip univ2linketh_a", func() {
		address := test_data.FlipUniV2LinkEthAddress()
		transformer := flip_univ2linketh_a_v1_2_6.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_6 flip univ2unieth_a", func() {
		address := test_data.FlipUniV2UniEthAddress()
		transformer := flip_univ2unieth_a_v1_2_6.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_4 flip univ2usdceth_a", func() {
		address := test_data.FlipUniV2UsdcEthAddress()
		transformer := flip_univ2usdceth_a_v1_2_4.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_7 flip univ2wbtcdai_a", func() {
		address := test_data.FlipUniV2WbtcDaiAddress()
		transformer := flip_univ2wbtcdai_a_v1_2_7.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_4 flip univ2wbtceth_a", func() {
		address := test_data.FlipUniV2WbtcEthAddress()
		transformer := flip_univ2wbtceth_a_v1_2_4.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_2_8 flip univ2daiusdt_a", func() {
		address := test_data.FlipUniV2DaiUsdtAddress()
		transformer := flip_univ2daiusdt_a_v1_2_8.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_4 flip usdc_a", func() {
		address := "0xE6ed1d09a19Bd335f051d78D5d22dF3bfF2c28B1"
		transformer := flip_usdc_a_v1_0_4.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip usdc_a", func() {
		address := "0x545521e0105C5698f75D6b3C3050CfCC62FB0C12"
		transformer := flip_usdc_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip usdc_a", func() {
		address := "0xbe359e53038E41a1ffA47DAE39645756C80e557a"
		transformer := flip_usdc_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_7 flip usdc_b", func() {
		address := "0xec25Ca3fFa512afbb1784E17f1D414E16D01794F"
		transformer := flip_usdc_b_v1_0_7.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip usdc_b", func() {
		address := "0x6002d3B769D64A9909b0B26fC00361091786fe48"
		transformer := flip_usdc_b_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip usdc_b", func() {
		address := "0x77282aD36aADAfC16bCA42c865c674F108c4a616"
		transformer := flip_usdc_b_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_1 flip usdt_a", func() {
		address := "0x667F41d0fDcE1945eE0f56A79dd6c142E37fCC26"
		transformer := flip_usdt_a_v1_1_1.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_6 flip wbtc_a", func() {
		address := "0x3E115d85D4d7253b05fEc9C0bB5b08383C2b0603"
		transformer := flip_wbtc_a_v1_0_6.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip wbtc_a", func() {
		address := "0xF70590Fa4AaBe12d3613f5069D02B8702e058569"
		transformer := flip_wbtc_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip wbtc_a", func() {
		address := "0x58CD24ac7322890382eE45A3E4F903a5B22Ee930"
		transformer := flip_wbtc_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_14 flip yfi_a", func() {
		address := test_data.FlipYfiAV1114Address()
		transformer := flip_yfi_a_v1_1_14.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_8 flip zrx_a", func() {
		address := "0x08c89251FC058cC97d5bA5F06F95026C0A5CF9B0"
		transformer := flip_zrx_a_v1_0_8.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_0_9 flip zrx_a", func() {
		address := "0x92645a34d07696395b6e5b8330b000D0436A9aAD"
		transformer := flip_zrx_a_v1_0_9.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})

	It("configures the v1_1_0 flip zrx_a", func() {
		address := "0xa4341cAf9F9F098ecb20fb2CeE2a0b8C78A18118"
		transformer := flip_zrx_a_v1_1_0.StorageTransformerInitializer(db)

		Expect(transformer.GetContractAddress().String()).To(Equal(address))
	})
})
