package new_collateral_test

import (
	"io/ioutil"
	"os"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/test_data"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	pluginConfig "github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCollateral", func() {
	var (
		filePath            = "./"
		fileName            = "test"
		fullConfigPath      = filePath + fileName + ".toml"
		configParser        test_data.MockConfigParser
		configUpdater       test_data.MockConfigUpdater
		collateralGenerator new_collateral.NewCollateralGenerator
	)

	BeforeEach(func() {
		configParser = test_data.MockConfigParser{}
		configUpdater = test_data.MockConfigUpdater{}
		collateralGenerator = new_collateral.NewCollateralGenerator{
			ConfigFileName: fileName,
			ConfigFilePath: filePath,
			ConfigParser:   &configParser,
			ConfigUpdater:  &configUpdater,
		}

	})
	Context("AddToConfig", func() {
		It("parses the current config", func() {
			err := collateralGenerator.AddToConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(configParser.ConfigFilePathPassedIn).To(Equal(filePath))
			Expect(configParser.ConfigFileNamePassedIn).To(Equal(fileName))
		})

		It("returns an error if parsing the config fails", func() {
			configParser.ParseErr = fakes.FakeError
			err := collateralGenerator.AddToConfig()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("sets the current config on the config updater", func() {
			testConfig := types.TransformersConfig{
				ExporterMetadata: types.ExporterMetaData{
					Home: "test",
				},
			}
			configParser.ConfigToReturn = testConfig
			err := collateralGenerator.AddToConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(configUpdater.SetCurrentConfigCalled).To(BeTrue())
			Expect(configUpdater.InitialConfigPassedIn).To(Equal(testConfig))
		})

		It("adds new the collateral to the current config", func() {
			err := collateralGenerator.AddToConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(configUpdater.AddNewCollateralCalled).To(BeTrue())
		})

		It("returns an error if adding to the current config fails", func() {
			configUpdater.AddNewCollateralErr = fakes.FakeError
			err := collateralGenerator.AddToConfig()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("gets updated config from the updater to write to the config file", func() {
			err := collateralGenerator.AddToConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(configUpdater.GetUpdatedConfigForTomlCalled).To(BeTrue())
		})

		It("returns an error if getting the updated config fails", func() {
			configUpdater.GetUpdatedConfigForTomlCalledErr = fakes.FakeError
			err := collateralGenerator.AddToConfig()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("writes the updated collateral to the config file", func() {
			_, createErr := os.Create(fullConfigPath)
			Expect(createErr).NotTo(HaveOccurred())

			configUpdater.UpdatedConfigForToml = test_data.UpdatedConfigForToml
			addErr := collateralGenerator.AddToConfig()
			Expect(addErr).NotTo(HaveOccurred())

			testConfigContent, readErr := ioutil.ReadFile(fullConfigPath)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(testConfigContent)).To(Equal(test_data.TestConfigFileContent))

			removeErr := os.Remove(fullConfigPath)
			Expect(removeErr).NotTo(HaveOccurred())
		})
	})

	Context("UpdatePluginExporter", func() {
		It("prepares the plugin.Config using the updated transformers config", func() {
			configUpdater.UpdatedConfig = test_data.UpdatedConfig
			config, pluginErr := collateralGenerator.PreparePluginConfig()
			Expect(pluginErr).NotTo(HaveOccurred())

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
	})
})
