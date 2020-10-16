package config_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Describe("Collateral", func() {
		var (
			collateral1 = config.NewCollateral("ETH-B", "1.2.3")
			collateral2 = config.NewCollateral("PAXG", "1_2_3")
			collateral3 = config.NewCollateral("Eth_B", "1_2_3")
		)

		Context("flip contracts", func() {
			It("formats the collateral as a flip transformer name", func() {
				Expect(collateral1.FormattedForFlipTransformerName()).To(Equal("eth_b_v1_2_3"))
				Expect(collateral2.FormattedForFlipTransformerName()).To(Equal("paxg_v1_2_3"))
				Expect(collateral3.FormattedForFlipTransformerName()).To(Equal("eth_b_v1_2_3"))
			})

			It("formats the collateral as a flip initializer file name", func() {
				Expect(collateral1.FormattedForFlipInitializerFileName()).To(Equal("eth_b/v1_2_3"))
				Expect(collateral2.FormattedForFlipInitializerFileName()).To(Equal("paxg/v1_2_3"))
				Expect(collateral3.FormattedForFlipInitializerFileName()).To(Equal("eth_b/v1_2_3"))
			})

			It("formats the collateral as a flip contract name", func() {
				Expect(collateral1.FormattedForFlipContractName()).To(Equal("MCD_FLIP_ETH_B_1_2_3"))
				Expect(collateral2.FormattedForFlipContractName()).To(Equal("MCD_FLIP_PAXG_1_2_3"))
				Expect(collateral3.FormattedForFlipContractName()).To(Equal("MCD_FLIP_ETH_B_1_2_3"))
			})
		})

		Context("median contracts", func() {
			It("formats the collateral as a median transformer name", func() {
				Expect(collateral1.FormattedForMedianTransformerName()).To(Equal("eth_b"))
				Expect(collateral2.FormattedForMedianTransformerName()).To(Equal("paxg"))
				Expect(collateral3.FormattedForMedianTransformerName()).To(Equal("eth_b"))
			})

			It("formats the collateral as a median contract name", func() {
				Expect(collateral1.FormattedForMedianContractName()).To(Equal("MEDIAN_ETH_B"))
				Expect(collateral2.FormattedForMedianContractName()).To(Equal("MEDIAN_PAXG"))
				Expect(collateral3.FormattedForMedianContractName()).To(Equal("MEDIAN_ETH_B"))
			})
		})

		Context("osm contracts", func() {
			It("formats the collateral as a osm contract name", func() {
				Expect(collateral1.FormattedForOsmContractName()).To(Equal("OSM_ETH_B"))
				Expect(collateral2.FormattedForOsmContractName()).To(Equal("OSM_PAXG"))
				Expect(collateral3.FormattedForOsmContractName()).To(Equal("OSM_ETH_B"))
			})
		})
	})
})
