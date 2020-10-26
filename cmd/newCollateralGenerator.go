package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/config"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	collateral             types.Collateral
	medianContractRequired bool
	osmContractRequired    bool
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
	var collateralErr error
	collateral, collateralErr = getCollateral()
	if collateralErr != nil {
		return collateralErr
	}

	contracts, getContractsErr := getContracts()
	if getContractsErr != nil {
		return getContractsErr
	}

	configUpdater := config.NewConfigUpdater(collateral, contracts, medianContractRequired, osmContractRequired)
	configFileName := "mcdTransformers"
	configFilePath := helpers.GetEnvironmentsPath()
	initializerWriter := initializer.Generator{
		Collateral:                collateral,
		MedianInitializerRequired: medianContractRequired,
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

func getCollateral() (types.Collateral, error) {
	namePrompt := promptui.Prompt{
		Label:    "Collateral Name",
		Validate: validateString,
	}

	versionPrompt := promptui.Prompt{
		Label:    "Collateral Version",
		Validate: validateString,
	}

	collateralName, nameErr := namePrompt.Run()
	if nameErr != nil {
		fmt.Printf("Prompt failed %v\n", nameErr)
		return types.Collateral{}, nameErr
	}
	collateralVersion, versionErr := versionPrompt.Run()
	if versionErr != nil {
		fmt.Printf("Prompt failed %v\n", versionErr)
		return types.Collateral{}, versionErr
	}

	return types.NewCollateral(collateralName, collateralVersion), nil
}

func validateString(input string) error {
	if input == "" {
		return errors.New("This string field is required")
	}
	return nil
}

func validateInt(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("Invalid number")
	}
	return nil
}

func getContracts() (map[string]types.Contract, error) {
	contracts := make(map[string]types.Contract)

	flipContract, flipErr := getContract("Flip")
	if flipErr != nil {
		return contracts, flipErr
	}
	contracts["flip"] = flipContract

	contracts, medianErr := addMedianContract(contracts)
	if medianErr != nil {
		return contracts, medianErr
	}

	contracts, osmErr := addOsmContract(contracts)
	if osmErr != nil {
		return contracts, osmErr
	}
	return contracts, nil
}

func addMedianContract(contracts map[string]types.Contract) (map[string]types.Contract, error) {
	var medianRequiredErr error
	medianContractRequired, medianRequiredErr = isContractRequired("Median")
	if medianRequiredErr != nil {
		return contracts, medianRequiredErr
	}

	if medianContractRequired {
		medianContract, medianContractErr := getContract("Median")
		if medianContractErr != nil {
			return contracts, medianContractErr
		}
		contracts["median"] = medianContract
	}

	return contracts, nil
}

func addOsmContract(contracts map[string]types.Contract) (map[string]types.Contract, error) {
	var osmRequiredErr error
	osmContractRequired, osmRequiredErr = isContractRequired("OSM")
	if osmRequiredErr != nil {
		return contracts, osmRequiredErr
	}

	if osmContractRequired {
		osmContract, osmContractErr := getContract("OSM")
		if osmContractErr != nil {
			return contracts, osmContractErr
		}
		contracts["osm"] = osmContract
	}

	return contracts, nil
}

func getContract(contractType string) (types.Contract, error) {
	addressPrompt := promptui.Prompt{
		Label:    fmt.Sprintf("%s Contract Address", contractType),
		Validate: validateString,
	}

	abiPrompt := promptui.Prompt{
		Label:    fmt.Sprintf("%s Contract ABI", contractType),
		Validate: validateString,
	}

	blockPrompt := promptui.Prompt{
		Label:    fmt.Sprintf("%s Contract deployment block", contractType),
		Validate: validateInt,
	}

	address, addressErr := addressPrompt.Run()
	if addressErr != nil {
		return types.Contract{}, addressErr
	}
	abi, abiErr := abiPrompt.Run()
	if abiErr != nil {
		return types.Contract{}, abiErr
	}
	blockResult, blockErr := blockPrompt.Run()
	if blockErr != nil {
		return types.Contract{}, blockErr
	}
	block, intErr := strconv.Atoi(blockResult)
	if intErr != nil {
		return types.Contract{}, intErr
	}

	return types.Contract{
		Address:  address,
		Abi:      abi,
		Deployed: block,
	}, nil
}

func isContractRequired(contractType string) (bool, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Is a %s contract required?", contractType),
		Items: []bool{true, false},
	}
	_, result, resultErr := prompt.Run()
	if resultErr != nil {
		return false, resultErr
	}
	return strconv.ParseBool(result)
}
