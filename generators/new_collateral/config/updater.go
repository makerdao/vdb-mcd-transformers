package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

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

func (cu *Updater) SetInitialConfig(initialConfig types.TransformersConfig) {
	cu.InitialConfig = initialConfig
}

func (cu *Updater) AddNewCollateralToConfig() error {
	copyErr := cu.copyInitialConfig()
	if copyErr != nil {
		return fmt.Errorf("error copying initialConfig to a new struct: %w", copyErr)
	}

	cu.addStorageTransformerNames()
	cu.addStorageExporters()
	addContractsToExportersErr := cu.addContractsToEventExporters()
	if addContractsToExportersErr != nil {
		return fmt.Errorf("error adding contracts to event exporters: %w", addContractsToExportersErr)
	}
	cu.addContracts()

	return nil
}

func (cu *Updater) copyInitialConfig() error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encErr := encoder.Encode(cu.InitialConfig)
	if encErr != nil {
		return fmt.Errorf("error encoding initial config: %w", encErr)
	}

	var updatedConfig types.TransformersConfig
	decoder := gob.NewDecoder(buf)
	decErr := decoder.Decode(&updatedConfig)
	if decErr != nil {
		return fmt.Errorf("error decoding updated config: %w", decErr)
	}

	cu.UpdatedConfig = updatedConfig
	return nil
}

func (cu *Updater) addStorageTransformerNames() {
	flipTransformerName := cu.Collateral.GetFlipTransformerName()
	newTransformerNames := []string{flipTransformerName}
	if cu.MedianContractRequired {
		medianTransformerName := cu.Collateral.GetMedianTransformerName()
		newTransformerNames = []string{flipTransformerName, medianTransformerName}
	}

	cu.UpdatedConfig.ExporterMetadata.TransformerNames = append(
		cu.UpdatedConfig.ExporterMetadata.TransformerNames,
		newTransformerNames...,
	)
}

func (cu *Updater) addStorageExporters() {
	flipStorageExporter := types.TransformerExporter{
		Path:       fmt.Sprintf("transformers/storage/flip/initializers/%s", cu.Collateral.GetFlipInitializerDirectory()),
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}
	transformerExporters := make(map[string]types.TransformerExporter)
	flipKey := cu.Collateral.GetFlipTransformerName()
	transformerExporters[flipKey] = flipStorageExporter

	if cu.MedianContractRequired {
		medianStorageExporter := types.TransformerExporter{
			Path:       fmt.Sprintf("transformers/storage/median/initializers/%s", cu.Collateral.GetMedianInitializerDirectory()),
			Type:       "eth_storage",
			Repository: "github.com/makerdao/vdb-mcd-transformers",
			Migrations: "db/migrations",
			Rank:       "0",
		}
		medianKey := cu.Collateral.GetMedianTransformerName()
		transformerExporters[medianKey] = medianStorageExporter
	}

	for k, v := range transformerExporters {
		cu.UpdatedConfig.TransformerExporters[k] = v
	}
}

func (cu *Updater) addContractsToEventExporters() error {
	flipErr := cu.addNewContractToFlipExporters()
	if flipErr != nil {
		return fmt.Errorf("error adding new contract to flip exporters: %w", flipErr)
	}

	medianErr := cu.addNewContractToMedianExporters()
	if medianErr != nil {
		return fmt.Errorf("error adding new contract to median exporters: %w", medianErr)
	}

	osmErr := cu.addNewContractToOsmExporters()
	if osmErr != nil {
		return fmt.Errorf("error adding new contract to osm exporters: %w", osmErr)
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

func (cu *Updater) addNewContractToFlipExporters() error {
	return cu.addNewContractToExporters(IsFlipExporter, cu.Collateral.GetFlipContractName)
}

func (cu *Updater) addNewContractToMedianExporters() error {
	if cu.MedianContractRequired {
		return cu.addNewContractToExporters(IsMedianExporter, cu.Collateral.GetMedianContractName)
	}
	return nil
}

func (cu *Updater) addNewContractToOsmExporters() error {
	if cu.OsmContractRequired {
		return cu.addNewContractToExporters(IsOsmExporter, cu.Collateral.GetOsmContractName)
	}
	return nil
}

type matcherFunc func(string) (bool, error)
type collateralFormatter func() string

func (cu *Updater) addNewContractToExporters(matcherFunc matcherFunc, collateralFormatter collateralFormatter) error {
	for name, exporter := range cu.UpdatedConfig.TransformerExporters {
		for _, contract := range exporter.Contracts {
			contractTypeMatched, matchErr := matcherFunc(contract)
			if matchErr != nil {
				return fmt.Errorf("error matching contract type: %w", matchErr)
			}
			if contractTypeMatched {
				exporter.Contracts = append(exporter.Contracts, collateralFormatter())
				// breaking out of this exporter's contract loop so that new collateral contract is not added multiple times
				break
			}
		}
		cu.UpdatedConfig.TransformerExporters[name] = exporter
	}

	return nil
}

func (cu *Updater) addContracts() {
	formattedContracts := make(map[string]types.Contract)

	flipContractKey := cu.Collateral.GetFlipContractName()
	formattedContracts[flipContractKey] = cu.Contracts["flip"]

	if cu.MedianContractRequired {
		medianContractKey := cu.Collateral.GetMedianContractName()
		formattedContracts[medianContractKey] = cu.Contracts["median"]
	}

	if cu.OsmContractRequired {
		osmContractKey := cu.Collateral.GetOsmContractName()
		formattedContracts[osmContractKey] = cu.Contracts["osm"]
	}

	for k, v := range formattedContracts {
		cu.UpdatedConfig.Contracts[k] = v
	}
}

func (cu *Updater) GetUpdatedConfig() types.TransformersConfig {
	return cu.UpdatedConfig
}

// GetUpdatedConfigForToml converts TransformersConfig.ExporterMetadata and TransformerExporter structs into a
// map[string]interface{} to allow for proper toml encoding when writing to the config file^
func (cu *Updater) GetUpdatedConfigForToml() (types.TransformersConfigForToml, error) {
	metadataMap, metadataMapErr := convertToLowerCaseStringToInterfaceMap(cu.UpdatedConfig.ExporterMetadata)
	if metadataMapErr != nil {
		return types.TransformersConfigForToml{}, fmt.Errorf("error converting metadata map: %w", metadataMapErr)
	}

	transformerExporterMap, transformerExporterMapErr := convertToLowerCaseStringToInterfaceMap(cu.UpdatedConfig.TransformerExporters)
	if transformerExporterMapErr != nil {
		return types.TransformersConfigForToml{}, fmt.Errorf("error converting transformer exporter map: %w", transformerExporterMapErr)
	}

	exporterMap := mergeMaps(metadataMap, transformerExporterMap)

	return types.TransformersConfigForToml{
		Exporter:  exporterMap,
		Contracts: cu.UpdatedConfig.Contracts,
	}, nil
}

func convertToLowerCaseStringToInterfaceMap(input interface{}) (map[string]interface{}, error) {
	var stringToInterfaceMap map[string]interface{}
	jsonBytes, _ := json.Marshal(input)
	unmarshalErr := json.Unmarshal(jsonBytes, &stringToInterfaceMap)
	if unmarshalErr != nil {
		return stringToInterfaceMap, fmt.Errorf("error unmarshaling interface map: %w", unmarshalErr)
	}
	return makeKeysLowercase(stringToInterfaceMap), nil
}

func makeKeysLowercase(inputMap map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{})
	for key, value := range inputMap {
		lowerCaseKey := downCaseFirstCharacter(key)
		valueMap, valueIsMap := value.(map[string]interface{})
		if valueIsMap {
			valueMap = makeKeysLowercase(valueMap)
			resultMap[lowerCaseKey] = valueMap
		} else {
			resultMap[lowerCaseKey] = value
		}
	}
	return resultMap
}

func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	for key, value := range map1 {
		map2[key] = value
	}
	return map2
}

func downCaseFirstCharacter(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
