package config

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/mitchellh/mapstructure"
)

type IParse interface {
	ParseCurrentConfig(configFilePath, configFileName string) (types.TransformersConfig, error)
}
type Parser struct{}

func (Parser) ParseCurrentConfig(configFilePath, configFileName string) (types.TransformersConfig, error) {
	var tomlConfig types.TransformersConfigForToml
	fullConfigFilePath := helpers.GetFullConfigFilePath(configFilePath, configFileName)
	_, decodeErr := toml.DecodeFile(fullConfigFilePath, &tomlConfig)
	if decodeErr != nil {
		return types.TransformersConfig{}, decodeErr
	}

	metadata, metadataErr := parseExporterMetaData(tomlConfig)
	if metadataErr != nil {
		return types.TransformersConfig{}, metadataErr
	}

	transformerExporters, transformerExportersErr := parseTransformerExporters(tomlConfig)
	if transformerExportersErr != nil {
		return types.TransformersConfig{}, transformerExportersErr
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
		return types.ExporterMetaData{}, errors.New("error asserting exporter meta data types")
	}

	var transformerNames []string
	decodeErr := mapstructure.Decode(tomlConfig.Exporter["transformerNames"], &transformerNames)
	if decodeErr != nil {
		return types.ExporterMetaData{}, decodeErr
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
	delete(tomlConfig.Exporter, "home")
	delete(tomlConfig.Exporter, "name")
	delete(tomlConfig.Exporter, "save")
	delete(tomlConfig.Exporter, "schema")
	delete(tomlConfig.Exporter, "transformerNames")

	var exporters = make(map[string]types.TransformerExporter)
	for exporterKey, exporterValue := range tomlConfig.Exporter {
		var result types.TransformerExporter
		decodeErr := mapstructure.Decode(exporterValue, &result)
		if decodeErr != nil {
			return types.TransformerExporters{}, decodeErr
		}

		exporters[exporterKey] = result
	}

	return exporters, nil
}
