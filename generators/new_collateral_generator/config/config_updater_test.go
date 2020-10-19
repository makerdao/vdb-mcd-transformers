package config_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config/test_data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewConfigUpdater", func() {
	Context("AddNewCollateralToConfig", func() {
		var configUpdater = config.NewConfigUpdater(test_data.EthBCollateral, test_data.EthBContracts)
		configUpdater.SetInitialConfig(test_data.InitialConfig)

		It("adds new transformer names to the exporter metadata for the new collateral", func() {
			expectedExporterMetadata := config.ExporterMetaData{
				Home:   "github.com/makerdao/vulcanizedb",
				Name:   "transformerExporter",
				Save:   false,
				Schema: "maker",
				TransformerNames: []string{
					"cat_v1_1_0",
					"cat_file_vow",
					"flip_eth_b_v1_1_3", // new storage flip transformer
					"median_eth_b",      // new median eth transformer
				},
			}

			addErr := configUpdater.AddNewCollateralToConfig()
			Expect(addErr).NotTo(HaveOccurred())
			Expect(configUpdater.UpdatedConfig.ExporterMetadata).To(Equal(expectedExporterMetadata))
		})

		It("adds new storage exporterTransformers for new collateral", func() {
			addEthBErr := configUpdater.AddNewCollateralToConfig()
			Expect(addEthBErr).NotTo(HaveOccurred())
			Expect(configUpdater.UpdatedConfig.TransformerExporters).To(
				HaveKeyWithValue("exporter.flip_eth_b_v1_1_3", test_data.FlipEthBStorageExporter))
			Expect(configUpdater.UpdatedConfig.TransformerExporters).To(
				HaveKeyWithValue("exporter.median_eth_b", test_data.MedianEthBStorageExporter))
		})

		It("adds the new collateral flip contract to event exporters", func() {
			addEthBErr := configUpdater.AddNewCollateralToConfig()
			Expect(addEthBErr).NotTo(HaveOccurred())
			denyExporter := configUpdater.UpdatedConfig.TransformerExporters["exporter.deny"]
			Expect(denyExporter.Contracts).To(ContainElement(test_data.FlipEthBContractName))
		})

		It("adds median contract to event exporters", func() {
			addEthBErr := configUpdater.AddNewCollateralToConfig()
			Expect(addEthBErr).NotTo(HaveOccurred())
			denyExporter := configUpdater.UpdatedConfig.TransformerExporters["exporter.deny"]
			logMedianExporter := configUpdater.UpdatedConfig.TransformerExporters["exporter.log_median_price"]
			Expect(denyExporter.Contracts).To(ContainElement(test_data.MedianEthBContractName))
			Expect(logMedianExporter.Contracts).To(ContainElement(test_data.MedianEthBContractName))
		})

		It("adds osm contract to event exporters", func() {
			addEthBErr := configUpdater.AddNewCollateralToConfig()
			Expect(addEthBErr).NotTo(HaveOccurred())
			denyExporter := configUpdater.UpdatedConfig.TransformerExporters["exporter.deny"]
			logValueExporter := configUpdater.UpdatedConfig.TransformerExporters["exporter.log_value"]
			Expect(denyExporter.Contracts).To(ContainElement(test_data.OsmEthBContractName))
			Expect(logValueExporter.Contracts).To(ContainElement(test_data.OsmEthBContractName))
		})

		It("does not add flip, median or osm contracts to event exporters that don't currently have those contract types", func() {
			addEthBErr := configUpdater.AddNewCollateralToConfig()
			Expect(addEthBErr).NotTo(HaveOccurred())
			catVowExporter := configUpdater.UpdatedConfig.TransformerExporters["exporter.cat_file_vow"]
			Expect(catVowExporter.Contracts).NotTo(ContainElement(test_data.FlipEthBContractName))
			Expect(catVowExporter.Contracts).NotTo(ContainElement(test_data.MedianEthBContractName))
			Expect(catVowExporter.Contracts).NotTo(ContainElement(test_data.OsmEthBContractName))
		})

		It("adds new flip, median and osm contracts for new collateral", func() {
			expectedContracts := config.Contracts{
				"MCD_CAT_1_0_0":        test_data.Cat100Contract,
				"MCD_CAT_1_1_0":        test_data.Cat110Contract,
				"MCD_FLIP_ETH_B_1_1_3": test_data.FlipEthBContract,
				"MEDIAN_ETH_B":         test_data.MedianEthBContract,
				"OSM_ETH_B":            test_data.OsmEthBContract,
			}
			addEthBErr := configUpdater.AddNewCollateralToConfig()
			Expect(addEthBErr).NotTo(HaveOccurred())
			Expect(configUpdater.UpdatedConfig.Contracts).To(Equal(expectedContracts))
		})

		It("doesn't update the initialConfig", func() {
			testCollateral := config.Collateral{Name: "TEST", Version: "1.0.0"}
			testContracts := config.Contracts{"flip": test_data.FlipEthBContract}
			configUpdater := config.NewConfigUpdater(testCollateral, testContracts)
			configUpdater.SetInitialConfig(test_data.InitialConfig)
			addErr := configUpdater.AddNewCollateralToConfig()
			Expect(addErr).NotTo(HaveOccurred())
			Expect(configUpdater.UpdatedConfig.ExporterMetadata).NotTo(Equal(configUpdater.InitialConfig.ExporterMetadata))
		})
	})
})
