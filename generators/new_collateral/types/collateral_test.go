package types_test

import (
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
		Expect(collateral1.FormattedVersionWithPrependedV()).To(Equal("v1_2_3"))
		Expect(collateral2.FormattedVersionWithPrependedV()).To(Equal("v1_2_3"))
		Expect(collateral3.FormattedVersionWithPrependedV()).To(Equal("v1_2_3"))
	})

	Context("get storage transformer name", func() {
		It("formats the collateral as a flip transformer name", func() {
			Expect(collateral1.GetFlipTransformerName()).To(Equal("flip_eth_b_v1_2_3"))
			Expect(collateral2.GetFlipTransformerName()).To(Equal("flip_paxg_v1_2_3"))
			Expect(collateral3.GetFlipTransformerName()).To(Equal("flip_eth_b_v1_2_3"))
		})

		It("formats the collateral as a median transformer name", func() {
			Expect(collateral1.GetMedianTransformerName()).To(Equal("median_eth_b_v1_2_3"))
			Expect(collateral2.GetMedianTransformerName()).To(Equal("median_paxg_v1_2_3"))
			Expect(collateral3.GetMedianTransformerName()).To(Equal("median_eth_b_v1_2_3"))
		})
	})

	Context("get contract names", func() {
		It("formats the collateral as a flip contract name", func() {
			Expect(collateral1.GetFlipContractName()).To(Equal("MCD_FLIP_ETH_B_1_2_3"))
			Expect(collateral2.GetFlipContractName()).To(Equal("MCD_FLIP_PAXG_1_2_3"))
			Expect(collateral3.GetFlipContractName()).To(Equal("MCD_FLIP_ETH_B_1_2_3"))
		})

		It("formats the collateral as a median contract name", func() {
			Expect(collateral1.GetMedianContractName()).To(Equal("MEDIAN_ETH_B_1_2_3"))
			Expect(collateral2.GetMedianContractName()).To(Equal("MEDIAN_PAXG_1_2_3"))
			Expect(collateral3.GetMedianContractName()).To(Equal("MEDIAN_ETH_B_1_2_3"))
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
			Expect(collateral1.GetMedianInitializerDirectory()).To(Equal("median_eth_b/v1_2_3"))
			Expect(collateral2.GetMedianInitializerDirectory()).To(Equal("median_paxg/v1_2_3"))
			Expect(collateral3.GetMedianInitializerDirectory()).To(Equal("median_eth_b/v1_2_3"))
		})
	})
})
