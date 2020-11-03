package config_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/test_data"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config Parser", func() {
	var (
		testConfigFilePath   = "../test_data/"
		testConfigFileName   = "testConfig"
		configParser         = config.NewParser()
		expectedParsedConfig = types.TransformersConfig{
			ExporterMetadata: types.ExporterMetaData{
				Home:             "github.com/makerdao/vulcanizedb",
				Name:             "transformerExporter",
				Save:             false,
				Schema:           "maker",
				TransformerNames: []string{"cat_v1_1_0", "cat_file_vow"},
			},
			TransformerExporters: types.TransformerExporters{
				"cat_v1_1_0":   test_data.Cat110Exporter,
				"cat_file_vow": test_data.CatFileVowExporter,
			},
			Contracts: types.Contracts{
				"MCD_CAT_1_0_0": test_data.Cat100Contract,
				"MCD_CAT_1_1_0": test_data.Cat110Contract,
			},
		}
	)

	Context("ParseCurrentConfigFile", func() {
		It("returns an error if it fails to decode the file", func() {
			configFile    := "non-existent-file"
			_, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, configFile)
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr).To(MatchError("error decoding config file: open ../test_data/non-existent-file.toml: no such file or directory"))
		})

		It("parses metadata", func() {
			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, testConfigFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.ExporterMetadata).To(Equal(expectedParsedConfig.ExporterMetadata))
		})

		It("returns an error if it fails to parse the exporter metadata", func() {
			configFileName := "testConfigWithBadMetadata"
			_, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, configFileName)
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr).To(MatchError(
				"error parsing exporter metadata from config file: error asserting exporterMetadata types - homeOk: false, nameOk: false, saveOk: false, schemaOk: false",
			))
		})

		It("can handle an empty transformerNames slice", func() {
			configFileName := "testConfigWithNoTransformerNames"
			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, configFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.ExporterMetadata).To(Equal(types.ExporterMetaData{
				Home:             "github.com/makerdao/vulcanizedb",
				Name:             "transformerExporter",
				Save:             true,
				Schema:           "maker",
				TransformerNames: nil,
			}))
		})

		It("parses transformerExporters", func() {
			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, testConfigFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config.TransformerExporters).To(Equal(expectedParsedConfig.TransformerExporters))
		})

		It("returns an error if it fails to decode an exporterValue", func() {
			// the exporter.cat_v1_1_0 will not properly decode into a types.TransformerExporter because the Path field doesn't match
			configFile := "testConfigWithBadTransformerExporter"
			_, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, configFile)
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr).To(MatchError(
				"error parsing transformer exporters from config file: error decoding transformerExporters:" +
					" 1 error(s) decoding:\n\n* 'Path' expected type 'string', got unconvertible type 'int64'",))
		})
	})
})
