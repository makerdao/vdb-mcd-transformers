package generator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
)

type ConfigGenerator struct {
	Collateral    Collateral
	Contracts     Contracts
	InitialConfig TransformersConfig
	UpdatedConfig TransformersConfig
}

func NewConfigGenerator(collateral Collateral, contracts Contracts, initialConfig TransformersConfig) *ConfigGenerator {
	return &ConfigGenerator{
		Collateral: collateral,
		Contracts: contracts,
		InitialConfig: initialConfig,
	}
}


func (cg *ConfigGenerator) AddNewCollateralToConfig() error {
	copyErr := cg.copyInitialConfig()
	if copyErr != nil {
		return copyErr
	}

	cg.addStorageTransformerNames()
	cg.addStorageExporters()
	cg.addContractsToEventExporters()
	cg.addContracts()

	return nil
}

func (cg *ConfigGenerator) copyInitialConfig() error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encErr := encoder.Encode(cg.InitialConfig)
	if encErr != nil {
		return encErr
	}

	var updatedConfig TransformersConfig
	decoder := gob.NewDecoder(buf)
	decErr := decoder.Decode(&updatedConfig)
	if decErr != nil {
		return decErr
	}

	cg.UpdatedConfig = updatedConfig
	return nil
}

func (cg *ConfigGenerator) addStorageTransformerNames() {
	flipTransformerName := "flip_" + cg.Collateral.FormattedForFlipTransformerName()
	medianTransformerName := "median_" + cg.Collateral.FormattedForMedianTransformerName()

	newTransformerNames := []string{flipTransformerName, medianTransformerName}

	cg.UpdatedConfig.ExporterMetadata.TransformerNames = append(
		cg.UpdatedConfig.ExporterMetadata.TransformerNames,
		newTransformerNames...
	)
}

func (cg *ConfigGenerator) addStorageExporters() {
	flipStorageExporter := TransformerExporter{
		Path:       fmt.Sprintf("transformers/storage/flip/initializers/%s", cg.Collateral.FormattedForFlipInitializerFileName()),
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}

	medianStorageExporter := TransformerExporter{
		Path:       fmt.Sprintf("transformers/storage/median/initializers/median_%s", cg.Collateral.FormattedForMedianTransformerName()),
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}
	transformerExporters := make(map[string]TransformerExporter)
	flipKey := fmt.Sprintf("exporter.flip_%s", cg.Collateral.FormattedForFlipTransformerName())
	transformerExporters[flipKey] = flipStorageExporter
	medianKey := fmt.Sprintf("exporter.median_%s", cg.Collateral.FormattedForMedianTransformerName())
	transformerExporters[medianKey] = medianStorageExporter

	for k, v := range transformerExporters {
		cg.UpdatedConfig.TransformerExporters[k] = v
	}
}

func (cg *ConfigGenerator) addContractsToEventExporters() error {
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

func (cg *ConfigGenerator) addNewContractToFlipExporters() error {
	return cg.addNewContractToExporters(IsFlipExporter, cg.Collateral.FormattedForFlipContractName)
}

func (cg *ConfigGenerator) addNewContractToMedianExporters() error {
	return cg.addNewContractToExporters(IsMedianExporter, cg.Collateral.FormattedForMedianContractName)
}

func (cg *ConfigGenerator) addNewContractToOsmExporters() error {
	return cg.addNewContractToExporters(IsOsmExporter, cg.Collateral.FormattedForOsmContractName)
}

type matcherFunc func(string) (bool, error)
type collateralFormatter func() string

func (cg *ConfigGenerator) addNewContractToExporters(matcherFunc matcherFunc, collateralFormatter collateralFormatter) error {
	for name, exporter := range cg.UpdatedConfig.TransformerExporters {
		for _, contract := range exporter.Contracts {
			contractTypeMatched, matchErr := matcherFunc(contract)
			if matchErr != nil {
				return matchErr
			}
			if contractTypeMatched {
				exporter.Contracts = append(exporter.Contracts, collateralFormatter())
				continue
			}
		}
		cg.UpdatedConfig.TransformerExporters[name] = exporter
	}

	return nil
}


func (cg *ConfigGenerator) addContracts() {
	formattedContracts := make(map[string]Contract)

	flipContractKey := cg.Collateral.FormattedForFlipContractName()
	formattedContracts[flipContractKey] = cg.Contracts["flip"]

	medianContractKey := cg.Collateral.FormattedForMedianContractName()
	formattedContracts[medianContractKey] = cg.Contracts["median"]

	osmContractKey := cg.Collateral.FormattedForOsmContractName()
	formattedContracts[osmContractKey] = cg.Contracts["osm"]

	for k, v := range formattedContracts {
		cg.UpdatedConfig.Contracts[k] = v
	}
}

