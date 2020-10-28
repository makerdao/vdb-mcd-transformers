package cmd

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/prompts"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	collateral types.Collateral
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
			logrus.Error("Failed to add new collateral to config: ", err)
			return
		}
		logrus.Infof("Successfully added %s config", collateral.Name)
		return
	},
}

func init() {
	rootCmd.AddCommand(addNewCollateralCmd)
}

func addNewCollateral() error {
	prompter := prompts.NewPrompter()
	collateral, collateralErr := prompter.GetCollateralDetails()
	if collateralErr != nil {
		return collateralErr
	}

	contracts, contractsErr := prompter.GetContractDetails()
	if contractsErr != nil {
		return contractsErr
	}

	configUpdater := config.NewConfigUpdater(collateral, contracts, prompter.MedianRequired, prompter.OsmRequired)
	configFileName := "mcdTransformers"
	configFilePath := helpers.GetEnvironmentsPath()
	initializerWriter := initializer.Generator{
		Collateral:                collateral,
		MedianInitializerRequired: prompter.MedianRequired,
	}
	newCollateralGenerator := new_collateral.NewCollateralGenerator{
		ConfigFileName:       configFileName,
		ConfigFilePath:       configFilePath,
		ConfigParser:         config.Parser{},
		ConfigUpdater:        configUpdater,
		InitializerGenerator: &initializerWriter,
	}

	fmt.Println(fmt.Sprintf("Adding %s to %s", collateral.Name, helpers.GetFullConfigFilePath(configFileName, configFilePath)))
	addErr := newCollateralGenerator.UpdateConfig()
	if addErr != nil {
		return addErr
	}

	fmt.Println(fmt.Sprintf("Adding %s to %s", collateral.Name, helpers.GetExecutePluginsPath()))
	updatePluginErr := newCollateralGenerator.UpdatePluginExporter()
	if updatePluginErr != nil {
		return updatePluginErr
	}

	fmt.Println(fmt.Sprintf("Writing initializers for %s", collateral.Name))
	writeInitializerErr := newCollateralGenerator.WriteInitializers()
	if writeInitializerErr != nil {
		return writeInitializerErr
	}

	return nil
}
