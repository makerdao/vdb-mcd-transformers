package types

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
)

type Collateral struct {
	Name    string
	Version string
}

func NewCollateral(name, version string) Collateral {
	return Collateral{
		Name:    name,
		Version: version,
	}
}

func (c Collateral) FormattedVersion() string {
	// example: v1_2_3
	return fmt.Sprintf("v%s", strings.Replace(c.Version, ".", "_", -1))
}

func (c Collateral) GetFlipTransformerName() string {
	// example: flip_eth_b_v1_2_3
	name := strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("flip_%s_%s", name, c.FormattedVersion())
}

func (c Collateral) GetMedianTransformerName() string {
	// example: median_eth_b
	name := strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("median_%s", name)
}

func (c Collateral) GetFlipInitializerDirectory() string {
	// example: eth_b/v1_2_3
	name := strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("%s/%s", name, c.FormattedVersion())
}

func (c Collateral) GetMedianInitializerDirectory() string {
	// example: median_eth_b
	name := strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("median_%s", name)
}

func (c Collateral) GetAbsoluteFlipStorageInitializersDirectoryPath() string {
	// example: $GOPATH/transformer/storage/flip/initializers/eth_b/v1_1_3
	return filepath.Join(
		helpers.GetProjectPath(), "transformers", "storage", "flip", "initializers", c.GetFlipInitializerDirectory())
}

func (c Collateral) GetAbsoluteMedianStorageInitializersDirectoryPath() string {
	// example: $GOPATH/transformer/storage/median/initializers/median_eth_b
	return filepath.Join(
		helpers.GetProjectPath(), "transformers", "storage", "median", "initializers", c.GetMedianInitializerDirectory())
}

func (c Collateral) GetAbsoluteFlipStorageInitializerFilePath() string {
	// example: $GOPATH/transformer/storage/flip/initializers/eth_b/v1_1_3/initializer.go
	return filepath.Join(
		helpers.GetProjectPath(), "transformers", "storage", "flip", "initializers", c.GetFlipInitializerDirectory(), "initializer.go")
}

func (c Collateral) GetAbsoluteMedianStorageInitializerFilePath() string {
	// example: $GOPATH/transformer/storage/median/initializers/median_eth_b/initializer.go
	return filepath.Join(
		helpers.GetProjectPath(), "transformers", "storage", "median", "initializers", c.GetMedianInitializerDirectory(), "initializer.go")
}

func (c Collateral) GetFlipContractName() string {
	// example: MCD_FLIP_ETH_B_1_1_3
	name := strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
	version := fmt.Sprintf("%s", strings.Replace(c.Version, ".", "_", -1))
	return fmt.Sprintf("MCD_FLIP_%s_%s", name, version)
}

func (c Collateral) GetMedianContractName() string {
	// example: MEDIAN_ETH_B
	name := strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("MEDIAN_%s", name)
}

func (c Collateral) GetOsmContractName() string {
	// example: MEDIAN_ETH_B
	name := strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
	return fmt.Sprintf("OSM_%s", name)
}
