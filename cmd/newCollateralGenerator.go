package cmd

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/prompts"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/transformer_exporter"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	collateral            types.Collateral
	collateralErr         error
	configFileName        string
	configFilePath        string
	defaultConfigFilePath = helpers.GetEnvironmentsPath()
)

const (
	defaultConfigFileName = "mcdTransformers"
)

var addNewCollateralCmd = &cobra.Command{
	Use:   "addNewCollateral",
	Short: "Adds configuration to track a new collateral's contracts.",
	Long: `This command generates many of the essential changes that are needed for tracking a new collateral:
- adds flip contract details to mcdTransformers.toml
- tracks flip events for the new collateral
- creates a new storage transformer for the new collateral's flip contract
- if median and/or osm contracts exist for the new collateral:
	- adds contract details to mcdTransformers.toml
	- tracks events for the new collateral contracts
	- creates a new storage transformer for the new collateral's median contract
- updates the transformerExporter.go

There are a few changes that still need to be made manually. Those include:
- creating integration tests for the new collateral contracts
- adding the new contracts to helper methods in transformers/shared/constants/method.go
and transformers/test_data/config_values.go
`,
	Example: `Run ./vdb-mcd-transformers addNewCollateral and then follow the prompts.`,
	PreRun:  setViperConfigs,
	Run: func(cmd *cobra.Command, args []string) {
		err := addNewCollateral()
		if err != nil {
			logrus.Errorf("Failed to add %s: %s", collateral.Name, err)
			return
		}
		logrus.Infof("Successfully added %s config", collateral.Name)
		return
	},
}

func init() {
	addNewCollateralCmd.Flags().StringVarP(&configFileName, "config-file-name", "n",
		defaultConfigFileName, fmt.Sprintf("config file name, defaults to %s", defaultConfigFileName))
	addNewCollateralCmd.Flags().StringVarP(&configFilePath, "config-file-path", "p",
		defaultConfigFilePath, fmt.Sprintf("config path where the config file is expected to be, defaults to %s", defaultConfigFileName))
	rootCmd.AddCommand(addNewCollateralCmd)
}

func addNewCollateral() error {
	prompter := prompts.NewPrompter()
	collateral, collateralErr = prompter.GetCollateralDetails()
	if collateralErr != nil {
		return fmt.Errorf("failed to get collateral from command line prompts: %w", collateralErr)
	}

	contracts, contractsErr := prompter.GetContractDetails()
	if contractsErr != nil {
		return fmt.Errorf("failed to get contract details from command line prompts: %w", contractsErr)
	}

	configUpdater := config.NewConfigUpdater(collateral, contracts, prompter.MedianRequired, prompter.OsmRequired)
	initializerWriter := initializer.NewInitializerGenerator(helpers.GetProjectPath(), collateral, prompter.MedianRequired)
	newCollateralGenerator := new_collateral.NewCollateralGenerator{
		ConfigFileName:             configFileName,
		ConfigFilePath:             configFilePath,
		ConfigParser:               config.NewParser(),
		ConfigUpdater:              configUpdater,
		TransformerExporterUpdater: transformer_exporter.NewTransformerExporterUpdater(),
		InitializerGenerator:       &initializerWriter,
	}

	return newCollateralGenerator.Execute()
}
