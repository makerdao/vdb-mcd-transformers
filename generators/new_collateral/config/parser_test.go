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
		badTestConfigFile = "non-existent-file"
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
		It("reads in the existing config file", func() {
			config, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, testConfigFileName)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config).To(Equal(expectedParsedConfig))
		})

		It("returns an error if it fails to decode the file", func() {
			_, parseErr := configParser.ParseCurrentConfig(testConfigFilePath, badTestConfigFile)
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr).To(MatchError("open ../test_data/non-existent-file.toml: no such file or directory"))
		})
	})

	Context("ParseExporterMetadata", func() {
		It("parses metadata", func() {
			tomlConfig := types.TransformersConfigForToml{
				Exporter:  map[string]interface{}{
					"home":"home", "name":"name", "save":true, "schema":"schema",
					"transformerNames": []string{"test-transformer"},
				},
			}
			config, parseErr := config.ParseExporterMetaData(tomlConfig)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(config).To(Equal(types.ExporterMetaData{
				Home:             "home",
				Name:             "name",
				Save:             true,
				Schema:           "schema",
				TransformerNames: []string{"test-transformer"},
			}))
		})

		It("returns an error if it fails to parse the exporter metadata", func() {
			tomlConfig := types.TransformersConfigForToml{Exporter:  nil}
			_, parseErr := config.ParseExporterMetaData(tomlConfig)
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr).To(MatchError(
				"error asserting exporterMetadata types - homeOk: false, nameOk: false, saveOk: false, schemaOk: false",
			))
		})

		It("can handle an empty transformerNames slice", func() {
			tomlConfig := types.TransformersConfigForToml{
				Exporter:  map[string]interface{}{
					"home":"home", "name":"name", "save":true, "schema":"schema",
				},
			}
			exporterMetadata, parseErr := config.ParseExporterMetaData(tomlConfig)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(exporterMetadata).To(Equal(types.ExporterMetaData{
				Home:             "home",
				Name:             "name",
				Save:             true,
				Schema:           "schema",
				TransformerNames: nil,
			}))
		})
	})

	Context("ParseTransformerExporters", func() {
		It("parses transformerExporters", func() {
			transformerExporterMap := map[string]interface{}{
				"path":"path",
				"type":"type",
				"repository":"repository",
				"migrations":"migrations",
				"contracts":[]string{"testContract"},
				"rank":"0",
			}
			tomlConfig := types.TransformersConfigForToml{
				Exporter:  map[string]interface{}{
					"test-exporter": transformerExporterMap,
				},
			}
			transformerExporters, parseErr := config.ParseTransformerExporters(tomlConfig)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(transformerExporters).To(Equal(types.TransformerExporters{
				"test-exporter": {
					Path:       "path",
					Type:       "type",
					Repository: "repository",
					Migrations: "migrations",
					Contracts:  []string{"testContract"},
					Rank:       "0",
				},
			}))
		})

		It("returns an error if it fails to decode an exporterValue", func() {
			// this map will not properly decode into a types.TransformerExporter because the Path field doesn't match
			notATransformerExporterMap := map[string]interface{}{
				"Path": 1,
			}
			tomlConfig := types.TransformersConfigForToml{
				Exporter:  map[string]interface{}{
					"test-exporter": notATransformerExporterMap,
				},
			}
			_, parseErr := config.ParseTransformerExporters(tomlConfig)
			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr.Error()).To(MatchRegexp("'Path' expected type 'string', got unconvertible type 'int'"))
		})
	})
})
