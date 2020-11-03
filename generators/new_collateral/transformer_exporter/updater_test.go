package transformer_exporter_test

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/test_data"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/transformer_exporter"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	pluginConfig "github.com/makerdao/vulcanizedb/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginWriter", func() {
	Context("PreparePluginConfig", func() {
		var pluginWriter = transformer_exporter.NewTransformerExporterUpdater()
		It("Formats the updated TransformerConfig into a plugin config struct", func() {
			updatedConfig := test_data.UpdatedConfig
			config, prepareErr := pluginWriter.PreparePluginConfig(updatedConfig)
			Expect(prepareErr).NotTo(HaveOccurred())

			expectedPluginConfig := pluginConfig.Plugin{
				Transformers: map[string]pluginConfig.Transformer{
					"test-1": {
						Path:           "path-test-1",
						Type:           pluginConfig.EthStorage,
						MigrationPath:  "test-migrations",
						MigrationRank:  0,
						RepositoryPath: "repo-1",
					},
					"test-2": {
						Path:           "path-test-2",
						Type:           pluginConfig.EthEvent,
						MigrationPath:  "test-migrations",
						MigrationRank:  0,
						RepositoryPath: "repo-2",
					},
				},
				FilePath: helpers.GetExecutePluginsPath(),
				FileName: test_data.UpdatedConfig.ExporterMetadata.Name,
				Save:     test_data.UpdatedConfig.ExporterMetadata.Save,
				Home:     test_data.UpdatedConfig.ExporterMetadata.Home,
				Schema:   test_data.UpdatedConfig.ExporterMetadata.Schema,
			}
			Expect(config).To(Equal(expectedPluginConfig))
		})

		It("returns an error if converting rank string to an int fails", func() {
			updatedConfig := types.TransformersConfig{
				ExporterMetadata: types.ExporterMetaData{},
				TransformerExporters: types.TransformerExporters{
					"test": types.TransformerExporter{
						Rank: "a",
					},
				},
			}
			_, prepareErr := pluginWriter.PreparePluginConfig(updatedConfig)
			Expect(prepareErr).To(HaveOccurred())
		})
	})
})
