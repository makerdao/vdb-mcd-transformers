package generator_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/generator"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/generator/test_data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config Parser", func() {
	var configParser generator.ConfigParser

	Context("ParseCurrentConfigFile", func() {
		It("reads in the exporter metadata", func() {
			expectedExporterMetadata := generator.ExporterMetaData{
				Home:             "github.com/makerdao/vulcanizedb",
				Name:             "transformerExporter",
				Save:             false,
				Schema:           "maker",
				TransformerNames: []string{"cat_v1_1_0", "cat_file_vow"},
			}

			config, parseErr := configParser.ParseCurrentConfig(test_data.TestConfigString)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.ExporterMetadata).To(Equal(expectedExporterMetadata))
		})

		It("reads in the exporterTransformers", func() {
			expectedTransformerExporters := generator.TransformerExporters{
				"exporter.cat_v1_1_0": test_data.Cat110Exporter,
				"exporter.cat_file_vow": test_data.CatFileVowExporter,
			}
			config, parseErr := configParser.ParseCurrentConfig(test_data.TestConfigString)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.TransformerExporters).To(Equal(expectedTransformerExporters))
		})

		It("reads in the contracts", func() {
			expectedContracts := generator.Contracts{
				"MCD_CAT_1_0_0": test_data.Cat100Contract,
				"MCD_CAT_1_1_0": test_data.Cat110Contract,
			}
			config, parseErr := configParser.ParseCurrentConfig(test_data.TestConfigString)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.Contracts).To(Equal(expectedContracts))
		})
	})
})

