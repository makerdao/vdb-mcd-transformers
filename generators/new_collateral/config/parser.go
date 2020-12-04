package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/mitchellh/mapstructure"
)

type Parser interface {
	ParseCurrentConfig(configFilePath, configFileName string) (types.TransformersConfig, error)
}

type parser struct{}

func NewParser() parser {
	return parser{}
}

var (
	ErrorDecodingConfigFile          = errors.New("error decoding config file")
	ErrorParsingExporterMetadata     = errors.New("error parsing exporter metadata from config file")
	ErrorParsingTransformerExporters = errors.New("error parsing transformer exporters from config file")
)

func (parser) ParseCurrentConfig(configFilePath, configFileName string) (types.TransformersConfig, error) {
	var tomlConfig types.TransformersConfigForToml
	fullConfigFilePath := helpers.GetFullConfigFilePath(configFilePath, configFileName)
	_, decodeErr := toml.DecodeFile(fullConfigFilePath, &tomlConfig)
	if decodeErr != nil {
		return types.TransformersConfig{}, fmt.Errorf("%w: %s", ErrorDecodingConfigFile, decodeErr)
	}

	metadata, metadataErr := parseExporterMetaData(tomlConfig)
	if metadataErr != nil {
		return types.TransformersConfig{}, fmt.Errorf("%w: %s", ErrorParsingExporterMetadata, metadataErr.Error())
	}

	transformerExporters, transformerExportersErr := parseTransformerExporters(tomlConfig)
	if transformerExportersErr != nil {
		return types.TransformersConfig{}, fmt.Errorf("%w: %s", ErrorParsingTransformerExporters, transformerExportersErr)
	}

	return types.TransformersConfig{
		ExporterMetadata:     metadata,
		TransformerExporters: transformerExporters,
		Contracts:            tomlConfig.Contracts,
	}, nil
}

func parseExporterMetaData(tomlConfig types.TransformersConfigForToml) (types.ExporterMetaData, error) {
	home, homeOk := tomlConfig.Exporter["home"].(string)
	name, nameOk := tomlConfig.Exporter["name"].(string)
	save, saveOk := tomlConfig.Exporter["save"].(bool)
	schema, schemaOk := tomlConfig.Exporter["schema"].(string)
	if !homeOk || !nameOk || !saveOk || !schemaOk {
		return types.ExporterMetaData{}, fmt.Errorf(
			"error asserting exporterMetadata types - homeOk: %t, nameOk: %t, saveOk: %t, schemaOk: %t",
			homeOk, nameOk, saveOk, schemaOk,
		)
	}

	var transformerNames []string
	decodeErr := mapstructure.Decode(tomlConfig.Exporter["transformerNames"], &transformerNames)
	if decodeErr != nil {
		return types.ExporterMetaData{}, fmt.Errorf("error decoding transformerNames: %w", decodeErr)
	}

	return types.ExporterMetaData{
		Home:             home,
		Name:             name,
		Save:             save,
		Schema:           schema,
		TransformerNames: transformerNames,
	}, nil
}

func parseTransformerExporters(tomlConfig types.TransformersConfigForToml) (types.TransformerExporters, error) {
	var exporters = make(map[string]types.TransformerExporter)
	for exporterKey, exporterValue := range tomlConfig.Exporter {
		if keyIsForMetadata(exporterKey) {
			continue
		} else {
			var result types.TransformerExporter
			decodeErr := mapstructure.Decode(exporterValue, &result)
			if decodeErr != nil {
				return types.TransformerExporters{}, fmt.Errorf("error decoding transformerExporters: %w", decodeErr)
			}

			exporters[exporterKey] = result
		}
	}

	return exporters, nil
}

func keyIsForMetadata(key string) bool {
	return key == "home" || key == "name" || key == "save" || key == "schema" || key == "transformerNames"
}
