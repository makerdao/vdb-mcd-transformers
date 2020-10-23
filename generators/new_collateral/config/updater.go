package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

type IUpdate interface {
	SetInitialConfig(initialConfig types.TransformersConfig)
	AddNewCollateralToConfig() error
	GetUpdatedConfig() types.TransformersConfig
	GetUpdatedConfigForToml() (types.TransformersConfigForToml, error)
}

type Updater struct {
	Collateral             types.Collateral
	Contracts              types.Contracts
	MedianContractRequired bool
	OsmContractRequired    bool
	InitialConfig          types.TransformersConfig
	UpdatedConfig          types.TransformersConfig
}

func NewConfigUpdater(collateral types.Collateral, contracts types.Contracts, medianContractRequired, osmContractRequired bool) *Updater {
	return &Updater{
		Collateral:             collateral,
		Contracts:              contracts,
		MedianContractRequired: medianContractRequired,
		OsmContractRequired:    osmContractRequired,
	}
}

func (cg *Updater) SetInitialConfig(initialConfig types.TransformersConfig) {
	cg.InitialConfig = initialConfig
}

func (cg *Updater) AddNewCollateralToConfig() error {
	copyErr := cg.copyInitialConfig()
	if copyErr != nil {
		return copyErr
	}

	cg.addStorageTransformerNames()
	cg.addStorageExporters()
	addContractsToExportersErr := cg.addContractsToEventExporters()
	if addContractsToExportersErr != nil {
		return addContractsToExportersErr
	}
	cg.addContracts()

	return nil
}

func (cg *Updater) copyInitialConfig() error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encErr := encoder.Encode(cg.InitialConfig)
	if encErr != nil {
		return encErr
	}

	var updatedConfig types.TransformersConfig
	decoder := gob.NewDecoder(buf)
	decErr := decoder.Decode(&updatedConfig)
	if decErr != nil {
		return decErr
	}

	cg.UpdatedConfig = updatedConfig
	return nil
}

func (cg *Updater) addStorageTransformerNames() {
	flipTransformerName := cg.Collateral.GetFlipTransformerName()
	newTransformerNames := []string{flipTransformerName}
	if cg.MedianContractRequired {
		medianTransformerName := cg.Collateral.GetMedianTransformerName()
		newTransformerNames = []string{flipTransformerName, medianTransformerName}
	}

	cg.UpdatedConfig.ExporterMetadata.TransformerNames = append(
		cg.UpdatedConfig.ExporterMetadata.TransformerNames,
		newTransformerNames...,
	)
}

func (cg *Updater) addStorageExporters() {
	flipStorageExporter := types.TransformerExporter{
		Path:       fmt.Sprintf("transformers/storage/flip/initializers/%s", cg.Collateral.GetFlipInitializerDirectory()),
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}
	transformerExporters := make(map[string]types.TransformerExporter)
	flipKey := cg.Collateral.GetFlipTransformerName()
	transformerExporters[flipKey] = flipStorageExporter

	if cg.MedianContractRequired {
		medianStorageExporter := types.TransformerExporter{
			Path:       fmt.Sprintf("transformers/storage/median/initializers/%s", cg.Collateral.GetMedianTransformerName()),
			Type:       "eth_storage",
			Repository: "github.com/makerdao/vdb-mcd-transformers",
			Migrations: "db/migrations",
			Rank:       "0",
		}
		medianKey := cg.Collateral.GetMedianTransformerName()
		transformerExporters[medianKey] = medianStorageExporter
	}

	for k, v := range transformerExporters {
		cg.UpdatedConfig.TransformerExporters[k] = v
	}
}

func (cg *Updater) addContractsToEventExporters() error {
	flipErr := cg.addNewContractToFlipExporters()
	if flipErr != nil {
		return flipErr
	}

	medianErr := cg.addNewContractToMedianExporters()
	if medianErr != nil {
		return medianErr
	}

	osmErr := cg.addNewContractToOsmExporters()
	if osmErr != nil {
		return osmErr
	}

	return nil
}

func IsFlipExporter(contractName string) (bool, error) {
	return regexp.Match("FLIP", []byte(contractName))
}

func IsMedianExporter(contractName string) (bool, error) {
	return regexp.Match("MEDIAN", []byte(contractName))
}

func IsOsmExporter(contractName string) (bool, error) {
	return regexp.Match("OSM", []byte(contractName))
}

func (cg *Updater) addNewContractToFlipExporters() error {
	return cg.addNewContractToExporters(IsFlipExporter, cg.Collateral.GetFlipContractName)
}

func (cg *Updater) addNewContractToMedianExporters() error {
	if cg.MedianContractRequired {
		return cg.addNewContractToExporters(IsMedianExporter, cg.Collateral.GetMedianContractName)
	}
	return nil
}

func (cg *Updater) addNewContractToOsmExporters() error {
	if cg.OsmContractRequired {
		return cg.addNewContractToExporters(IsOsmExporter, cg.Collateral.GetOsmContractName)
	}
	return nil
}

type matcherFunc func(string) (bool, error)
type collateralFormatter func() string

func (cg *Updater) addNewContractToExporters(matcherFunc matcherFunc, collateralFormatter collateralFormatter) error {
	for name, exporter := range cg.UpdatedConfig.TransformerExporters {
		for _, contract := range exporter.Contracts {
			contractTypeMatched, matchErr := matcherFunc(contract)
			if matchErr != nil {
				return matchErr
			}
			if contractTypeMatched {
				exporter.Contracts = append(exporter.Contracts, collateralFormatter())
				// breaking out of this exporter's contract loop so that new collateral contract is not added multiple times
				break
			}
		}
		cg.UpdatedConfig.TransformerExporters[name] = exporter
	}

	return nil
}

func (cg *Updater) addContracts() {
	formattedContracts := make(map[string]types.Contract)

	flipContractKey := cg.Collateral.GetFlipContractName()
	formattedContracts[flipContractKey] = cg.Contracts["flip"]

	if cg.MedianContractRequired {
		medianContractKey := cg.Collateral.GetMedianContractName()
		formattedContracts[medianContractKey] = cg.Contracts["median"]
	}

	if cg.OsmContractRequired {
		osmContractKey := cg.Collateral.GetOsmContractName()
		formattedContracts[osmContractKey] = cg.Contracts["osm"]
	}

	for k, v := range formattedContracts {
		cg.UpdatedConfig.Contracts[k] = v
	}
}

func (cg *Updater) GetUpdatedConfig() types.TransformersConfig {
	return cg.UpdatedConfig
}

func (cg *Updater) GetUpdatedConfigForToml() (types.TransformersConfigForToml, error) {
	var metaDataMap map[string]interface{}
	metadata, _ := json.Marshal(cg.UpdatedConfig.ExporterMetadata)
	metadataUnmarshalErr := json.Unmarshal(metadata, &metaDataMap)
	if metadataUnmarshalErr != nil {
		return types.TransformersConfigForToml{}, metadataUnmarshalErr
	}

	var exporterMap map[string]interface{}
	exporters, _ := json.Marshal(cg.UpdatedConfig.TransformerExporters)
	exportersUnmarshalErr := json.Unmarshal(exporters, &exporterMap)
	if exportersUnmarshalErr != nil {
		return types.TransformersConfigForToml{}, metadataUnmarshalErr
	}

	for field, val := range exporterMap {
		metaDataMap[field] = val
	}

	return types.TransformersConfigForToml{
		Exporter:  metaDataMap,
		Contracts: cg.UpdatedConfig.Contracts,
	}, nil
}
