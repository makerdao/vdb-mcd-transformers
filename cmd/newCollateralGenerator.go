package cmd

import (
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral_generator/generator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	collateralName         string
	collateralVersion      string
	flipAddress            string
	flipAbi                string
	flipBlock              int
	medianContractRequired bool
	medianAddress          string
	medianAbi              string
	medianBlock            int
	osmContractRequired    bool
	osmAddress             string
	osmAbi                 string
	osmBlock               int
)

var addNewCollateralCmd = &cobra.Command{
	Use:     "addNewCollateral",
	Short:   "",
	Long:    ``,
	Example: ``,
	PreRun:  setViperConfigs,
	Run: func(cmd *cobra.Command, args []string) {
		err := addNewCollateral()
		if err != nil {
			logrus.Errorf("Failed to add new collateral to config: %w", err)
			return
		}
		logrus.Infof("Successfully added %s config", collateralName)
		return
	},
}

func init() {
	rootCmd.AddCommand(addNewCollateralCmd)
	addNewCollateralCmd.Flags().StringVarP(&collateralName, "collateral-name", "n", "", "new collateral name")
	addNewCollateralCmd.Flags().StringVarP(&collateralVersion, "collateral-version", "v", "", "new collateral version")
	addNewCollateralCmd.Flags().StringVar(&flipAddress, "flip-address", "", "new collateral's flip contract address")
	addNewCollateralCmd.Flags().StringVar(&flipAbi, "flip-abi", "", "new collateral's flip abi")
	addNewCollateralCmd.Flags().IntVar(&flipBlock, "flip-block", 0, "new collateral's flip deployed block")

	addNewCollateralCmd.Flags().BoolVar(&medianContractRequired, "median-contract-required", false, "pass this flag in when a median contract is required")
	addNewCollateralCmd.Flags().StringVar(&medianAddress, "median-address", "", "new collateral's median contract address")
	addNewCollateralCmd.Flags().StringVar(&medianAbi, "median-abi", "", "new collateral's median abi")
	addNewCollateralCmd.Flags().IntVar(&medianBlock, "median-block", 0, "new collateral's median deployed block")

	addNewCollateralCmd.Flags().BoolVar(&osmContractRequired, "osm-contract-required", false, "pass this flag in when a median contract is required")
	addNewCollateralCmd.Flags().StringVar(&osmAddress, "osm-address", "", "new collateral's osm contract address")
	addNewCollateralCmd.Flags().StringVar(&osmAbi, "osm-abi", "", "new collateral's osm abi")
	addNewCollateralCmd.Flags().IntVar(&osmBlock, "osm-block", 0, "new collateral's osm deployed block")
}

func addNewCollateral() error {
	println("Adding the new collateral to mcdTransformer.toml.")
	collateral := config.NewCollateral(collateralName, collateralVersion)
	flipContract := config.Contract{
		Address:  flipAddress,
		Abi:      flipAbi,
		Deployed: flipBlock,
	}
	medianContract := config.Contract{
		Address:  medianAddress,
		Abi:      medianAbi,
		Deployed: medianBlock,
	}
	osmContract := config.Contract{
		Address:  osmAddress,
		Abi:      osmAbi,
		Deployed: osmBlock,
	}
	contracts := make(map[string]config.Contract)
	contracts["flip"] = flipContract
	contracts["median"] = medianContract
	contracts["osm"] = osmContract

	configParser := config.Parser{}
	configUpdater := config.NewConfigUpdater(collateral, contracts, medianContractRequired, osmContractRequired)
	newCollateralGenerator := generator.NewCollateralGenerator{
		ConfigFileName: "mcdTransformers",
		ConfigFilePath: "../environments/",
		ConfigParser:   configParser,
		ConfigUpdater:  configUpdater,
	}
	addErr := newCollateralGenerator.AddToConfig()
	if addErr != nil {
		return addErr
	}
	return nil
}
