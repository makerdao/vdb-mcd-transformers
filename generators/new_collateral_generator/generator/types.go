package generator

import (
	"fmt"
	"strings"
)

type Collateral struct {
	Name string
	Version string
}

func NewCollateral(name, version string) Collateral {
	return Collateral{
		Name:    name,
		Version: version,
	}
}

func (c Collateral) FormattedForFlipTransformerName() string{
	// example: eth_b_v1_2_3
	name := strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
	version := fmt.Sprintf("v%s", strings.Replace(c.Version, ".", "_", -1))
	return fmt.Sprintf("%s_%s", name, version)
}

func (c Collateral) FormattedForFlipInitializerFileName() string{
	// example: eth_b/v1_2_3
	name := strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
	version := fmt.Sprintf("v%s", strings.Replace(c.Version, ".", "_", -1))
	return fmt.Sprintf("%s/%s", name, version)
}

func (c Collateral) FormattedForMedianTransformerName() string{
	// example: eth_b
	return strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
}

func (c Collateral) FormattedForFlipContractName() string {
	// example: MCD_FLIP_ETH_B_1_1_3
	name := strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
	version := fmt.Sprintf("%s", strings.Replace(c.Version, ".", "_", -1))
	return fmt.Sprintf("MCD_FLIP_%s_%s", name, version)
}

func (c Collateral) FormattedForMedianContractName() string {
	// example: MEDIAN_ETH_B
	name := strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("MEDIAN_%s", name)
}

func (c Collateral) FormattedForOsmContractName() string {
	// example: MEDIAN_ETH_B
	name := strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("OSM_%s", name)
}

type Contract struct {
	Address  string
	Abi      string
	Deployed int
}

type Contracts map[string]Contract

type TransformersConfig struct {
	ExporterMetadata     ExporterMetaData `toml:"exporter"`
	TransformerExporters TransformerExporters
	Contracts            Contracts `toml:"contract"`
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
	Type       string `toml:"type"`
	Repository string
	Migrations string
	Contracts  []string
	Rank       string
}

type TransformerExporters map[string]TransformerExporter

