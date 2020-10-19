package config

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
)

type IUpdate interface {
	SetInitialConfig(initialConfig TransformersConfig)
	AddNewCollateralToConfig() error
	GetUpdatedConfig() TransformersConfig
}

type Updater struct {
	Collateral             Collateral
	Contracts              Contracts
	MedianContractRequired bool
	OsmContractRequired    bool
	InitialConfig          TransformersConfig
	UpdatedConfig          TransformersConfig
}

func NewConfigUpdater(collateral Collateral, contracts Contracts, medianContractRequired, osmContractRequired bool) *Updater {
	return &Updater{
		Collateral:             collateral,
		Contracts:              contracts,
		MedianContractRequired: medianContractRequired,
		OsmContractRequired:    osmContractRequired,
	}
}

func (cg *Updater) SetInitialConfig(initialConfig TransformersConfig) {
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

	var updatedConfig TransformersConfig
	decoder := gob.NewDecoder(buf)
	decErr := decoder.Decode(&updatedConfig)
	if decErr != nil {
		return decErr
	}

	cg.UpdatedConfig = updatedConfig
	return nil
}

func (cg *Updater) addStorageTransformerNames() {
	flipTransformerName := "flip_" + cg.Collateral.FormattedForFlipTransformerName()
	newTransformerNames := []string{flipTransformerName}
	if cg.MedianContractRequired {
		medianTransformerName := "median_" + cg.Collateral.FormattedForMedianTransformerName()
		newTransformerNames = []string{flipTransformerName, medianTransformerName}
	}

	cg.UpdatedConfig.ExporterMetadata.TransformerNames = append(
		cg.UpdatedConfig.ExporterMetadata.TransformerNames,
		newTransformerNames...,
	)
}

func (cg *Updater) addStorageExporters() {
	flipStorageExporter := TransformerExporter{
		Path:       fmt.Sprintf("transformers/storage/flip/initializers/%s", cg.Collateral.FormattedForFlipInitializerFileName()),
		Type:       "eth_storage",
		Repository: "github.com/makerdao/vdb-mcd-transformers",
		Migrations: "db/migrations",
		Rank:       "0",
	}
	transformerExporters := make(map[string]TransformerExporter)
	flipKey := fmt.Sprintf("flip_%s", cg.Collateral.FormattedForFlipTransformerName())
	transformerExporters[flipKey] = flipStorageExporter

	if cg.MedianContractRequired {
		medianStorageExporter := TransformerExporter{
			Path:       fmt.Sprintf("transformers/storage/median/initializers/median_%s", cg.Collateral.FormattedForMedianTransformerName()),
			Type:       "eth_storage",
			Repository: "github.com/makerdao/vdb-mcd-transformers",
			Migrations: "db/migrations",
			Rank:       "0",
		}
		medianKey := fmt.Sprintf("median_%s", cg.Collateral.FormattedForMedianTransformerName())
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
	if cg.MedianContractRequired {
		medianErr := cg.addNewContractToMedianExporters()
		if medianErr != nil {
			return medianErr
		}
	}

	if cg.OsmContractRequired {
		osmErr := cg.addNewContractToOsmExporters()
		if osmErr != nil {
			return osmErr
		}
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
	return cg.addNewContractToExporters(IsFlipExporter, cg.Collateral.FormattedForFlipContractName)
}

func (cg *Updater) addNewContractToMedianExporters() error {
	return cg.addNewContractToExporters(IsMedianExporter, cg.Collateral.FormattedForMedianContractName)
}

func (cg *Updater) addNewContractToOsmExporters() error {
	return cg.addNewContractToExporters(IsOsmExporter, cg.Collateral.FormattedForOsmContractName)
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
				continue
			}
		}
		cg.UpdatedConfig.TransformerExporters[name] = exporter
	}

	return nil
}

func (cg *Updater) addContracts() {
	formattedContracts := make(map[string]Contract)

	flipContractKey := cg.Collateral.FormattedForFlipContractName()
	formattedContracts[flipContractKey] = cg.Contracts["flip"]

	if cg.MedianContractRequired {
		medianContractKey := cg.Collateral.FormattedForMedianContractName()
		formattedContracts[medianContractKey] = cg.Contracts["median"]
	}

	if cg.OsmContractRequired {
		osmContractKey := cg.Collateral.FormattedForOsmContractName()
		formattedContracts[osmContractKey] = cg.Contracts["osm"]
	}

	for k, v := range formattedContracts {
		cg.UpdatedConfig.Contracts[k] = v
	}
}

func (cg *Updater) GetUpdatedConfig() TransformersConfig {
	//configToWrite := config.TransformersConfig{
	//	ExporterMetadata:     g.ConfigUpdater.UpdatedConfig.ExporterMetadata,
	//	TransformerExporters: g.ConfigUpdater.UpdatedConfig.TransformerExporters,
	//	Contracts:            g.ConfigUpdater.UpdatedConfig.Contracts,
	//}
	return cg.UpdatedConfig
}
