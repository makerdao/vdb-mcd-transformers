package config_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config/test_data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config Parser", func() {
	var (
		testConfigFilePath = "/Users/elizabethengelman/go/src/github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config/test_data/"
		testConfigFileName = "testConfig"
		configParser       = config.Parser{}
	)

	Context("ParseCurrentConfigFile", func() {
		It("reads in the exporter metadata", func() {
			expectedExporterMetadata := config.ExporterMetaData{
				Home:             "github.com/makerdao/vulcanizedb",
				Name:             "transformerExporter",
				Save:             false,
				Schema:           "maker",
				TransformerNames: []string{"cat_v1_1_0", "cat_file_vow"},
			}

			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, testConfigFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.ExporterMetadata).To(Equal(expectedExporterMetadata))
		})

		It("reads in the exporterTransformers", func() {
			expectedTransformerExporters := config.TransformerExporters{
				"cat_v1_1_0":   test_data.Cat110Exporter,
				"cat_file_vow": test_data.CatFileVowExporter,
			}
			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, testConfigFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.TransformerExporters).To(Equal(expectedTransformerExporters))
		})

		It("reads in the contracts", func() {
			expectedContracts := config.Contracts{
				"MCD_CAT_1_0_0": test_data.Cat100Contract,
				"MCD_CAT_1_1_0": test_data.Cat110Contract,
			}
			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, testConfigFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.Contracts).To(Equal(expectedContracts))
		})
	})
})
