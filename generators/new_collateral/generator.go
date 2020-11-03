package new_collateral

import (
	"fmt"
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
		return fmt.Errorf("error updating config: %w", addErr)
	}

	logrus.Infof("Adding new collateral to %s", helpers.GetExecutePluginsPath())
	updatePluginErr := g.updatePluginExporter()
	if updatePluginErr != nil {
		return fmt.Errorf("error updating tranformerExporter.go: %w", updatePluginErr)
	}

	logrus.Info("Writing initializers for new collateral")
	writeInitializerErr := g.writeInitializers()
	if writeInitializerErr != nil {
		return fmt.Errorf("error creating initializer files: %w", writeInitializerErr)
	}
	return nil
}

func (g NewCollateralGenerator) updateConfig() error {
	initialConfig, parseConfigErr := g.ConfigParser.ParseCurrentConfig(g.ConfigFilePath, g.ConfigFileName)
	if parseConfigErr != nil {
		return fmt.Errorf("error parsing current config: %w", parseConfigErr)
	}

	g.ConfigUpdater.SetInitialConfig(initialConfig)

	updateErr := g.ConfigUpdater.AddNewCollateralToConfig()
	if updateErr != nil {
		return fmt.Errorf("error adding new collateral to config: %w", updateErr)
	}

	file, fileOpenErr := os.Create(helpers.GetFullConfigFilePath(g.ConfigFilePath, g.ConfigFileName))
	if fileOpenErr != nil {
		return fmt.Errorf("opening config file: %w", fileOpenErr)
	}

	updatedConfig, updatedConfigErr := g.ConfigUpdater.GetUpdatedConfigForToml()
	if updatedConfigErr != nil {
		return fmt.Errorf("formatting config for toml: %w", updatedConfigErr)
	}

	encodingErr := toml.NewEncoder(file).Encode(updatedConfig)
	if encodingErr != nil {
		return fmt.Errorf("error encoding config for toml: %w", encodingErr)
	}

	return file.Close()
}

func (g *NewCollateralGenerator) updatePluginExporter() error {
	updatedConfig := g.ConfigUpdater.GetUpdatedConfig()
	pluginConfig, pluginErr := g.TransformerExporterUpdater.PreparePluginConfig(updatedConfig)
	if pluginErr != nil {
		return fmt.Errorf("error preparing plugin config: %w", pluginErr)
	}

	return g.TransformerExporterUpdater.WritePlugin(pluginConfig)
}

func (g *NewCollateralGenerator) writeInitializers() error {
	flipInitializersErr := g.InitializerGenerator.GenerateFlipInitializer()
	if flipInitializersErr != nil {
		return fmt.Errorf("error generating flip initializer: %w", flipInitializersErr)
	}
	return g.InitializerGenerator.GenerateMedianInitializer()
}
