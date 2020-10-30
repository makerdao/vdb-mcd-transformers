package new_collateral

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/transformer_exporter"
	"github.com/sirupsen/logrus"
)

type NewCollateralGenerator struct {
	ConfigFileName             string
	ConfigFilePath             string
	ConfigParser               config.Parser
	ConfigUpdater              config.IUpdate
	TransformerExporterUpdater transformer_exporter.TransformerExporterUpdater
	InitializerGenerator       initializer.IGenerate
}

func (g NewCollateralGenerator) Execute() error {
	logrus.Infof("Adding new collateral to %s", helpers.GetFullConfigFilePath(g.ConfigFileName, g.ConfigFilePath))
	addErr := g.updateConfig()
	if addErr != nil {
		return addErr
	}

	logrus.Infof("Adding new collateral to %s", helpers.GetExecutePluginsPath())
	updatePluginErr := g.updatePluginExporter()
	if updatePluginErr != nil {
		return updatePluginErr
	}

	logrus.Info("Writing initializers for new collateral")
	writeInitializerErr := g.writeInitializers()
	if writeInitializerErr != nil {
		return writeInitializerErr
	}
	return nil
}

func (g NewCollateralGenerator) updateConfig() error {
	initialConfig, parseConfigErr := g.ConfigParser.ParseCurrentConfig(g.ConfigFilePath, g.ConfigFileName)
	if parseConfigErr != nil {
		return parseConfigErr
	}

	g.ConfigUpdater.SetInitialConfig(initialConfig)

	updateErr := g.ConfigUpdater.AddNewCollateralToConfig()
	if updateErr != nil {
		return updateErr
	}

	file, fileOpenErr := os.Create(helpers.GetFullConfigFilePath(g.ConfigFilePath, g.ConfigFileName))
	if fileOpenErr != nil {
		return fileOpenErr
	}

	updatedConfig, updatedConfigErr := g.ConfigUpdater.GetUpdatedConfigForToml()
	if updatedConfigErr != nil {
		return updatedConfigErr
	}

	encodingErr := toml.NewEncoder(file).Encode(updatedConfig)
	if encodingErr != nil {
		return encodingErr
	}

	return file.Close()
}

func (g *NewCollateralGenerator) updatePluginExporter() error {
	updatedConfig := g.ConfigUpdater.GetUpdatedConfig()
	pluginConfig, pluginErr := g.TransformerExporterUpdater.PreparePluginConfig(updatedConfig)
	if pluginErr != nil {
		return pluginErr
	}

	return g.TransformerExporterUpdater.WritePlugin(pluginConfig)
}

func (g *NewCollateralGenerator) writeInitializers() error {
	flipInitializersErr := g.InitializerGenerator.GenerateFlipInitializer()
	if flipInitializersErr != nil {
		return flipInitializersErr
	}
	return g.InitializerGenerator.GenerateMedianInitializer()
}
