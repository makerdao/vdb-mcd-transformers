package types

import (
	"fmt"
	"strings"
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

func (c Collateral) formattedLowerCaseName() string {
	return strings.ToLower(strings.Replace(c.Name, "-", "_", -1))
}

func (c Collateral) formattedUpperCaseName() string {
	return strings.ToUpper(strings.Replace(c.Name, "-", "_", -1))
}

func (c Collateral) FormattedLowerCaseNameWithoutVersionLetter() string {
	// example: ETH_B => ETH
	return strings.Split(c.formattedLowerCaseName(), "_")[0]
}

func (c Collateral) FormattedUpperCaseNameWithoutVersionLetter() string {
	// example: ETH_B => ETH
	return strings.Split(c.formattedUpperCaseName(), "_")[0]
}

func (c Collateral) FormattedVersionWithPrependedV() string {
	// example: v1_2_3
	return fmt.Sprintf("v%s", c.FormattedVersion())
}

func (c Collateral) FormattedVersion() string {
	// example: 1_2_3
	return fmt.Sprintf("%s", strings.Replace(c.Version, ".", "_", -1))
}

func (c Collateral) GetFlipTransformerName() string {
	// example: flip_eth_b_v1_2_3
	return fmt.Sprintf("flip_%s_%s", c.formattedLowerCaseName(), c.FormattedVersionWithPrependedV())
}

func (c Collateral) GetMedianTransformerName() string {
	// example: median_eth_v1_2_3
	return fmt.Sprintf("median_%s_%s", c.FormattedLowerCaseNameWithoutVersionLetter(), c.FormattedVersionWithPrependedV())
}

func (c Collateral) GetFlipInitializerDirectory() string {
	// example: eth_b/v1_2_3
	return fmt.Sprintf("%s/%s", c.formattedLowerCaseName(), c.FormattedVersionWithPrependedV())
}

func (c Collateral) GetMedianInitializerDirectory() string {
	// example: median_eth_b
	return fmt.Sprintf("median_%s/%s", c.FormattedLowerCaseNameWithoutVersionLetter(), c.FormattedVersionWithPrependedV())
}

func (c Collateral) GetFlipContractName() string {
	// example: MCD_FLIP_ETH_B_1_1_3
	version := fmt.Sprintf("%s", strings.Replace(c.Version, ".", "_", -1))
	return fmt.Sprintf("MCD_FLIP_%s_%s", c.formattedUpperCaseName(), version)
}

func (c Collateral) GetMedianContractName() string {
	// example: MEDIAN_ETH_1_1_3
	version := fmt.Sprintf("%s", strings.Replace(c.Version, ".", "_", -1))
	return fmt.Sprintf("MEDIAN_%s_%s", c.FormattedUpperCaseNameWithoutVersionLetter(), version)
}

func (c Collateral) GetOsmContractName() string {
	// example: OSM_ETH_B
	return fmt.Sprintf("OSM_%s", c.FormattedUpperCaseNameWithoutVersionLetter())
}
