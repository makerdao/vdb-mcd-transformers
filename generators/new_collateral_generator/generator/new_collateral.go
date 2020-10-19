package generator

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"
)

type NewCollateralGenerator struct {
	ConfigFileName string
	ConfigFilePath string
	ConfigParser   config.IParse
	ConfigUpdater  config.IUpdate
}

func (g NewCollateralGenerator) AddToConfig() error {
	initialConfig, parseConfigErr := g.ConfigParser.ParseCurrentConfig(g.ConfigFilePath, g.ConfigFileName)
	if parseConfigErr != nil {
		return parseConfigErr
	}

	g.ConfigUpdater.SetInitialConfig(initialConfig)

	updateErr := g.ConfigUpdater.AddNewCollateralToConfig()
	if updateErr != nil {
		return updateErr
	}

	file, fileOpenErr := os.Create(g.ConfigFilePath + g.ConfigFileName + ".toml")
	if fileOpenErr != nil {
		return fileOpenErr
	}

	updatedConfig := g.ConfigUpdater.GetUpdatedConfig()
	configToWrite := config.TransformersConfigForEncoding{
		ExporterMetadata:     updatedConfig.ExporterMetadata,
		TransformerExporters: updatedConfig.TransformerExporters,
		Contracts:            updatedConfig.Contracts,
	}

	encodingErr := toml.NewEncoder(file).Encode(configToWrite)
	if encodingErr != nil {
		return encodingErr
	}

	return file.Close()
}
