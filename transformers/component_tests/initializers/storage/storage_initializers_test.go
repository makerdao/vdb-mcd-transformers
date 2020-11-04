package storage

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	cat_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_0_0/initializer"
	cat_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_1_0/initializer"
	cdp_manager "github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager/initializer"
	flap_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializers/v1_0_0"
	flap_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap/initializers/v1_0_9"
	flip_bat_a_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_0_0"
	flip_bat_a_v1_0_9 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_0_9"
	flip_bat_a_v1_1_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/bat_a/v1_1_0"
	flip_comp_a_v1_1_2 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/comp_a/v1_1_2"
	flip_eth_a_v1_0_0 "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers/eth_a/v1_0_0"
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
})