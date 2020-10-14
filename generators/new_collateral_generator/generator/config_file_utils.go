package generator

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

func getTransformerExporters(transformerNames []string, viperConfig *viper.Viper) (TransformerExporters, error) {
	exporters := make(TransformerExporters)
	for _, transformerName := range transformerNames {
		transformer := viperConfig.GetStringMapString("exporter." + transformerName)
		transformerContracts := viperConfig.GetStringMapStringSlice("exporter." + transformerName)
		//TODO: handle errors in case one of the values doesn't exist
		exporters["exporter."+transformerName] = TransformerExporter{
			Path:       transformer["path"],
			Kind:       transformer["type"],
			Repository: transformer["repository"],
			Migrations: transformer["migrations"],
			Contracts:  transformerContracts["contracts"],
			Rank:       transformer["rank"],
		}
	}
	return exporters, nil
}

func ParseCurrentConfig(configString string) (TomlConfig, error) {
	//TODO: maybe just pass this in as a byte array?
	configBytes := []byte(configString)
	viperConfig := viper.New()
	viperConfig.SetConfigType("toml")
	readConfigErr := viperConfig.ReadConfig(bytes.NewBuffer(configBytes))
	if readConfigErr != nil {
		return TomlConfig{}, readConfigErr
	}

	var tomlConfig TomlConfig
	_, err := toml.Decode(configString, &tomlConfig)
	if err != nil {
		return TomlConfig{}, err
	}

	exporters, exporterErr := getTransformerExporters(tomlConfig.ExporterMetadata.TransformerNames, viperConfig)
	if exporterErr != nil {
		return TomlConfig{}, exporterErr
	}

	tomlConfig.TransformerExporters = exporters

	return tomlConfig, nil
}

type TomlConfig struct {
	ExporterMetadata     ExporterMetaData `toml:"exporter"`
	TransformerExporters TransformerExporters
	Contracts            Contracts `toml:"contract"`
	Servers              map[string]server
}
type server struct {
	IP string
	DC string
}
type ExporterMetaData struct {
	Home             string
	Name             string
	Save             bool
	Schema           string
	TransformerNames []string
}

type TransformerExporter struct {
	Path       string
	Kind       string `toml:"type"`
	Repository string
	Migrations string
	Contracts  []string
	Rank       string
}

type TransformerExporters map[string]TransformerExporter

type Contract struct {
	Address  string
	Abi      string
	Deployed int
}

type Contracts map[string]Contract
