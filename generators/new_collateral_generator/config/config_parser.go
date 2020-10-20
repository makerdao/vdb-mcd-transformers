package config

import (
	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator"
	"github.com/spf13/viper"
)

type IParse interface {
	ParseCurrentConfig(configFilePath, configFileName string) (TransformersConfig, error)
}
type Parser struct{}

func (Parser) ParseCurrentConfig(configFilePath, configFileName string) (TransformersConfig, error) {
	viperConfig := viper.New()
	viperConfig.AddConfigPath(configFilePath)
	viperConfig.SetConfigName(configFileName)
	readConfigErr := viperConfig.ReadInConfig()
	if readConfigErr != nil {
		return TransformersConfig{}, readConfigErr
	}

	var tomlConfig TransformersConfig
	fullConfigFilePath := new_collateral_generator.GetFullConfigFilePath(configFilePath, configFileName)
	_, decodeErr := toml.DecodeFile(fullConfigFilePath, &tomlConfig)
	if decodeErr != nil {
		return TransformersConfig{}, decodeErr
	}

	//TODO: if we update the config file format to separate the exporter metadata from the exporters, this step should
	// no longer be necessary - toml.DecodeFile should be able to properly decode those exporters
	exporters, exporterErr := getTransformerExporters(tomlConfig.ExporterMetadata.TransformerNames, viperConfig)
	if exporterErr != nil {
		return TransformersConfig{}, exporterErr
	}

	tomlConfig.TransformerExporters = exporters
	return tomlConfig, nil
}

func getTransformerExporters(transformerNames []string, viperConfig *viper.Viper) (TransformerExporters, error) {
	exporters := make(TransformerExporters)
	for _, transformerName := range transformerNames {
		transformer := viperConfig.GetStringMapString("exporter." + transformerName)
		transformerContracts := viperConfig.GetStringMapStringSlice("exporter." + transformerName)
		//TODO: handle errors in case one of the values doesn't exist
		exporters[transformerName] = TransformerExporter{
			Path:       transformer["path"],
			Type:       transformer["type"],
			Repository: transformer["repository"],
			Migrations: transformer["migrations"],
			Contracts:  transformerContracts["contracts"],
			Rank:       transformer["rank"],
		}
	}
	return exporters, nil
}
