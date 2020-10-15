package generator

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type ConfigParser struct {}

func (cp *ConfigParser) ParseCurrentConfig(configString string) (TransformersConfig, error) {
	//TODO: maybe just pass this in as a byte array?
	configBytes := []byte(configString)
	viperConfig := viper.New()
	viperConfig.SetConfigType("toml")
	readConfigErr := viperConfig.ReadConfig(bytes.NewBuffer(configBytes))
	if readConfigErr != nil {
		return TransformersConfig{}, readConfigErr
	}

	var tomlConfig TransformersConfig
	_, decodeErr := toml.Decode(configString, &tomlConfig)
	if decodeErr != nil {
		return TransformersConfig{}, decodeErr
	}

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
		exporters["exporter."+transformerName] = TransformerExporter{
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
