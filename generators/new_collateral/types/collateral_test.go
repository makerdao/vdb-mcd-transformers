package types_test

import (
	"path/filepath"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collateral", func() {
	var (
		collateral1 = types.NewCollateral("ETH-B", "1.2.3")
		collateral2 = types.NewCollateral("PAXG", "1_2_3")
		collateral3 = types.NewCollateral("Eth_B", "1_2_3")
	)

	It("formats the version", func() {
		Expect(collateral1.FormattedVersion()).To(Equal("v1_2_3"))
		Expect(collateral2.FormattedVersion()).To(Equal("v1_2_3"))
		Expect(collateral3.FormattedVersion()).To(Equal("v1_2_3"))
	})

	Context("get storage transformer name", func() {
		It("formats the collateral as a flip transformer name", func() {
			Expect(collateral1.GetFlipTransformerName()).To(Equal("flip_eth_b_v1_2_3"))
			Expect(collateral2.GetFlipTransformerName()).To(Equal("flip_paxg_v1_2_3"))
			Expect(collateral3.GetFlipTransformerName()).To(Equal("flip_eth_b_v1_2_3"))
		})

		It("formats the collateral as a median transformer name", func() {
			Expect(collateral1.GetMedianTransformerName()).To(Equal("median_eth_b"))
			Expect(collateral2.GetMedianTransformerName()).To(Equal("median_paxg"))
			Expect(collateral3.GetMedianTransformerName()).To(Equal("median_eth_b"))
		})
	})

	Context("get contract names", func() {
		It("formats the collateral as a flip contract name", func() {
			Expect(collateral1.GetFlipContractName()).To(Equal("MCD_FLIP_ETH_B_1_2_3"))
			Expect(collateral2.GetFlipContractName()).To(Equal("MCD_FLIP_PAXG_1_2_3"))
			Expect(collateral3.GetFlipContractName()).To(Equal("MCD_FLIP_ETH_B_1_2_3"))
		})

		It("formats the collateral as a median contract name", func() {
			Expect(collateral1.GetMedianContractName()).To(Equal("MEDIAN_ETH_B"))
			Expect(collateral2.GetMedianContractName()).To(Equal("MEDIAN_PAXG"))
			Expect(collateral3.GetMedianContractName()).To(Equal("MEDIAN_ETH_B"))
		})

		It("formats the collateral as a osm contract name", func() {
			Expect(collateral1.GetOsmContractName()).To(Equal("OSM_ETH_B"))
			Expect(collateral2.GetOsmContractName()).To(Equal("OSM_PAXG"))
			Expect(collateral3.GetOsmContractName()).To(Equal("OSM_ETH_B"))
		})
	})

	Context("get storage initializer directory", func() {
		It("formats the collateral as a flip initializer directory", func() {
			Expect(collateral1.GetFlipInitializerDirectory()).To(Equal("eth_b/v1_2_3"))
			Expect(collateral2.GetFlipInitializerDirectory()).To(Equal("paxg/v1_2_3"))
			Expect(collateral3.GetFlipInitializerDirectory()).To(Equal("eth_b/v1_2_3"))
		})

		It("formats the collateral as a median initializer directory", func() {
			Expect(collateral1.GetMedianInitializerDirectory()).To(Equal("median_eth_b"))
			Expect(collateral2.GetMedianInitializerDirectory()).To(Equal("median_paxg"))
			Expect(collateral3.GetMedianInitializerDirectory()).To(Equal("median_eth_b"))
		})
	})

	Context("get absolute initializer directory path", func() {
		It("formats the collateral as a flip initializer path", func() {
			flipInitializersPath := filepath.Join(helpers.GetProjectPath(), "transformers", "storage", "flip", "initializers")
			Expect(collateral1.GetAbsoluteFlipStorageInitializersDirectoryPath()).To(Equal(filepath.Join(flipInitializersPath, "eth_b/v1_2_3")))
			Expect(collateral2.GetAbsoluteFlipStorageInitializersDirectoryPath()).To(Equal(filepath.Join(flipInitializersPath, "paxg/v1_2_3")))
			Expect(collateral3.GetAbsoluteFlipStorageInitializersDirectoryPath()).To(Equal(filepath.Join(flipInitializersPath, "eth_b/v1_2_3")))
		})

		It("formats the collateral as a median initializer path", func() {
			medianInitializersPath := filepath.Join(helpers.GetProjectPath(), "transformers", "storage", "median", "initializers")
			Expect(collateral1.GetAbsoluteMedianStorageInitializersDirectoryPath()).To(Equal(filepath.Join(medianInitializersPath, "median_eth_b")))
			Expect(collateral2.GetAbsoluteMedianStorageInitializersDirectoryPath()).To(Equal(filepath.Join(medianInitializersPath, "median_paxg")))
			Expect(collateral3.GetAbsoluteMedianStorageInitializersDirectoryPath()).To(Equal(filepath.Join(medianInitializersPath, "median_eth_b")))
		})
	})

	Context("get absolute initializer file path", func() {
		It("formats the collateral as a flip initializer file path", func() {
			flipInitializersPath := filepath.Join(helpers.GetProjectPath(), "transformers", "storage", "flip", "initializers")
			Expect(collateral1.GetAbsoluteFlipStorageInitializerFilePath()).To(Equal(filepath.Join(flipInitializersPath, "eth_b/v1_2_3", "initializer.go")))
			Expect(collateral2.GetAbsoluteFlipStorageInitializerFilePath()).To(Equal(filepath.Join(flipInitializersPath, "paxg/v1_2_3", "initializer.go")))
			Expect(collateral3.GetAbsoluteFlipStorageInitializerFilePath()).To(Equal(filepath.Join(flipInitializersPath, "eth_b/v1_2_3", "initializer.go")))
		})

		It("formats the collateral as a median initializer file path", func() {
			medianInitializersPath := filepath.Join(helpers.GetProjectPath(), "transformers", "storage", "median", "initializers")
			Expect(collateral1.GetAbsoluteMedianStorageInitializerFilePath()).To(Equal(filepath.Join(medianInitializersPath, "median_eth_b", "initializer.go")))
			Expect(collateral2.GetAbsoluteMedianStorageInitializerFilePath()).To(Equal(filepath.Join(medianInitializersPath, "median_paxg", "initializer.go")))
			Expect(collateral3.GetAbsoluteMedianStorageInitializerFilePath()).To(Equal(filepath.Join(medianInitializersPath, "median_eth_b", "initializer.go")))
		})
	})
})
