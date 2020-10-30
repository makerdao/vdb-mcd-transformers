package new_collateral_test

import (
	"io/ioutil"
	"os"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/test_data"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewCollateral", func() {
	var (
		filePath                   = "./"
		fileName                   = "test"
		fullConfigPath             = filePath + fileName + ".toml"
		configParser               test_data.MockConfigParser
		configUpdater              test_data.MockConfigUpdater
		transformerExporterUpdater test_data.MockTransformerExporterUpdater
		initializerGenerator       test_data.MockInitializerGenerator
		collateralGenerator        new_collateral.NewCollateralGenerator
	)

	BeforeEach(func() {
		configParser = test_data.MockConfigParser{}
		configUpdater = test_data.MockConfigUpdater{}
		transformerExporterUpdater = test_data.MockTransformerExporterUpdater{}
		initializerGenerator = test_data.MockInitializerGenerator{}
		collateralGenerator = new_collateral.NewCollateralGenerator{
			ConfigFileName:             fileName,
			ConfigFilePath:             filePath,
			ConfigParser:               &configParser,
			ConfigUpdater:              &configUpdater,
			TransformerExporterUpdater: &transformerExporterUpdater,
			InitializerGenerator:       &initializerGenerator,
		}
	})

	Describe("Execute", func() {
		Context("update config", func() {
			It("parses the current config", func() {
				err := collateralGenerator.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(configParser.ConfigFilePathPassedIn).To(Equal(filePath))
				Expect(configParser.ConfigFileNamePassedIn).To(Equal(fileName))
			})

			It("returns an error if parsing the config fails", func() {
				configParser.ParseErr = fakes.FakeError
				err := collateralGenerator.Execute()
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
				err := collateralGenerator.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(configUpdater.SetCurrentConfigCalled).To(BeTrue())
				Expect(configUpdater.InitialConfigPassedIn).To(Equal(testConfig))
			})

			It("adds new the collateral to the current config", func() {
				err := collateralGenerator.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(configUpdater.AddNewCollateralCalled).To(BeTrue())
			})

			It("returns an error if adding to the current config fails", func() {
				configUpdater.AddNewCollateralErr = fakes.FakeError
				err := collateralGenerator.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("gets updated config from the updater to write to the config file", func() {
				err := collateralGenerator.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(configUpdater.GetUpdatedConfigForTomlCalled).To(BeTrue())
			})

			It("returns an error if getting the updated config fails", func() {
				configUpdater.GetUpdatedConfigForTomlCalledErr = fakes.FakeError
				err := collateralGenerator.Execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fakes.FakeError))
			})

			It("writes the updated collateral to the config file", func() {
				configUpdater.UpdatedConfigForToml = test_data.UpdatedConfigForToml
				addErr := collateralGenerator.Execute()
				Expect(addErr).NotTo(HaveOccurred())

				testConfigContent, readErr := ioutil.ReadFile(fullConfigPath)
				Expect(readErr).NotTo(HaveOccurred())
				Expect(string(testConfigContent)).To(Equal(test_data.TestConfigFileContent))

				removeErr := os.Remove(fullConfigPath)
				Expect(removeErr).NotTo(HaveOccurred())
			})
		})

		Context("writes the transformer exporter file", func() {
			It("prepares the plugin config", func() {
				err := collateralGenerator.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(transformerExporterUpdater.PreparePluginConfigCalled).To(BeTrue())
				Expect(configUpdater.GetUpdatedConfigCalled).To(BeTrue())
			})

			It("writes the transformer exporter file with the VDB plugin writer", func() {
				pluginConfigToWrite := config.Plugin{
					FilePath: "filePath",
					FileName: "fileName",
					Save:     true,
					Home:     "home",
					Schema:   "schema",
				}
				transformerExporterUpdater.PluginConfigToReturn = pluginConfigToWrite
				err := collateralGenerator.Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(transformerExporterUpdater.WritePluginCalled).To(BeTrue())
				Expect(transformerExporterUpdater.PluginConfigPassedIn).To(Equal(pluginConfigToWrite))
			})
		})

		Context("write initializer files", func() {
			It("writes the flip initializer file", func() {
				initializerErr := collateralGenerator.Execute()
				Expect(initializerErr).NotTo(HaveOccurred())
				Expect(initializerGenerator.GenerateFlipInitializerCalled).To(BeTrue())
			})

			It("returns an error if writing the flip initializer fails", func() {
				initializerGenerator.FlipInitializerErr = fakes.FakeError
				initializerErr := collateralGenerator.Execute()
				Expect(initializerErr).To(HaveOccurred())
				Expect(initializerErr).To(MatchError(fakes.FakeError))
			})

			It("writes the median initializer file", func() {
				initializerErr := collateralGenerator.Execute()
				Expect(initializerErr).NotTo(HaveOccurred())
				Expect(initializerGenerator.GenerateMedianInitializerCalled).To(BeTrue())
			})

			It("returns an error if writing the median initializer fails", func() {
				initializerGenerator.MedianInitializerErr = fakes.FakeError
				initializerErr := collateralGenerator.Execute()
				Expect(initializerErr).To(HaveOccurred())
				Expect(initializerErr).To(MatchError(fakes.FakeError))
			})
		})
	})
})
